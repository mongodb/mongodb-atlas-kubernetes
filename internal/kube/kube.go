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

package kube

import (
	"fmt"
	"regexp"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/validation"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var invalidStartEnd = regexp.MustCompile(`(^[^a-z0-9]+)|([^a-z0-9]+$)`)

// https://github.com/kubernetes/apimachinery/blob/master/pkg/util/validation/validation.go#L177
var nonIdentifierRegexp = regexp.MustCompile(`[^a-z0-9.]+`)

// https://github.com/kubernetes/apimachinery/blob/master/pkg/util/validation/validation.go#L155
var nonLabelRegexp = regexp.MustCompile(`[^a-z0-9._]+`)

func ObjectKey(namespace, name string) client.ObjectKey {
	return types.NamespacedName{Name: name, Namespace: namespace}
}

func ObjectKeyFromObject(obj metav1.Object) client.ObjectKey {
	return ObjectKey(obj.GetNamespace(), obj.GetName())
}

// NormalizeIdentifier returns the 'name' "normalized" for the standard identifier in Kubernetes. All non-allowed symbols are replaced with
// dashes.
// See https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#dns-subdomain-names
func NormalizeIdentifier(name string) string {
	if errs := validation.IsDNS1123Subdomain(name); len(errs) == 0 {
		return name
	}
	return normalize(name, 253, nonIdentifierRegexp)
}

// NormalizeLabelValue returns the 'name' "normalized" for the label value. All non-allowed symbols are replaced with
// dashes.
// See https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#dns-label-names
func NormalizeLabelValue(name string) string {
	if errs := validation.IsValidLabelValue(name); len(errs) == 0 {
		return name
	}
	return normalize(name, 63, nonLabelRegexp)
}

// ParseDeploymentNameFromPodName returns the name of Deployment by Pod Name. The Pods for Deployments have two hashes
// parts as the first one is generated for the ReplicaSet resource created and the second - for the Pod itself.
// Example:
// - Deployment: "prometheus-adapter"
// - ReplicaSet: "prometheus-adapter-65c6cb864f"
// - Pod: "prometheus-adapter-797f946f88-97f2q"
func ParseDeploymentNameFromPodName(podName string) (string, error) {
	parts := strings.Split(podName, "-")
	if len(parts) <= 2 {
		return "", fmt.Errorf(`the Pod name must follow the format "<deployment_name>-797f946f88-97f2q" but got %s`, podName)
	}
	return strings.Join(parts[0:len(parts)-2], "-"), nil
}

// Dev note: the algorithm tries to replace the invalid characters with '-' (or simply omit it replacing is not possible)
// Note, that this algorithm is not ideal - e.g. it won't fix the following: "a.#b" ("a._b" is still not a valid output - as
// nonalphanumeric symbols cannot go together) though this doesn't intend to work in ALL the cases but in the MAJORITY instead
func normalize(name string, limit int, regexp *regexp.Regexp) string {
	if len(name) >= limit {
		name = name[:limit]
	}
	name = strings.ToLower(name)
	name = invalidStartEnd.ReplaceAllString(name, "") // makes sure start & end are alphanumeric
	name = regexp.ReplaceAllString(name, "-")         // replaces every sequence of invalid runes with a single "-"
	return name
}
