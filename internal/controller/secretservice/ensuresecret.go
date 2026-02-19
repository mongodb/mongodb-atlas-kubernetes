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

package secretservice

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/kube"
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

	schemeMongoDBSRV string = "mongodb+srv://"
	schemeMongoDB    string = "mongodb://"
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
func Ensure(ctx context.Context, client client.Client, namespace, projectName, projectID, clusterName string, data ConnectionData) (string, error) {
	var getError error
	s := &corev1.Secret{ObjectMeta: metav1.ObjectMeta{
		Name:      formatSecretName(projectName, clusterName, data.DBUserName),
		Namespace: namespace,
	}}
	if getError = client.Get(ctx, kube.ObjectKeyFromObject(s), s); getError != nil && !apierrors.IsNotFound(getError) {
		return "", getError
	}
	if err := fillSecret(s, projectID, clusterName, data); err != nil {
		return "", err
	}
	if getError != nil {
		// Creating
		return s.Name, client.Create(ctx, s)
	}

	return s.Name, client.Update(ctx, s)
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
	if connURL == "" {
		return "", nil
	}

	var prefix string
	switch {
	case strings.HasPrefix(connURL, schemeMongoDBSRV):
		prefix = schemeMongoDBSRV
	case strings.HasPrefix(connURL, schemeMongoDB):
		prefix = schemeMongoDB
	default:
		return "", fmt.Errorf("unsupported MongoDB connection string scheme: %q", connURL)
	}

	rest := connURL[len(prefix):]
	end := len(rest)
	if i := strings.IndexAny(rest, "/?"); i >= 0 {
		end = i
	}
	authority := rest[:end]
	tail := rest[end:]

	if strings.Contains(authority, "@") {
		parts := strings.SplitN(authority, "@", 2)
		authority = parts[1]
	}

	userinfo := url.UserPassword(userName, password).String()
	return prefix + userinfo + "@" + authority + tail, nil
}
