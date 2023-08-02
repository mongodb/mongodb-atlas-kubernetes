package connectionsecret

import (
	"context"
	"fmt"
	"net/url"

	corev1 "k8s.io/api/core/v1"
	apiErrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/kube"
)

const (
	ProjectLabelKey string = "atlas.mongodb.com/project-id"
	ClusterLabelKey string = "atlas.mongodb.com/cluster-name"
	TypeLabelKey           = "atlas.mongodb.com/type"
	CredLabelVal           = "credentials"

	standardKey     string = "connectionStringStandard"
	standardKeySrv  string = "connectionStringStandardSrv"
	privateKey      string = "connectionStringPrivate"
	privateKeySrv   string = "connectionStringPrivateSrv"
	privateShardKey string = "connectionStringPrivateShard"
	userNameKey     string = "username"
	passwordKey     string = "password"
)

type ConnectionData struct {
	DBUserName      string
	Password        string
	ConnURL         string
	SrvConnURL      string
	PrivateConnURLs []PrivateLinkConnURLs
}

type PrivateLinkConnURLs struct {
	PvtConnURL      string
	PvtSrvConnURL   string
	PvtShardConnURL string
}

// Ensure creates or updates the connection Secret for the specific cluster and db user. Returns the name of the Secret
// created.
func Ensure(client client.Client, namespace, projectName, projectID, clusterName string, data ConnectionData) (string, error) {
	var getError error
	s := &corev1.Secret{ObjectMeta: metav1.ObjectMeta{
		Name:      formatSecretName(projectName, clusterName, data.DBUserName),
		Namespace: namespace,
	}}
	if getError = client.Get(context.Background(), kube.ObjectKeyFromObject(s), s); getError != nil && !apiErrors.IsNotFound(getError) {
		return "", getError
	}
	if err := fillSecret(s, projectID, clusterName, data); err != nil {
		return "", err
	}
	if getError != nil {
		// Creating
		return s.Name, client.Create(context.Background(), s)
	}

	return s.Name, client.Update(context.Background(), s)
}

func fillSecret(secret *corev1.Secret, projectID string, clusterName string, data ConnectionData) error {
	var err error
	if data.ConnURL, err = AddCredentialsToConnectionURL(data.ConnURL, data.DBUserName, data.Password); err != nil {
		return err
	}
	if data.SrvConnURL, err = AddCredentialsToConnectionURL(data.SrvConnURL, data.DBUserName, data.Password); err != nil {
		return err
	}
	for idx, privateConn := range data.PrivateConnURLs {
		if data.PrivateConnURLs[idx].PvtConnURL, err = AddCredentialsToConnectionURL(privateConn.PvtConnURL, data.DBUserName, data.Password); err != nil {
			return err
		}
		if data.PrivateConnURLs[idx].PvtSrvConnURL, err = AddCredentialsToConnectionURL(privateConn.PvtSrvConnURL, data.DBUserName, data.Password); err != nil {
			return err
		}
		if data.PrivateConnURLs[idx].PvtShardConnURL, err = AddCredentialsToConnectionURL(privateConn.PvtShardConnURL, data.DBUserName, data.Password); err != nil {
			return err
		}
	}

	secret.Labels = map[string]string{
		TypeLabelKey:    CredLabelVal,
		ProjectLabelKey: projectID,
		ClusterLabelKey: kube.NormalizeLabelValue(clusterName),
	}

	secret.Data = map[string][]byte{
		userNameKey:    []byte(data.DBUserName),
		passwordKey:    []byte(data.Password),
		standardKey:    []byte(data.ConnURL),
		standardKeySrv: []byte(data.SrvConnURL),
		privateKey:     []byte(""),
		privateKeySrv:  []byte(""),
	}

	for idx, privateConn := range data.PrivateConnURLs {
		suffix := getSuffix(idx)
		secret.Data[privateKey+suffix] = []byte(privateConn.PvtConnURL)
		secret.Data[privateKeySrv+suffix] = []byte(privateConn.PvtSrvConnURL)
		secret.Data[privateShardKey+suffix] = []byte(privateConn.PvtShardConnURL)
	}

	return nil
}

func getSuffix(idx int) string {
	if idx == 0 {
		return ""
	}

	return fmt.Sprint(idx)
}

func formatSecretName(projectName, clusterName, dbUserName string) string {
	name := fmt.Sprintf("%s-%s-%s",
		kube.NormalizeIdentifier(projectName),
		kube.NormalizeIdentifier(clusterName),
		kube.NormalizeIdentifier(dbUserName))
	return kube.NormalizeIdentifier(name)
}

func AddCredentialsToConnectionURL(connURL, userName, password string) (string, error) {
	cs, err := url.Parse(connURL)
	if err != nil {
		return "", err
	}
	cs.User = url.UserPassword(userName, password)
	return cs.String(), nil
}
