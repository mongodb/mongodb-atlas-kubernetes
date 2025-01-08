/*
Copyright 2025 MongoDB.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/common"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/cel"
)

func TestIPAccessListProjectRefCELValidations(t *testing.T) {
	launchProjectRefCELTests(
		t,
		func(pdr *ProjectDualReference) AtlasCustomResource {
			pe := AtlasIPAccessList{}
			if pdr != nil {
				setDualRef(pe.ProjectDualRef(), pdr)
			}
			return &pe
		},
		"../../config/crd/bases/atlas.mongodb.com_atlasipaccesslists.yaml",
	)
}

func TestEntryUniqueness(t *testing.T) {
	testCases := map[string]struct {
		obj            AtlasCustomResource
		expectedErrors []string
	}{
		"should fail if IP and CIDR are set in the same entry": {
			obj: &AtlasIPAccessList{
				Spec: AtlasIPAccessListSpec{
					ProjectDualReference: ProjectDualReference{
						ProjectRef: &common.ResourceRefNamespaced{
							Name: "my-project",
						},
					},
					Entries: []IPAccessEntry{
						{
							IPAddress: "10.0.0.1",
							CIDRBlock: "192.168.0.0/24",
						},
					},
				},
			},
			expectedErrors: []string{"spec.entries[0]: Invalid value: \"object\": Only one of ipAddress, cidrBlock, or awsSecurityGroup may be set."},
		},
		"should fail if IP and AWS SG are set in the same entry": {
			obj: &AtlasIPAccessList{
				Spec: AtlasIPAccessListSpec{
					ProjectDualReference: ProjectDualReference{
						ProjectRef: &common.ResourceRefNamespaced{
							Name: "my-project",
						},
					},
					Entries: []IPAccessEntry{
						{
							IPAddress:        "10.0.0.1",
							AwsSecurityGroup: "sg-123456",
						},
					},
				},
			},
			expectedErrors: []string{"spec.entries[0]: Invalid value: \"object\": Only one of ipAddress, cidrBlock, or awsSecurityGroup may be set."},
		},
		"should fail if CIDR and AWS SG are set in the same entry": {
			obj: &AtlasIPAccessList{
				Spec: AtlasIPAccessListSpec{
					ProjectDualReference: ProjectDualReference{
						ProjectRef: &common.ResourceRefNamespaced{
							Name: "my-project",
						},
					},
					Entries: []IPAccessEntry{
						{
							CIDRBlock:        "192.168.0.0/24",
							AwsSecurityGroup: "sg-123456",
						},
					},
				},
			},
			expectedErrors: []string{"spec.entries[0]: Invalid value: \"object\": Only one of ipAddress, cidrBlock, or awsSecurityGroup may be set."},
		},
		"should fail with multiple entries wrongly set": {
			obj: &AtlasIPAccessList{
				Spec: AtlasIPAccessListSpec{
					ProjectDualReference: ProjectDualReference{
						ProjectRef: &common.ResourceRefNamespaced{
							Name: "my-project",
						},
					},
					Entries: []IPAccessEntry{
						{
							IPAddress: "10.0.0.1",
							CIDRBlock: "192.168.0.0/24",
						},
						{
							CIDRBlock:        "192.168.1.0/24",
							AwsSecurityGroup: "sg-123456",
						},
					},
				},
			},
			expectedErrors: []string{
				"spec.entries[0]: Invalid value: \"object\": Only one of ipAddress, cidrBlock, or awsSecurityGroup may be set.",
				"spec.entries[1]: Invalid value: \"object\": Only one of ipAddress, cidrBlock, or awsSecurityGroup may be set.",
			},
		},
		"should pass when correctly configured": {
			obj: &AtlasIPAccessList{
				Spec: AtlasIPAccessListSpec{
					ProjectDualReference: ProjectDualReference{
						ProjectRef: &common.ResourceRefNamespaced{
							Name: "my-project",
						},
					},
					Entries: []IPAccessEntry{
						{
							IPAddress:       "10.0.0.1",
							DeleteAfterDate: pointer.MakePtr(metav1.Now()),
						},
						{
							CIDRBlock: "192.168.0.0/24",
							Comment:   "My Network",
						},
						{
							AwsSecurityGroup: "sg-123456",
						},
					},
				},
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			unstructuredObject, err := runtime.DefaultUnstructuredConverter.ToUnstructured(&tc.obj)
			require.NoError(t, err)

			validator, err := cel.VersionValidatorFromFile(t, "../../config/crd/bases/atlas.mongodb.com_atlasipaccesslists.yaml", "v1")
			assert.NoError(t, err)
			errs := validator(unstructuredObject, nil)

			require.Equal(t, len(tc.expectedErrors), len(errs))

			for i, err := range errs {
				assert.Equal(t, tc.expectedErrors[i], err.Error())
			}
		})
	}
}
