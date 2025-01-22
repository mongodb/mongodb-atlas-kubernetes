//go:build !ignore_autogenerated

/*
Copyright (C) MongoDB, Inc. 2020-present.

Licensed under the Apache License, Version 2.0 (the "License"); you may
not use this file except in compliance with the License. You may obtain
a copy of the License at http://www.apache.org/licenses/LICENSE-2.0
*/

// Code generated by controller-gen. DO NOT EDIT.

package status

import (
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/authmode"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/common"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/project"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AWSStatus) DeepCopyInto(out *AWSStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AWSStatus.
func (in *AWSStatus) DeepCopy() *AWSStatus {
	if in == nil {
		return nil
	}
	out := new(AWSStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AlertConfiguration) DeepCopyInto(out *AlertConfiguration) {
	*out = *in
	if in.Enabled != nil {
		in, out := &in.Enabled, &out.Enabled
		*out = new(bool)
		**out = **in
	}
	if in.CurrentValue != nil {
		in, out := &in.CurrentValue, &out.CurrentValue
		*out = new(CurrentValue)
		**out = **in
	}
	if in.Matchers != nil {
		in, out := &in.Matchers, &out.Matchers
		*out = make([]Matcher, len(*in))
		copy(*out, *in)
	}
	if in.MetricThreshold != nil {
		in, out := &in.MetricThreshold, &out.MetricThreshold
		*out = new(MetricThreshold)
		**out = **in
	}
	if in.Threshold != nil {
		in, out := &in.Threshold, &out.Threshold
		*out = new(Threshold)
		**out = **in
	}
	if in.Notifications != nil {
		in, out := &in.Notifications, &out.Notifications
		*out = make([]Notification, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AlertConfiguration.
func (in *AlertConfiguration) DeepCopy() *AlertConfiguration {
	if in == nil {
		return nil
	}
	out := new(AlertConfiguration)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AtlasCustomRoleStatus) DeepCopyInto(out *AtlasCustomRoleStatus) {
	*out = *in
	in.Common.DeepCopyInto(&out.Common)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AtlasCustomRoleStatus.
func (in *AtlasCustomRoleStatus) DeepCopy() *AtlasCustomRoleStatus {
	if in == nil {
		return nil
	}
	out := new(AtlasCustomRoleStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AtlasDatabaseUserStatus) DeepCopyInto(out *AtlasDatabaseUserStatus) {
	*out = *in
	in.Common.DeepCopyInto(&out.Common)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AtlasDatabaseUserStatus.
func (in *AtlasDatabaseUserStatus) DeepCopy() *AtlasDatabaseUserStatus {
	if in == nil {
		return nil
	}
	out := new(AtlasDatabaseUserStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AtlasDeploymentStatus) DeepCopyInto(out *AtlasDeploymentStatus) {
	*out = *in
	in.Common.DeepCopyInto(&out.Common)
	if in.ConnectionStrings != nil {
		in, out := &in.ConnectionStrings, &out.ConnectionStrings
		*out = new(ConnectionStrings)
		(*in).DeepCopyInto(*out)
	}
	if in.ReplicaSets != nil {
		in, out := &in.ReplicaSets, &out.ReplicaSets
		*out = make([]ReplicaSet, len(*in))
		copy(*out, *in)
	}
	if in.ServerlessPrivateEndpoints != nil {
		in, out := &in.ServerlessPrivateEndpoints, &out.ServerlessPrivateEndpoints
		*out = make([]ServerlessPrivateEndpoint, len(*in))
		copy(*out, *in)
	}
	if in.CustomZoneMapping != nil {
		in, out := &in.CustomZoneMapping, &out.CustomZoneMapping
		*out = new(CustomZoneMapping)
		(*in).DeepCopyInto(*out)
	}
	if in.ManagedNamespaces != nil {
		in, out := &in.ManagedNamespaces, &out.ManagedNamespaces
		*out = make([]ManagedNamespace, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.SearchIndexes != nil {
		in, out := &in.SearchIndexes, &out.SearchIndexes
		*out = make([]DeploymentSearchIndexStatus, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AtlasDeploymentStatus.
func (in *AtlasDeploymentStatus) DeepCopy() *AtlasDeploymentStatus {
	if in == nil {
		return nil
	}
	out := new(AtlasDeploymentStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AtlasFederatedAuthStatus) DeepCopyInto(out *AtlasFederatedAuthStatus) {
	*out = *in
	in.Common.DeepCopyInto(&out.Common)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AtlasFederatedAuthStatus.
func (in *AtlasFederatedAuthStatus) DeepCopy() *AtlasFederatedAuthStatus {
	if in == nil {
		return nil
	}
	out := new(AtlasFederatedAuthStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AtlasIPAccessListStatus) DeepCopyInto(out *AtlasIPAccessListStatus) {
	*out = *in
	in.Common.DeepCopyInto(&out.Common)
	if in.Entries != nil {
		in, out := &in.Entries, &out.Entries
		*out = make([]IPAccessEntryStatus, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AtlasIPAccessListStatus.
func (in *AtlasIPAccessListStatus) DeepCopy() *AtlasIPAccessListStatus {
	if in == nil {
		return nil
	}
	out := new(AtlasIPAccessListStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AtlasNetworkPeer) DeepCopyInto(out *AtlasNetworkPeer) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AtlasNetworkPeer.
func (in *AtlasNetworkPeer) DeepCopy() *AtlasNetworkPeer {
	if in == nil {
		return nil
	}
	out := new(AtlasNetworkPeer)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AtlasNetworkPeeringStatus) DeepCopyInto(out *AtlasNetworkPeeringStatus) {
	*out = *in
	in.Common.DeepCopyInto(&out.Common)
	if in.AWSStatus != nil {
		in, out := &in.AWSStatus, &out.AWSStatus
		*out = new(AWSStatus)
		**out = **in
	}
	if in.AzureStatus != nil {
		in, out := &in.AzureStatus, &out.AzureStatus
		*out = new(AzureStatus)
		**out = **in
	}
	if in.GoogleStatus != nil {
		in, out := &in.GoogleStatus, &out.GoogleStatus
		*out = new(GoogleStatus)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AtlasNetworkPeeringStatus.
func (in *AtlasNetworkPeeringStatus) DeepCopy() *AtlasNetworkPeeringStatus {
	if in == nil {
		return nil
	}
	out := new(AtlasNetworkPeeringStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AtlasPrivateEndpointStatus) DeepCopyInto(out *AtlasPrivateEndpointStatus) {
	*out = *in
	in.Common.DeepCopyInto(&out.Common)
	if in.ServiceAttachmentNames != nil {
		in, out := &in.ServiceAttachmentNames, &out.ServiceAttachmentNames
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.Endpoints != nil {
		in, out := &in.Endpoints, &out.Endpoints
		*out = make([]EndpointInterfaceStatus, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AtlasPrivateEndpointStatus.
func (in *AtlasPrivateEndpointStatus) DeepCopy() *AtlasPrivateEndpointStatus {
	if in == nil {
		return nil
	}
	out := new(AtlasPrivateEndpointStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AtlasProjectStatus) DeepCopyInto(out *AtlasProjectStatus) {
	*out = *in
	in.Common.DeepCopyInto(&out.Common)
	if in.ExpiredIPAccessList != nil {
		in, out := &in.ExpiredIPAccessList, &out.ExpiredIPAccessList
		*out = make([]project.IPAccessList, len(*in))
		copy(*out, *in)
	}
	if in.PrivateEndpoints != nil {
		in, out := &in.PrivateEndpoints, &out.PrivateEndpoints
		*out = make([]ProjectPrivateEndpoint, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.NetworkPeers != nil {
		in, out := &in.NetworkPeers, &out.NetworkPeers
		*out = make([]AtlasNetworkPeer, len(*in))
		copy(*out, *in)
	}
	if in.AuthModes != nil {
		in, out := &in.AuthModes, &out.AuthModes
		*out = make(authmode.AuthModes, len(*in))
		copy(*out, *in)
	}
	if in.AlertConfigurations != nil {
		in, out := &in.AlertConfigurations, &out.AlertConfigurations
		*out = make([]AlertConfiguration, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.CloudProviderIntegrations != nil {
		in, out := &in.CloudProviderIntegrations, &out.CloudProviderIntegrations
		*out = make([]CloudProviderIntegration, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.CustomRoles != nil {
		in, out := &in.CustomRoles, &out.CustomRoles
		*out = make([]CustomRole, len(*in))
		copy(*out, *in)
	}
	if in.Teams != nil {
		in, out := &in.Teams, &out.Teams
		*out = make([]ProjectTeamStatus, len(*in))
		copy(*out, *in)
	}
	if in.Prometheus != nil {
		in, out := &in.Prometheus, &out.Prometheus
		*out = new(Prometheus)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AtlasProjectStatus.
func (in *AtlasProjectStatus) DeepCopy() *AtlasProjectStatus {
	if in == nil {
		return nil
	}
	out := new(AtlasProjectStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AtlasSearchIndexConfigStatus) DeepCopyInto(out *AtlasSearchIndexConfigStatus) {
	*out = *in
	in.Common.DeepCopyInto(&out.Common)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AtlasSearchIndexConfigStatus.
func (in *AtlasSearchIndexConfigStatus) DeepCopy() *AtlasSearchIndexConfigStatus {
	if in == nil {
		return nil
	}
	out := new(AtlasSearchIndexConfigStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AtlasStreamConnectionStatus) DeepCopyInto(out *AtlasStreamConnectionStatus) {
	*out = *in
	in.Common.DeepCopyInto(&out.Common)
	if in.Instances != nil {
		in, out := &in.Instances, &out.Instances
		*out = make([]common.ResourceRefNamespaced, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AtlasStreamConnectionStatus.
func (in *AtlasStreamConnectionStatus) DeepCopy() *AtlasStreamConnectionStatus {
	if in == nil {
		return nil
	}
	out := new(AtlasStreamConnectionStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AtlasStreamInstanceStatus) DeepCopyInto(out *AtlasStreamInstanceStatus) {
	*out = *in
	in.Common.DeepCopyInto(&out.Common)
	if in.Hostnames != nil {
		in, out := &in.Hostnames, &out.Hostnames
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.Connections != nil {
		in, out := &in.Connections, &out.Connections
		*out = make([]StreamConnection, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AtlasStreamInstanceStatus.
func (in *AtlasStreamInstanceStatus) DeepCopy() *AtlasStreamInstanceStatus {
	if in == nil {
		return nil
	}
	out := new(AtlasStreamInstanceStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AzureStatus) DeepCopyInto(out *AzureStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AzureStatus.
func (in *AzureStatus) DeepCopy() *AzureStatus {
	if in == nil {
		return nil
	}
	out := new(AzureStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BackupCompliancePolicyStatus) DeepCopyInto(out *BackupCompliancePolicyStatus) {
	*out = *in
	in.Common.DeepCopyInto(&out.Common)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BackupCompliancePolicyStatus.
func (in *BackupCompliancePolicyStatus) DeepCopy() *BackupCompliancePolicyStatus {
	if in == nil {
		return nil
	}
	out := new(BackupCompliancePolicyStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BackupPolicyStatus) DeepCopyInto(out *BackupPolicyStatus) {
	*out = *in
	in.Common.DeepCopyInto(&out.Common)
	if in.BackupScheduleIDs != nil {
		in, out := &in.BackupScheduleIDs, &out.BackupScheduleIDs
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BackupPolicyStatus.
func (in *BackupPolicyStatus) DeepCopy() *BackupPolicyStatus {
	if in == nil {
		return nil
	}
	out := new(BackupPolicyStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BackupScheduleStatus) DeepCopyInto(out *BackupScheduleStatus) {
	*out = *in
	in.Common.DeepCopyInto(&out.Common)
	if in.DeploymentIDs != nil {
		in, out := &in.DeploymentIDs, &out.DeploymentIDs
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BackupScheduleStatus.
func (in *BackupScheduleStatus) DeepCopy() *BackupScheduleStatus {
	if in == nil {
		return nil
	}
	out := new(BackupScheduleStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CloudProviderIntegration) DeepCopyInto(out *CloudProviderIntegration) {
	*out = *in
	if in.FeatureUsages != nil {
		in, out := &in.FeatureUsages, &out.FeatureUsages
		*out = make([]FeatureUsage, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CloudProviderIntegration.
func (in *CloudProviderIntegration) DeepCopy() *CloudProviderIntegration {
	if in == nil {
		return nil
	}
	out := new(CloudProviderIntegration)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ConnectionStrings) DeepCopyInto(out *ConnectionStrings) {
	*out = *in
	if in.PrivateEndpoint != nil {
		in, out := &in.PrivateEndpoint, &out.PrivateEndpoint
		*out = make([]PrivateEndpoint, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ConnectionStrings.
func (in *ConnectionStrings) DeepCopy() *ConnectionStrings {
	if in == nil {
		return nil
	}
	out := new(ConnectionStrings)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CurrentValue) DeepCopyInto(out *CurrentValue) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CurrentValue.
func (in *CurrentValue) DeepCopy() *CurrentValue {
	if in == nil {
		return nil
	}
	out := new(CurrentValue)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CustomRole) DeepCopyInto(out *CustomRole) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CustomRole.
func (in *CustomRole) DeepCopy() *CustomRole {
	if in == nil {
		return nil
	}
	out := new(CustomRole)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CustomZoneMapping) DeepCopyInto(out *CustomZoneMapping) {
	*out = *in
	if in.CustomZoneMapping != nil {
		in, out := &in.CustomZoneMapping, &out.CustomZoneMapping
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CustomZoneMapping.
func (in *CustomZoneMapping) DeepCopy() *CustomZoneMapping {
	if in == nil {
		return nil
	}
	out := new(CustomZoneMapping)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DataFederationStatus) DeepCopyInto(out *DataFederationStatus) {
	*out = *in
	in.Common.DeepCopyInto(&out.Common)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DataFederationStatus.
func (in *DataFederationStatus) DeepCopy() *DataFederationStatus {
	if in == nil {
		return nil
	}
	out := new(DataFederationStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DeploymentSearchIndexStatus) DeepCopyInto(out *DeploymentSearchIndexStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DeploymentSearchIndexStatus.
func (in *DeploymentSearchIndexStatus) DeepCopy() *DeploymentSearchIndexStatus {
	if in == nil {
		return nil
	}
	out := new(DeploymentSearchIndexStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Endpoint) DeepCopyInto(out *Endpoint) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Endpoint.
func (in *Endpoint) DeepCopy() *Endpoint {
	if in == nil {
		return nil
	}
	out := new(Endpoint)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EndpointInterfaceStatus) DeepCopyInto(out *EndpointInterfaceStatus) {
	*out = *in
	if in.GCPForwardingRules != nil {
		in, out := &in.GCPForwardingRules, &out.GCPForwardingRules
		*out = make([]GCPForwardingRule, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EndpointInterfaceStatus.
func (in *EndpointInterfaceStatus) DeepCopy() *EndpointInterfaceStatus {
	if in == nil {
		return nil
	}
	out := new(EndpointInterfaceStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FeatureUsage) DeepCopyInto(out *FeatureUsage) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FeatureUsage.
func (in *FeatureUsage) DeepCopy() *FeatureUsage {
	if in == nil {
		return nil
	}
	out := new(FeatureUsage)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GCPEndpoint) DeepCopyInto(out *GCPEndpoint) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GCPEndpoint.
func (in *GCPEndpoint) DeepCopy() *GCPEndpoint {
	if in == nil {
		return nil
	}
	out := new(GCPEndpoint)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GCPForwardingRule) DeepCopyInto(out *GCPForwardingRule) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GCPForwardingRule.
func (in *GCPForwardingRule) DeepCopy() *GCPForwardingRule {
	if in == nil {
		return nil
	}
	out := new(GCPForwardingRule)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GoogleStatus) DeepCopyInto(out *GoogleStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GoogleStatus.
func (in *GoogleStatus) DeepCopy() *GoogleStatus {
	if in == nil {
		return nil
	}
	out := new(GoogleStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *IPAccessEntryStatus) DeepCopyInto(out *IPAccessEntryStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new IPAccessEntryStatus.
func (in *IPAccessEntryStatus) DeepCopy() *IPAccessEntryStatus {
	if in == nil {
		return nil
	}
	out := new(IPAccessEntryStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ManagedNamespace) DeepCopyInto(out *ManagedNamespace) {
	*out = *in
	if in.IsCustomShardKeyHashed != nil {
		in, out := &in.IsCustomShardKeyHashed, &out.IsCustomShardKeyHashed
		*out = new(bool)
		**out = **in
	}
	if in.IsShardKeyUnique != nil {
		in, out := &in.IsShardKeyUnique, &out.IsShardKeyUnique
		*out = new(bool)
		**out = **in
	}
	if in.PresplitHashedZones != nil {
		in, out := &in.PresplitHashedZones, &out.PresplitHashedZones
		*out = new(bool)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ManagedNamespace.
func (in *ManagedNamespace) DeepCopy() *ManagedNamespace {
	if in == nil {
		return nil
	}
	out := new(ManagedNamespace)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Matcher) DeepCopyInto(out *Matcher) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Matcher.
func (in *Matcher) DeepCopy() *Matcher {
	if in == nil {
		return nil
	}
	out := new(Matcher)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MetricThreshold) DeepCopyInto(out *MetricThreshold) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MetricThreshold.
func (in *MetricThreshold) DeepCopy() *MetricThreshold {
	if in == nil {
		return nil
	}
	out := new(MetricThreshold)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Notification) DeepCopyInto(out *Notification) {
	*out = *in
	if in.DelayMin != nil {
		in, out := &in.DelayMin, &out.DelayMin
		*out = new(int)
		**out = **in
	}
	if in.EmailEnabled != nil {
		in, out := &in.EmailEnabled, &out.EmailEnabled
		*out = new(bool)
		**out = **in
	}
	if in.SMSEnabled != nil {
		in, out := &in.SMSEnabled, &out.SMSEnabled
		*out = new(bool)
		**out = **in
	}
	if in.Roles != nil {
		in, out := &in.Roles, &out.Roles
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Notification.
func (in *Notification) DeepCopy() *Notification {
	if in == nil {
		return nil
	}
	out := new(Notification)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PrivateEndpoint) DeepCopyInto(out *PrivateEndpoint) {
	*out = *in
	if in.Endpoints != nil {
		in, out := &in.Endpoints, &out.Endpoints
		*out = make([]Endpoint, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PrivateEndpoint.
func (in *PrivateEndpoint) DeepCopy() *PrivateEndpoint {
	if in == nil {
		return nil
	}
	out := new(PrivateEndpoint)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ProjectPrivateEndpoint) DeepCopyInto(out *ProjectPrivateEndpoint) {
	*out = *in
	if in.ServiceAttachmentNames != nil {
		in, out := &in.ServiceAttachmentNames, &out.ServiceAttachmentNames
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.Endpoints != nil {
		in, out := &in.Endpoints, &out.Endpoints
		*out = make([]GCPEndpoint, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ProjectPrivateEndpoint.
func (in *ProjectPrivateEndpoint) DeepCopy() *ProjectPrivateEndpoint {
	if in == nil {
		return nil
	}
	out := new(ProjectPrivateEndpoint)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ProjectTeamStatus) DeepCopyInto(out *ProjectTeamStatus) {
	*out = *in
	out.TeamRef = in.TeamRef
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ProjectTeamStatus.
func (in *ProjectTeamStatus) DeepCopy() *ProjectTeamStatus {
	if in == nil {
		return nil
	}
	out := new(ProjectTeamStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Prometheus) DeepCopyInto(out *Prometheus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Prometheus.
func (in *Prometheus) DeepCopy() *Prometheus {
	if in == nil {
		return nil
	}
	out := new(Prometheus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ReplicaSet) DeepCopyInto(out *ReplicaSet) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ReplicaSet.
func (in *ReplicaSet) DeepCopy() *ReplicaSet {
	if in == nil {
		return nil
	}
	out := new(ReplicaSet)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ServerlessPrivateEndpoint) DeepCopyInto(out *ServerlessPrivateEndpoint) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ServerlessPrivateEndpoint.
func (in *ServerlessPrivateEndpoint) DeepCopy() *ServerlessPrivateEndpoint {
	if in == nil {
		return nil
	}
	out := new(ServerlessPrivateEndpoint)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *StreamConnection) DeepCopyInto(out *StreamConnection) {
	*out = *in
	out.ResourceRef = in.ResourceRef
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new StreamConnection.
func (in *StreamConnection) DeepCopy() *StreamConnection {
	if in == nil {
		return nil
	}
	out := new(StreamConnection)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TeamProject) DeepCopyInto(out *TeamProject) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TeamProject.
func (in *TeamProject) DeepCopy() *TeamProject {
	if in == nil {
		return nil
	}
	out := new(TeamProject)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TeamStatus) DeepCopyInto(out *TeamStatus) {
	*out = *in
	in.Common.DeepCopyInto(&out.Common)
	if in.Projects != nil {
		in, out := &in.Projects, &out.Projects
		*out = make([]TeamProject, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TeamStatus.
func (in *TeamStatus) DeepCopy() *TeamStatus {
	if in == nil {
		return nil
	}
	out := new(TeamStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Threshold) DeepCopyInto(out *Threshold) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Threshold.
func (in *Threshold) DeepCopy() *Threshold {
	if in == nil {
		return nil
	}
	out := new(Threshold)
	in.DeepCopyInto(out)
	return out
}
