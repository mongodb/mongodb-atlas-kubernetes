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

package k8s

import (
	"context"
	"encoding/pem"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	. "github.com/onsi/gomega"
	"github.com/sethvargo/go-password/password"
	"gopkg.in/yaml.v2"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	k8scfg "sigs.k8s.io/controller-runtime/pkg/client/config"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/secretservice"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/config"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/utils"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e2/kube"
)

func init() {
	err := akov2.AddToScheme(scheme.Scheme)
	if err != nil {
		log.Fatalf("failed to preload Kubernetes schemas: %v", err)
	}
}

func CreateNewClient() (client.Client, error) {
	cfg, err := k8scfg.GetConfig()
	if err != nil {
		return nil, err
	}
	k8sClient, err := client.New(cfg, client.Options{Scheme: scheme.Scheme})
	if err != nil {
		return nil, err
	}
	return k8sClient, nil
}

// GetPodStatus status.phase
func GetPodStatus(ctx context.Context, k8sClient client.Client, ns string) (string, error) {
	pod := &corev1.PodList{}
	err := k8sClient.List(ctx, pod, client.InNamespace(ns), client.MatchingLabels{"app.kubernetes.io/instance": "mongodb-atlas-kubernetes-operator"})
	if err != nil {
		return "", err
	}
	if len(pod.Items) == 0 {
		return "", fmt.Errorf("no pods found")
	}
	return string(pod.Items[0].Status.Phase), nil
}

func GetDeploymentObservedGeneration(ctx context.Context, k8sClient client.Client, ns, resourceName string) (int, error) {
	deployment := &akov2.AtlasDeployment{}
	err := k8sClient.Get(ctx, client.ObjectKey{Namespace: ns, Name: resourceName}, deployment)
	if err != nil {
		return 0, err
	}
	return int(deployment.Status.ObservedGeneration), nil
}

func GetProjectObservedGeneration(ctx context.Context, k8sClient client.Client, ns, resourceName string) (int, error) {
	project := &akov2.AtlasProject{}
	err := k8sClient.Get(ctx, client.ObjectKey{Namespace: ns, Name: resourceName}, project)
	if err != nil {
		return 0, err
	}
	return int(project.Status.ObservedGeneration), nil
}

func GetProjectStatusCondition(ctx context.Context, k8sClient client.Client, statusType api.ConditionType, ns string, name string) (string, error) {
	project := &akov2.AtlasProject{}
	err := k8sClient.Get(ctx, client.ObjectKey{Namespace: ns, Name: name}, project)
	if err != nil {
		return "", err
	}
	for _, condition := range project.Status.Conditions {
		if condition.Type == statusType {
			return string(condition.Status), nil
		}
	}
	return "", fmt.Errorf("condition %s not found. found %v", statusType, project.Status.Conditions)
}

func GetStatusCondition(ctx context.Context, k8sClient client.Client, obj kube.ObjectWithStatus, statusType string) (metav1.Condition, error) {
	err := k8sClient.Get(ctx, client.ObjectKeyFromObject(obj), obj)
	if err != nil {
		return metav1.Condition{}, err
	}

	for _, condition := range obj.GetConditions() {
		if condition.Type == statusType {
			return condition, nil
		}
	}

	return metav1.Condition{}, fmt.Errorf("condition %s not found. found %v", statusType, obj.GetConditions())
}

func GetDeploymentStatusCondition(ctx context.Context, k8sClient client.Client, statusType api.ConditionType, ns string, name string) (string, error) {
	deployment := &akov2.AtlasDeployment{}
	err := k8sClient.Get(ctx, client.ObjectKey{Namespace: ns, Name: name}, deployment)
	if err != nil {
		return "", err
	}

	if deployment.GetObjectMeta().GetGeneration() != deployment.GetStatus().GetObservedGeneration() {
		return "", errors.New("object is not updated")
	}

	for _, condition := range deployment.Status.Conditions {
		if condition.Type == statusType {
			return string(condition.Status), nil
		}
	}
	return "", fmt.Errorf("condition %s not found. found %v", statusType, deployment.Status.Conditions)
}

func GetDBUserStatusCondition(ctx context.Context, k8sClient client.Client, statusType api.ConditionType, ns string, name string) (string, error) {
	user := &akov2.AtlasDatabaseUser{}
	err := k8sClient.Get(ctx, client.ObjectKey{Namespace: ns, Name: name}, user)
	if err != nil {
		return "", err
	}
	for _, condition := range user.Status.Conditions {
		if condition.Type == statusType {
			return string(condition.Status), nil
		}
	}
	return "", fmt.Errorf("condition %s not found. found %v", statusType, user.Status.Conditions)
}

func GetPodStatusPhaseByLabel(ctx context.Context, k8sClient client.Client, ns, labelKey, labelValue string) (string, error) {
	pod := &corev1.PodList{}
	err := k8sClient.List(ctx, pod, client.InNamespace(ns), client.MatchingLabels{labelKey: labelValue})
	if err != nil {
		return "", err
	}
	if len(pod.Items) == 0 {
		return "", fmt.Errorf("no pods found")
	}
	return string(pod.Items[0].Status.Phase), nil
}

func GetDeploymentResource(ctx context.Context, k8sClient client.Client, namespace, rName string) (*akov2.AtlasDeployment, error) {
	deployment := &akov2.AtlasDeployment{}
	err := k8sClient.Get(ctx, client.ObjectKey{Namespace: namespace, Name: rName}, &akov2.AtlasDeployment{})
	if err != nil {
		return nil, err
	}
	return deployment, nil
}

func GetK8sDeploymentStateName(ctx context.Context, k8sClient client.Client, ns, rName string) (string, error) {
	deployment := &akov2.AtlasDeployment{}
	err := k8sClient.Get(ctx, client.ObjectKey{Namespace: ns, Name: rName}, deployment)
	if err != nil {
		return "", err
	}
	return deployment.Status.StateName, nil
}

func DeleteNamespace(ctx context.Context, k8sClient client.Client, ns string) error {
	namespace := &corev1.Namespace{}
	err := k8sClient.Get(ctx, client.ObjectKey{Name: ns}, namespace)
	if err != nil {
		return err
	}
	err = k8sClient.Delete(ctx, namespace)
	if err != nil {
		return err
	}
	return nil
}

func DeleteDeployment(ctx context.Context, k8sClient client.Client, ns, name string) error {
	deployment := &akov2.AtlasDeployment{}
	err := k8sClient.Get(ctx, client.ObjectKey{Namespace: ns, Name: name}, deployment)
	if err != nil {
		return err
	}
	err = k8sClient.Delete(ctx, deployment)
	if err != nil {
		return err
	}
	return nil
}

func CreateNamespace(ctx context.Context, k8sClient client.Client, ns string) error {
	namespace := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: ns,
		},
	}
	err := k8sClient.Create(ctx, namespace)
	if err != nil {
		if !k8serrors.IsAlreadyExists(err) {
			return err
		}
	}
	return nil
}

func CreateRandomNamespace(ctx context.Context, k8sClient client.Client, generateName string) (*corev1.Namespace, error) {
	namespace := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: generateName,
		},
	}
	err := k8sClient.Create(ctx, namespace)
	if err != nil {
		if !k8serrors.IsAlreadyExists(err) {
			return nil, err
		}
	}
	return namespace, nil
}

func CreateRandomUserSecret(ctx context.Context, k8sClient client.Client, name, ns string) error {
	secret, _ := password.Generate(10, 3, 0, false, false)
	return CreateUserSecret(ctx, k8sClient, secret, name, ns)
}

func CreateUserSecret(ctx context.Context, k8sClient client.Client, secret, name, ns string) error {
	secretObj := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: ns,
		},
		StringData: map[string]string{
			"password": secret,
		},
	}
	secretObj.SetLabels(map[string]string{
		secretservice.TypeLabelKey: secretservice.CredLabelVal,
	})
	err := k8sClient.Create(ctx, secretObj)
	if err != nil {
		return err
	}
	return nil
}

func CreateDefaultSecret(ctx context.Context, k8sClient client.Client, name, ns string) {
	Expect(CreateSecret(ctx, k8sClient, os.Getenv("MCLI_PUBLIC_API_KEY"), os.Getenv("MCLI_PRIVATE_API_KEY"), name, ns)).Should(Succeed())
}

func CreateSecret(ctx context.Context, k8sClient client.Client, publicKey, privateKey, name, ns string) error {
	secretservice := corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: ns,
			Labels: map[string]string{
				secretservice.TypeLabelKey: secretservice.CredLabelVal,
			},
		},
		StringData: map[string]string{"orgId": os.Getenv("MCLI_ORG_ID"), "publicApiKey": publicKey, "privateApiKey": privateKey},
	}
	err := k8sClient.Create(ctx, &secretservice)
	if err != nil && !k8serrors.IsAlreadyExists(err) {
		return fmt.Errorf("error creating secret: %w", err)
	}
	err = k8sClient.Update(ctx, &secretservice)
	if err != nil {
		return fmt.Errorf("error updating existing secret: %w", err)
	}
	return nil
}

func DeleteKey(ctx context.Context, k8sClient client.Client, keyName, ns string) error {
	secret := &corev1.Secret{}
	err := k8sClient.Get(ctx, client.ObjectKey{Namespace: ns, Name: keyName}, secret)
	if err != nil {
		return err
	}
	err = k8sClient.Delete(ctx, secret)
	if err != nil {
		return err
	}
	return nil
}

func CreateCertificateX509(ctx context.Context, k8sClient client.Client, name, ns string) error {
	cert, _, _, err := utils.GenerateX509Cert()
	if err != nil {
		return fmt.Errorf("error generating x509 cert: %w", err)
	}

	certFile, err := os.Create(filepath.Clean(config.PEMCertFileName))
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", config.PEMCertFileName, err)
	}

	err = pem.Encode(certFile, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: cert,
	})
	if err != nil {
		return fmt.Errorf("failed to write data to %s: %w", config.PEMCertFileName, err)
	}
	err = certFile.Close()
	if err != nil {
		return fmt.Errorf("cant close file: %w", err)
	}

	var rawCert []byte
	rawCert, err = os.ReadFile(filepath.Clean(config.PEMCertFileName))
	if err != nil {
		return fmt.Errorf("failed to read cert file: %w", err)
	}

	certificateSecret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: ns,
		},
		Data: map[string][]byte{
			filepath.Base(config.PEMCertFileName): rawCert,
		},
	}
	certificateSecret.Labels = map[string]string{
		secretservice.TypeLabelKey: secretservice.CredLabelVal,
	}
	err = k8sClient.Create(ctx, certificateSecret)
	if err != nil {
		if !k8serrors.IsAlreadyExists(err) {
			return fmt.Errorf("error creating certificate secret: %w, %v", err, certificateSecret)
		}
	}
	return nil
}

func ProjectListYaml(ctx context.Context, k8sClient client.Client, ns string) ([]byte, error) {
	projectList := &akov2.AtlasProjectList{}
	err := k8sClient.List(ctx, projectList, client.InNamespace(ns))
	if err != nil {
		return nil, err
	}
	return yaml.Marshal(projectList)
}

func UserListYaml(ctx context.Context, k8sClient client.Client, ns string) ([]byte, error) {
	userList := &akov2.AtlasDatabaseUserList{}
	err := k8sClient.List(ctx, userList, client.InNamespace(ns))
	if err != nil {
		return nil, err
	}
	return yaml.Marshal(userList)
}

func TeamListYaml(ctx context.Context, k8sClient client.Client, ns string) ([]byte, error) {
	teamList := &akov2.AtlasTeamList{}
	err := k8sClient.List(ctx, teamList, client.InNamespace(ns))
	if err != nil {
		return nil, err
	}
	return yaml.Marshal(teamList)
}

func AtlasOrgSettingsListYaml(ctx context.Context, k8sClient client.Client, ns string) ([]byte, error) {
	orgSettingsList := &akov2.AtlasOrgSettingsList{}
	err := k8sClient.List(ctx, orgSettingsList, client.InNamespace(ns))
	if err != nil {
		return nil, err
	}
	return yaml.Marshal(orgSettingsList)
}

func DeploymentListYml(ctx context.Context, k8sClient client.Client, ns string) ([]byte, error) {
	deploymentList := &akov2.AtlasDeploymentList{}
	err := k8sClient.List(ctx, deploymentList, client.InNamespace(ns))
	if err != nil {
		return nil, err
	}
	return yaml.Marshal(deploymentList)
}
