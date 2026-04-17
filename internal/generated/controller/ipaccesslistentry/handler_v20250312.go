// Copyright 2025 MongoDB Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ipaccesslistentry

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strings"
	"time"

	v20250312sdk "go.mongodb.org/atlas-sdk/v20250312018/admin"
	k8smeta "k8s.io/apimachinery/pkg/api/meta"
	controllerruntime "sigs.k8s.io/controller-runtime"
	builder "sigs.k8s.io/controller-runtime/pkg/builder"
	client "sigs.k8s.io/controller-runtime/pkg/client"
	controller "sigs.k8s.io/controller-runtime/pkg/controller"
	reconcile "sigs.k8s.io/controller-runtime/pkg/reconcile"

	akov2generated "github.com/mongodb/mongodb-atlas-kubernetes/v2/generated/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/customresource"
	ctrlstate "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/state"
	crapi "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/crapi"
	result "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/result"
	state "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/state"
)

// expiredStateMsg is stored in the State condition's Message field to signal that the entry
// has expired. checkExpiry reads it back on subsequent reconciles so we avoid hitting Atlas
// once the entry is past its deleteAfterDate.
const expiredStateMsg = "Expired"

const atlasAccessListNotFound = "NOT_IN_IP_RANGE"
const atlasAccessListEntryNotFound = "ATLAS_NETWORK_PERMISSION_ENTRY_NOT_FOUND"

type Handlerv20250312 struct {
	kubeClient         client.Client
	atlasClient        *v20250312sdk.APIClient
	translator         crapi.Translator
	deletionProtection bool
}

func NewHandlerv20250312(kubeClient client.Client, atlasClient *v20250312sdk.APIClient, translator crapi.Translator, deletionProtection bool) *Handlerv20250312 {
	return &Handlerv20250312{
		atlasClient:        atlasClient,
		deletionProtection: deletionProtection,
		kubeClient:         kubeClient,
		translator:         translator,
	}
}

// HandleInitial creates a new IP access list entry in Atlas.
func (h *Handlerv20250312) HandleInitial(ctx context.Context, ipaccesslistentry *akov2generated.IPAccessListEntry) (ctrlstate.Result, error) {
	groupID, entryValue, err := h.resolveIdentity(ctx, ipaccesslistentry)
	if err != nil {
		return result.Error(state.StateInitial, fmt.Errorf("failed to resolve IPAccessListEntry identity: %w", err))
	}

	entry := buildNetworkPermissionEntry(ipaccesslistentry)
	_, _, err = h.atlasClient.ProjectIPAccessListApi.CreateAccessListEntry(ctx, groupID, &[]v20250312sdk.NetworkPermissionEntry{entry}).Execute()
	if err != nil {
		return result.Error(state.StateInitial, fmt.Errorf("failed to create IP access list entry %q: %w", entryValue, err))
	}

	if err := h.persistIdentity(ctx, ipaccesslistentry, groupID, entryValue); err != nil {
		return result.Error(state.StateInitial, err)
	}

	return result.NextState(state.StateCreating, "IP access list entry created. Waiting for activation.")
}

// HandleImportRequested imports an existing Atlas IP access list entry.
// The annotation mongodb.com/external-id must be set to the entry value (IP address, CIDR block or AWS security group)
func (h *Handlerv20250312) HandleImportRequested(ctx context.Context, ipaccesslistentry *akov2generated.IPAccessListEntry) (ctrlstate.Result, error) {
	entryValue, ok := ipaccesslistentry.GetAnnotations()["mongodb.com/external-id"]
	if !ok || entryValue == "" {
		return result.Error(state.StateImportRequested, errors.New("missing annotation mongodb.com/external-id: set it to the entry value (IP, CIDR, or AWS SG)"))
	}

	groupID, err := h.resolveGroupID(ctx, ipaccesslistentry)
	if err != nil {
		return result.Error(state.StateImportRequested, fmt.Errorf("failed to resolve groupId: %w", err))
	}

	_, _, err = h.atlasClient.ProjectIPAccessListApi.GetAccessListEntry(ctx, groupID, entryValue).Execute()
	if err != nil {
		return result.Error(state.StateImportRequested, fmt.Errorf("failed to get IP access list entry %q in group %q: %w", entryValue, groupID, err))
	}

	if err := h.persistIdentity(ctx, ipaccesslistentry, groupID, entryValue); err != nil {
		return result.Error(state.StateImportRequested, err)
	}

	return result.NextState(state.StateImported, "IP access list entry imported.")
}

// HandleImported handles the imported state.
func (h *Handlerv20250312) HandleImported(ctx context.Context, ipaccesslistentry *akov2generated.IPAccessListEntry) (ctrlstate.Result, error) {
	if checkExpiry(ipaccesslistentry) {
		return expiredResult(state.StateImported), nil
	}

	return h.handleSteadyState(ctx, state.StateImported, ipaccesslistentry)
}

// HandleCreating polls Atlas until the entry is ACTIVE.
func (h *Handlerv20250312) HandleCreating(ctx context.Context, ipaccesslistentry *akov2generated.IPAccessListEntry) (ctrlstate.Result, error) {
	groupID, entryValue, err := h.identityFromStatusOrDeps(ctx, ipaccesslistentry)
	if err != nil {
		return result.Error(state.StateCreating, err)
	}

	atlasStatus, err := h.getEntryStatus(ctx, groupID, entryValue)
	if err != nil {
		return result.Error(state.StateCreating, fmt.Errorf("failed to get IP access list status: %w", err))
	}

	if atlasStatus != "ACTIVE" {
		return result.NextState(state.StateCreating, fmt.Sprintf("IP access list entry is %s. Waiting for activation.", atlasStatus))
	}

	deps, err := h.getDependencies(ctx, ipaccesslistentry)
	if err != nil {
		return result.Error(state.StateCreating, fmt.Errorf("failed to resolve dependencies: %w", err))
	}

	copy := ipaccesslistentry.DeepCopy()
	if err := ctrlstate.NewPatcher(copy).UpdateStateTracker(deps...).Patch(ctx, h.kubeClient); err != nil {
		return result.Error(state.StateCreating, fmt.Errorf("failed to update state tracker: %w", err))
	}

	// If the entry has a deleteAfterDate, schedule a requeue instead of relying on watch events
	// (which are filtered by predicates when only metadata/status changes). This ensures
	// HandleCreated fires and sets up the expiry polling loop.
	if hasDeleteAfterDate(ipaccesslistentry) {
		return ctrlstate.Result{
			Result:    reconcile.Result{RequeueAfter: requeueAfterExpiry(ipaccesslistentry)},
			NextState: state.StateCreated,
			StateMsg:  "IP access list entry is active.",
		}, nil
	}

	return result.NextState(state.StateCreated, "IP access list entry is active.")
}

// HandleCreated handles the created steady state.
func (h *Handlerv20250312) HandleCreated(ctx context.Context, ipaccesslistentry *akov2generated.IPAccessListEntry) (ctrlstate.Result, error) {
	if checkExpiry(ipaccesslistentry) {
		return expiredResult(state.StateCreated), nil
	}

	return h.handleSteadyState(ctx, state.StateCreated, ipaccesslistentry)
}

// HandleUpdating polls Atlas until the entry is ACTIVE after an update.
func (h *Handlerv20250312) HandleUpdating(ctx context.Context, ipaccesslistentry *akov2generated.IPAccessListEntry) (ctrlstate.Result, error) {
	groupID, entryValue, err := h.identityFromStatusOrDeps(ctx, ipaccesslistentry)
	if err != nil {
		return result.Error(state.StateUpdating, err)
	}

	atlasStatus, err := h.getEntryStatus(ctx, groupID, entryValue)
	if err != nil {
		return result.Error(state.StateUpdating, fmt.Errorf("failed to get IP access list status: %w", err))
	}

	if atlasStatus != "ACTIVE" {
		return result.NextState(state.StateUpdating, fmt.Sprintf("IP access list entry is %s. Waiting for activation.", atlasStatus))
	}

	deps, err := h.getDependencies(ctx, ipaccesslistentry)
	if err != nil {
		return result.Error(state.StateUpdating, fmt.Errorf("failed to resolve dependencies: %w", err))
	}

	copy := ipaccesslistentry.DeepCopy()
	if err := ctrlstate.NewPatcher(copy).UpdateStateTracker(deps...).Patch(ctx, h.kubeClient); err != nil {
		return result.Error(state.StateUpdating, fmt.Errorf("failed to update state tracker: %w", err))
	}

	if hasDeleteAfterDate(ipaccesslistentry) {
		return ctrlstate.Result{
			Result:    reconcile.Result{RequeueAfter: requeueAfterExpiry(ipaccesslistentry)},
			NextState: state.StateUpdated,
			StateMsg:  "IP access list entry is active.",
		}, nil
	}

	return result.NextState(state.StateUpdated, "IP access list entry is up to date.")
}

// HandleUpdated handles the updated steady state.
func (h *Handlerv20250312) HandleUpdated(ctx context.Context, ipaccesslistentry *akov2generated.IPAccessListEntry) (ctrlstate.Result, error) {
	if checkExpiry(ipaccesslistentry) {
		return expiredResult(state.StateUpdated), nil
	}

	return h.handleSteadyState(ctx, state.StateUpdated, ipaccesslistentry)
}

// HandleDeletionRequested deletes the IP access list entry from Atlas.
func (h *Handlerv20250312) HandleDeletionRequested(ctx context.Context, ipaccesslistentry *akov2generated.IPAccessListEntry) (ctrlstate.Result, error) {
	if customresource.IsResourcePolicyKeepOrDefault(ipaccesslistentry, h.deletionProtection) {
		return result.NextState(state.StateDeleted, "IP access list entry skipped deletion due to retention policy.")
	}

	groupID, entryValue, err := h.identityFromStatusOrDeps(ctx, ipaccesslistentry)
	if err != nil {
		return result.Error(state.StateDeletionRequested, err)
	}

	_, err = h.atlasClient.ProjectIPAccessListApi.DeleteAccessListEntry(ctx, groupID, entryValue).Execute()
	if v20250312sdk.IsErrorCode(err, atlasAccessListNotFound) || v20250312sdk.IsErrorCode(err, atlasAccessListEntryNotFound) {
		return result.NextState(state.StateDeleted, "IP access list entry deleted.")
	}
	if err != nil {
		return result.Error(state.StateDeletionRequested, fmt.Errorf("failed to delete IP access list entry %q: %w", entryValue, err))
	}

	return result.NextState(state.StateDeleted, "IP access list entry deleted.")
}

// HandleDeleting handles the deleting state.
func (h *Handlerv20250312) HandleDeleting(ctx context.Context, ipaccesslistentry *akov2generated.IPAccessListEntry) (ctrlstate.Result, error) {
	groupID, entryValue, err := h.identityFromStatusOrDeps(ctx, ipaccesslistentry)
	if err != nil {
		return result.Error(state.StateDeleting, err)
	}

	_, _, err = h.atlasClient.ProjectIPAccessListApi.GetAccessListEntry(ctx, groupID, entryValue).Execute()
	if v20250312sdk.IsErrorCode(err, atlasAccessListNotFound) || v20250312sdk.IsErrorCode(err, atlasAccessListEntryNotFound) {
		return result.NextState(state.StateDeleted, "IP access list entry deleted.")
	}
	if err != nil {
		return result.Error(state.StateDeleting, fmt.Errorf("failed to check IP access list deletion: %w", err))
	}

	return result.NextState(state.StateDeleting, "Waiting for IP access list entry deletion.")
}

// For returns the resource and predicates for the controller.
func (h *Handlerv20250312) For() (client.Object, builder.Predicates) {
	return &akov2generated.IPAccessListEntry{}, builder.WithPredicates()
}

// SetupWithManager is not used for version-specific handlers but required by the StateHandler interface.
func (h *Handlerv20250312) SetupWithManager(mgr controllerruntime.Manager, rec reconcile.Reconciler, defaultOptions controller.Options) error {
	return nil
}

// handleSteadyState checks whether the spec changed and recreates the entry if needed.
func (h *Handlerv20250312) handleSteadyState(ctx context.Context, currentState state.ResourceState, ipaccesslistentry *akov2generated.IPAccessListEntry) (ctrlstate.Result, error) {
	deps, err := h.getDependencies(ctx, ipaccesslistentry)
	if err != nil {
		return result.Error(currentState, fmt.Errorf("failed to resolve IPAccessListEntry dependencies: %w", err))
	}

	// First, check if the entry has expired
	if hasDeleteAfterDate(ipaccesslistentry) {
		groupID, entryValue, err := h.identityFromStatusOrDeps(ctx, ipaccesslistentry)
		if err != nil {
			return result.Error(currentState, err)
		}
		// Check if the entry exists in Atlas
		_, _, err = h.atlasClient.ProjectIPAccessListApi.GetAccessListEntry(ctx, groupID, entryValue).Execute()
		if v20250312sdk.IsErrorCode(err, atlasAccessListNotFound) || v20250312sdk.IsErrorCode(err, atlasAccessListEntryNotFound) {
			// Atlas deleted the entry (deleteAfterDate elapsed).
			return expiredResult(currentState), nil
		}
		if err != nil {
			return result.Error(currentState, fmt.Errorf("failed to verify IP access list entry in Atlas: %w", err))
		}
		// Entry is still active in Atlas. Skip the normal ShouldUpdate/update cycle:
		// entries with deleteAfterDate are ephemeral and re-creating them after the
		// date has passed would produce EXPIRATION_DATE_IN_PAST errors. Just schedule
		// a requeue so we poll Atlas again soon.
		return ctrlstate.Result{
			Result:    reconcile.Result{RequeueAfter: requeueAfterExpiry(ipaccesslistentry)},
			NextState: currentState,
			StateMsg:  "IP access list entry is active.",
		}, nil
	}

	update, err := ctrlstate.ShouldUpdate(ipaccesslistentry, deps...)
	if err != nil {
		return result.Error(currentState, reconcile.TerminalError(err))
	}

	if !update {
		return result.NextState(currentState, "IP access list entry is up to date. No update required.")
	}

	groupID, oldEntryValue, err := h.identityFromStatusOrDeps(ctx, ipaccesslistentry)
	if err != nil {
		return result.Error(currentState, err)
	}

	newEntryValue := entryValueFromSpec(ipaccesslistentry)
	if newEntryValue == "" {
		return result.Error(currentState, errors.New("spec must set exactly one of: ipAddress, cidrBlock, awsSecurityGroup"))
	}

	// Delete the old one if the entry value changed (e.g. CIDR, IP, AWS Security Group were modified).
	if oldEntryValue != "" && oldEntryValue != newEntryValue {
		_, err = h.atlasClient.ProjectIPAccessListApi.DeleteAccessListEntry(ctx, groupID, oldEntryValue).Execute()
		if err != nil && !v20250312sdk.IsErrorCode(err, atlasAccessListNotFound) && !v20250312sdk.IsErrorCode(err, atlasAccessListEntryNotFound) {
			return result.Error(currentState, fmt.Errorf("failed to delete old IP access list entry %q: %w", oldEntryValue, err))
		}
	}

	entry := buildNetworkPermissionEntry(ipaccesslistentry)
	_, _, err = h.atlasClient.ProjectIPAccessListApi.CreateAccessListEntry(ctx, groupID, &[]v20250312sdk.NetworkPermissionEntry{entry}).Execute()
	if err != nil {
		return result.Error(currentState, fmt.Errorf("failed to create IP access list entry %q: %w", newEntryValue, err))
	}

	if err := h.persistIdentity(ctx, ipaccesslistentry, groupID, newEntryValue); err != nil {
		return result.Error(currentState, err)
	}

	copy := ipaccesslistentry.DeepCopy()
	if err := ctrlstate.NewPatcher(copy).UpdateStateTracker(deps...).Patch(ctx, h.kubeClient); err != nil {
		return result.Error(currentState, fmt.Errorf("failed to update state tracker: %w", err))
	}

	return result.NextState(state.StateUpdating, "IP access list entry updated. Waiting for activation.")
}

// resolveIdentity resolves the groupID and entry value from spec + dependencies.
func (h *Handlerv20250312) resolveIdentity(ctx context.Context, ipaccesslistentry *akov2generated.IPAccessListEntry) (groupID, entryValue string, err error) {
	groupID, err = h.resolveGroupID(ctx, ipaccesslistentry)
	if err != nil {
		return "", "", err
	}

	entryValue = entryValueFromSpec(ipaccesslistentry)
	if entryValue == "" {
		return "", "", errors.New("spec must set exactly one of: ipAddress, cidrBlock, awsSecurityGroup")
	}

	return groupID, entryValue, nil
}

// resolveGroupID resolves the Atlas project ID from spec.v20250312.groupId or spec.v20250312.groupRef.
func (h *Handlerv20250312) resolveGroupID(ctx context.Context, ipaccesslistentry *akov2generated.IPAccessListEntry) (string, error) {
	if ipaccesslistentry.Spec.V20250312 == nil {
		return "", errors.New("spec.v20250312 is not set")
	}

	if ipaccesslistentry.Spec.V20250312.GroupId != nil {
		return *ipaccesslistentry.Spec.V20250312.GroupId, nil
	}

	if ipaccesslistentry.Spec.V20250312.GroupRef != nil {
		deps, err := h.getDependencies(ctx, ipaccesslistentry)
		if err != nil {
			return "", fmt.Errorf("failed to resolve dependencies: %w", err)
		}

		for _, dep := range deps {
			if g, ok := dep.(*akov2generated.Group); ok {
				if g.Status.V20250312 != nil && g.Status.V20250312.Id != nil {
					return *g.Status.V20250312.Id, nil
				}
			}
		}
	}

	return "", errors.New("could not resolve groupId: neither groupId nor groupRef with a ready Group is available")
}

func (h *Handlerv20250312) identityFromStatusOrDeps(ctx context.Context, ipaccesslistentry *akov2generated.IPAccessListEntry) (groupID, entryValue string, err error) {
	if groupID, entryValue, ok := readIdentityFromStatus(ipaccesslistentry); ok {
		return groupID, entryValue, nil
	}

	return h.resolveIdentity(ctx, ipaccesslistentry)
}

func readIdentityFromStatus(ipaccesslistentry *akov2generated.IPAccessListEntry) (groupID, entryValue string, ok bool) {
	if ipaccesslistentry.Status.V20250312 == nil {
		return "", "", false
	}

	s := ipaccesslistentry.Status.V20250312
	if s.GroupId == nil {
		return "", "", false
	}

	groupID = *s.GroupId

	switch {
	case s.CidrBlock != nil && *s.CidrBlock != "":
		entryValue = *s.CidrBlock
	case s.IpAddress != nil && *s.IpAddress != "":
		entryValue = *s.IpAddress
	case s.AwsSecurityGroup != nil && *s.AwsSecurityGroup != "":
		entryValue = *s.AwsSecurityGroup
	}

	if entryValue == "" {
		return "", "", false
	}

	return groupID, entryValue, true
}

func (h *Handlerv20250312) persistIdentity(ctx context.Context, ipaccesslistentry *akov2generated.IPAccessListEntry, groupID, entryValue string) error {
	copy := ipaccesslistentry.DeepCopy()
	if copy.Status.V20250312 == nil {
		copy.Status.V20250312 = &akov2generated.IPAccessListEntryStatusV20250312{}
	}

	copy.Status.V20250312.GroupId = &groupID

	if strings.Contains(entryValue, "/") {
		copy.Status.V20250312.CidrBlock = &entryValue
	} else if strings.HasPrefix(entryValue, "sg-") {
		copy.Status.V20250312.AwsSecurityGroup = &entryValue
	} else {
		copy.Status.V20250312.IpAddress = &entryValue
	}

	if err := ctrlstate.NewPatcher(copy).UpdateStatus().Patch(ctx, h.kubeClient); err != nil {
		return fmt.Errorf("failed to patch IPAccessListEntry status: %w", err)
	}

	return nil
}

// getEntryStatus returns the Atlas activation status for the entry.
func (h *Handlerv20250312) getEntryStatus(ctx context.Context, groupID, entryValue string) (string, error) {
	resp, _, err := h.atlasClient.ProjectIPAccessListApi.GetAccessListStatus(ctx, groupID, entryValue).Execute()
	if err != nil {
		return "", err
	}

	return resp.GetSTATUS(), nil
}

func entryValueFromSpec(ipaccesslistentry *akov2generated.IPAccessListEntry) string {
	if ipaccesslistentry.Spec.V20250312 == nil || ipaccesslistentry.Spec.V20250312.Entry == nil {
		return ""
	}

	entry := ipaccesslistentry.Spec.V20250312.Entry

	switch {
	case entry.CidrBlock != nil && *entry.CidrBlock != "":
		return normalizeCIDR(*entry.CidrBlock)
	case entry.IpAddress != nil && *entry.IpAddress != "":
		return *entry.IpAddress
	case entry.AwsSecurityGroup != nil && *entry.AwsSecurityGroup != "":
		return *entry.AwsSecurityGroup
	}

	return ""
}

func buildNetworkPermissionEntry(ipaccesslistentry *akov2generated.IPAccessListEntry) v20250312sdk.NetworkPermissionEntry {
	entry := v20250312sdk.NetworkPermissionEntry{}

	if ipaccesslistentry.Spec.V20250312 == nil || ipaccesslistentry.Spec.V20250312.Entry == nil {
		return entry
	}

	e := ipaccesslistentry.Spec.V20250312.Entry

	if e.CidrBlock != nil {
		normalized := normalizeCIDR(*e.CidrBlock)
		entry.CidrBlock = &normalized
	} else if e.IpAddress != nil {
		entry.IpAddress = e.IpAddress
	} else if e.AwsSecurityGroup != nil {
		entry.AwsSecurityGroup = e.AwsSecurityGroup
	}

	if e.Comment != nil {
		entry.Comment = e.Comment
	}

	if e.DeleteAfterDate != nil {
		t, err := time.Parse(time.RFC3339, *e.DeleteAfterDate)
		if err == nil {
			entry.DeleteAfterDate = &t
		}
	}

	return entry
}

// hasDeleteAfterDate returns true if the entry spec contains a deleteAfterDate field.
func hasDeleteAfterDate(ipaccesslistentry *akov2generated.IPAccessListEntry) bool {
	return ipaccesslistentry.Spec.V20250312 != nil &&
		ipaccesslistentry.Spec.V20250312.Entry != nil &&
		ipaccesslistentry.Spec.V20250312.Entry.DeleteAfterDate != nil
}

// requeueAfterExpiry returns the duration to wait before the next Atlas GET check.
// When the deleteAfterDate is still in the future, schedules a requeue shortly after
// the date so expiry is detected promptly. Once the date has passed, polls every
// pollAfterExpiry seconds until Atlas confirms deletion via 404.
// WARNING: This is best effort for purely informational purposes, it is not guaranteed to be accurate.
func requeueAfterExpiry(ipaccesslistentry *akov2generated.IPAccessListEntry) time.Duration {
	const expiryBuffer = 30 * time.Second
	const pollAfterExpiry = 30 * time.Second
	const minRequeue = 30 * time.Second
	if ipaccesslistentry.Spec.V20250312 != nil &&
		ipaccesslistentry.Spec.V20250312.Entry != nil &&
		ipaccesslistentry.Spec.V20250312.Entry.DeleteAfterDate != nil {
		if t, err := time.Parse(time.RFC3339, *ipaccesslistentry.Spec.V20250312.Entry.DeleteAfterDate); err == nil {
			remaining := time.Until(t)
			if remaining <= 0 {
				// Date has already passed — poll frequently until Atlas confirms 404.
				return pollAfterExpiry
			}
			if d := remaining + expiryBuffer; d > minRequeue {
				return d
			}
		}
	}
	return minRequeue
}

// checkExpiry returns true if the entry was already marked expired in a previous reconcile.
// This is a fast path that avoids an Atlas API call when the entry is known to be gone.
// The authoritative expiry detection happens in handleSteadyState via Atlas GET.
func checkExpiry(ipaccesslistentry *akov2generated.IPAccessListEntry) bool {
	stateCondition := k8smeta.FindStatusCondition(ipaccesslistentry.GetConditions(), "State")
	return stateCondition != nil && stateCondition.Message == expiredStateMsg
}

// expiredResult returns a reconcile result that marks the resource as expired without
// transitioning to a new state. The 24 h requeue prevents reconcileReapply from running
// while still allowing the controller to notice if the resource is updated.
func expiredResult(currentState state.ResourceState) ctrlstate.Result {
	return ctrlstate.Result{
		Result:    reconcile.Result{RequeueAfter: 24 * time.Hour},
		NextState: currentState,
		StateMsg:  expiredStateMsg,
	}
}

func normalizeCIDR(cidr string) string {
	_, network, err := net.ParseCIDR(cidr)
	if err != nil {
		return cidr
	}

	return network.String()
}
