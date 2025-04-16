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

package atlas

const (
	// Error codes that Atlas may return that we are concerned about
	GroupExistsAPIErrorCode = "GROUP_ALREADY_EXISTS"

	// The error that Atlas API returns if the GET request is sent to read the project that either doesn't exist
	// or the user doesn't have permissions for
	NotInGroup = "NOT_IN_GROUP"

	// Error indicates that the project is being removed while it still has clusters
	CannotCloseGroupActiveAtlasDeployment = "CANNOT_CLOSE_GROUP_ACTIVE_ATLAS_CLUSTERS"

	// Error indicates that the database user doesn't exist
	UsernameNotFound = "USERNAME_NOT_FOUND"

	// Error indicates that the database user doesn't exist
	UserNotfound = "USER_NOT_FOUND"

	// Error indicates that the cluster doesn't exist
	ClusterNotFound = "CLUSTER_NOT_FOUND"

	// ServerlessClusterNotFound indicates that the serverless cluster doesn't exist
	ServerlessInstanceNotFound = "SERVERLESS_INSTANCE_NOT_FOUND"

	// ServerlessClusterFromClusterAPI indicates that we are trying to access
	// a serverless instance from the cluster API, which is not allowed
	ServerlessInstanceFromClusterAPI = "CANNOT_USE_SERVERLESS_INSTANCE_IN_CLUSTER_API"

	ClusterInstanceFromServerlessAPI = "CANNOT_USE_CLUSTER_IN_SERVERLESS_INSTANCE_API"

	// Resource not found
	ResourceNotFound = "RESOURCE_NOT_FOUND"

	// Instance for the passed {groupId, tenantName} pair does not exist
	DataFederationTenantNotFound = "DATA_FEDERATION_TENANT_NOT_FOUND_FOR_NAME"

	// Backup Compliance Policy rejected, as there are existing backup policies which do not meet the requirements
	BackupComplianceNotMet = "BACKUP_POLICIES_NOT_MEETING_BACKUP_COMPLIANCE_POLICY_REQUIREMENTS"

	ProviderUnsupported = "PROVIDER_UNSUPPORTED"

	// Cannot use the Flex API to interact with non-Flex clusters
	NonFlexInFlexAPI = "CANNOT_USE_NON_FLEX_CLUSTER_IN_FLEX_API"

	FeatureUnsupported = "FEATURE_UNSUPPORTED"
)
