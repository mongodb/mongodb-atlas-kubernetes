package k8s

import (
	"context"
	"encoding/pem"
	"fmt"
	"os"

	. "github.com/onsi/gomega"
	"github.com/sethvargo/go-password/password"
	"gopkg.in/yaml.v2"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"

	v1 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/connectionsecret"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/e2e/utils"
)

func CreateNewClient() (client.Client, error) {
	cfg, err := config.GetConfig()
	if err != nil {
		return nil, err
	}
	err = v1.AddToScheme(scheme.Scheme)
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
	deployment := &v1.AtlasDeployment{}
	err := k8sClient.Get(ctx, client.ObjectKey{Namespace: ns, Name: resourceName}, deployment)
	if err != nil {
		return 0, err
	}
	return int(deployment.Status.ObservedGeneration), nil
}

func GetProjectObservedGeneration(ctx context.Context, k8sClient client.Client, ns, resourceName string) (int, error) {
	project := &v1.AtlasProject{}
	err := k8sClient.Get(ctx, client.ObjectKey{Namespace: ns, Name: resourceName}, project)
	if err != nil {
		return 0, err
	}
	return int(project.Status.ObservedGeneration), nil
}

func GetProjectStatusCondition(ctx context.Context, k8sClient client.Client, statusType status.ConditionType, ns string, name string) (string, error) {
	project := &v1.AtlasProject{}
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

func GetDeploymentStatusCondition(ctx context.Context, k8sClient client.Client, statusType status.ConditionType, ns string, name string) (string, error) {
	deployment := &v1.AtlasDeployment{}
	err := k8sClient.Get(ctx, client.ObjectKey{Namespace: ns, Name: name}, deployment)
	if err != nil {
		return "", err
	}
	for _, condition := range deployment.Status.Conditions {
		if condition.Type == statusType {
			return string(condition.Status), nil
		}
	}
	return "", fmt.Errorf("condition %s not found. found %v", statusType, deployment.Status.Conditions)
}

func GetDBUserStatusCondition(ctx context.Context, k8sClient client.Client, statusType status.ConditionType, ns string, name string) (string, error) {
	user := &v1.AtlasDatabaseUser{}
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

func GetDeploymentResource(ctx context.Context, k8sClient client.Client, namespace, rName string) (*v1.AtlasDeployment, error) {
	deployment := &v1.AtlasDeployment{}
	err := k8sClient.Get(ctx, client.ObjectKey{Namespace: namespace, Name: rName}, &v1.AtlasDeployment{})
	if err != nil {
		return nil, err
	}
	return deployment, nil
}

func GetK8sDeploymentStateName(ctx context.Context, k8sClient client.Client, ns, rName string) (string, error) {
	deployment := &v1.AtlasDeployment{}
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
	deployment := &v1.AtlasDeployment{}
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
		connectionsecret.TypeLabelKey: connectionsecret.CredLabelVal,
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
	connectionSecret := corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: ns,
			Labels: map[string]string{
				connectionsecret.TypeLabelKey: connectionsecret.CredLabelVal,
			},
		},
		StringData: map[string]string{"orgId": os.Getenv("MCLI_ORG_ID"), "publicApiKey": publicKey, "privateApiKey": privateKey},
	}
	err := k8sClient.Create(ctx, &connectionSecret)
	if err != nil && !k8serrors.IsAlreadyExists(err) {
		return fmt.Errorf("error creating secret: %w", err)
	}
	err = k8sClient.Update(ctx, &connectionSecret)
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

	certFileName := "x509cert.pem"
	certFile, err := os.Create(certFileName)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", certFileName, err)
	}

	err = pem.Encode(certFile, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: cert,
	})
	if err != nil {
		return fmt.Errorf("failed to write data to %s: %w", certFileName, err)
	}
	err = certFile.Close()
	if err != nil {
		return fmt.Errorf("cant close file: %w", err)
	}

	var rawCert []byte
	rawCert, err = os.ReadFile(certFileName)
	if err != nil {
		return fmt.Errorf("failed to read cert file: %w", err)
	}

	certificateSecret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: ns,
		},
		Data: map[string][]byte{
			certFileName: rawCert,
		},
	}
	certificateSecret.Labels = map[string]string{
		connectionsecret.TypeLabelKey: connectionsecret.CredLabelVal,
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
	projectList := &v1.AtlasProjectList{}
	err := k8sClient.List(ctx, projectList, client.InNamespace(ns))
	if err != nil {
		return nil, err
	}
	return yaml.Marshal(projectList)
}

func UserListYaml(ctx context.Context, k8sClient client.Client, ns string) ([]byte, error) {
	userList := &v1.AtlasDatabaseUserList{}
	err := k8sClient.List(ctx, userList, client.InNamespace(ns))
	if err != nil {
		return nil, err
	}
	return yaml.Marshal(userList)
}

func TeamListYaml(ctx context.Context, k8sClient client.Client, ns string) ([]byte, error) {
	teamList := &v1.AtlasTeamList{}
	err := k8sClient.List(ctx, teamList, client.InNamespace(ns))
	if err != nil {
		return nil, err
	}
	return yaml.Marshal(teamList)
}

func DeploymentListYml(ctx context.Context, k8sClient client.Client, ns string) ([]byte, error) {
	deploymentList := &v1.AtlasDeploymentList{}
	err := k8sClient.List(ctx, deploymentList, client.InNamespace(ns))
	if err != nil {
		return nil, err
	}
	return yaml.Marshal(deploymentList)
}
