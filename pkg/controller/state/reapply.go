// Copyright 2025 MongoDB Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package state

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func ShouldReapply(obj metav1.Object) (bool, error) {
	timestamp, hasTimestamp, err := ReapplyTimestamp(obj)
	if err != nil {
		return false, err
	}

	if !hasTimestamp {
		return false, nil
	}

	period, hasPeriod, err := ReapplyPeriod(obj)
	if err != nil {
		return false, err
	}

	if !hasPeriod {
		return false, nil
	}

	diff := time.Until(timestamp.Add(period))

	return diff <= 0, nil
}

func ReapplyPeriod(obj metav1.Object) (time.Duration, bool, error) {
	annotationPeriod, ok := obj.GetAnnotations()["mongodb.com/reapply-period"]
	if !ok {
		return 0, false, nil
	}

	period, err := time.ParseDuration(annotationPeriod)
	if err != nil {
		return 0, false, fmt.Errorf("failed to parse reapply period: %w", err)
	}

	if period < 60*time.Second {
		return 0, false, errors.New("reapply period is invalid: must be greater than 60m")
	}

	return period, true, nil
}

func ReapplyTimestamp(obj metav1.Object) (time.Time, bool, error) {
	annotationTimestamp, ok := obj.GetAnnotations()[AnnotationReapplyTimestamp]
	if !ok {
		return time.Time{}, false, nil
	}

	timestampMillis, err := strconv.ParseInt(annotationTimestamp, 10, 0)
	if err != nil {
		return time.Time{}, false, fmt.Errorf("failed to parse reapply timestamp: %w", err)
	}

	return time.UnixMilli(timestampMillis), true, nil
}

func PatchReapplyTimestamp(ctx context.Context, kubeClient client.Client, obj client.Object) (time.Duration, error) {
	period, hasPeriod, err := ReapplyPeriod(obj)
	if err != nil {
		return 0, err
	}

	if !hasPeriod {
		return 0, nil
	}

	timestamp, hasTimestamp, err := ReapplyTimestamp(obj)
	if err != nil {
		return 0, err
	}

	now := time.Now()
	diff := timestamp.Add(period).Sub(now)
	if hasTimestamp && diff > 0 {
		return diff, nil
	}

	patch := fmt.Appendf(nil, `[{
	"op":    "replace",
	"path":  "/metadata/annotations/%v",
	"value": "%v"
}]`, jsonPatchReplacer.Replace(AnnotationReapplyTimestamp), now.UnixMilli())

	if err := kubeClient.Patch(ctx, obj, client.RawPatch(types.JSONPatchType, patch)); err != nil {
		return 0, fmt.Errorf("failed to patch object: %w", err)
	}

	return period, nil
}
