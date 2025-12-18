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

package testparams

import "strings"

// TestParams holds all test parameters for test isolation purposes.
// Shared values (OrgID, OperatorNamespace, CredentialsSecretName) are typically
// set from input config and remain constant across tests. Per-test values
// (GroupID, GroupName) are set per test case.
type TestParams struct {
	// GroupID is the Atlas group ID, assigned by Atlas after Group creation.
	// This is per-test and may be empty initially.
	GroupID string
	// OrgID is the Atlas organization ID, set from input config (e.g., MCLI_ORG_ID env var).
	OrgID string
	// GroupName is a randomized name for test isolation, per test case.
	GroupName string
	// OperatorNamespace is the namespace where the operator is running, set from input config.
	OperatorNamespace string
	// CredentialsSecretName is the name of the credentials secret, set from input config.
	CredentialsSecretName string
}

// New creates a new TestParams struct with shared configuration values.
// These values are typically constant across all tests in a suite:
//   - orgID: Atlas organization ID (from MCLI_ORG_ID env var)
//   - operatorNamespace: Namespace where operator is running (from OPERATOR_NAMESPACE env var)
//   - credentialsSecretName: Name of the credentials secret (e.g., DefaultGlobalCredentials)
//
// Per-test values (GroupID, GroupName) should be set using WithGroupID and WithGroupName.
func New(orgID, operatorNamespace, credentialsSecretName string) *TestParams {
	return &TestParams{
		OrgID:                 orgID,
		OperatorNamespace:     operatorNamespace,
		CredentialsSecretName: credentialsSecretName,
	}
}

// WithGroupID returns a copy of the TestParams with GroupID set.
// GroupID is assigned by Atlas after Group creation and is per-test.
func (p *TestParams) WithGroupID(groupID string) *TestParams {
	copy := *p
	copy.GroupID = groupID
	return &copy
}

// WithGroupName returns a copy of the TestParams with GroupName set.
// GroupName should contain a randomized portion for test isolation.
func (p *TestParams) WithGroupName(groupName string) *TestParams {
	copy := *p
	copy.GroupName = groupName
	return &copy
}

// ReplaceYAML replaces all placeholders in the YAML template with actual values.
// Supported placeholders:
//   - __GROUP_ID__ -> GroupID
//   - __ORG_ID__ -> OrgID
//   - __GROUP_NAME__ -> GroupName
//   - __OPERATOR_NAMESPACE__ -> OperatorNamespace
//   - __CREDENTIALS_SECRET_NAME__ -> CredentialsSecretName
func (p *TestParams) ReplaceYAML(yaml string) string {
	result := yaml
	result = strings.ReplaceAll(result, "__GROUP_ID__", p.GroupID)
	result = strings.ReplaceAll(result, "__ORG_ID__", p.OrgID)
	result = strings.ReplaceAll(result, "__GROUP_NAME__", p.GroupName)
	result = strings.ReplaceAll(result, "__OPERATOR_NAMESPACE__", p.OperatorNamespace)
	result = strings.ReplaceAll(result, "__CREDENTIALS_SECRET_NAME__", p.CredentialsSecretName)
	return result
}
