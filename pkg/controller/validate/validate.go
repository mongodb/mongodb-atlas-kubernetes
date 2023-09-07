package validate

import (
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"net"
	"net/url"
	"reflect"
	"regexp"
	"strings"

	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/timeutil"

	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/project"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
)

type googleServiceAccountKey struct {
	Type                    string `json:"type"`
	ProjectID               string `json:"project_id"`
	PrivateKeyID            string `json:"private_key_id"`
	PrivateKey              string `json:"private_key"`  // Expects valid PEM key
	ClientEmail             string `json:"client_email"` // Expects a valid email
	ClientID                string `json:"client_id"`
	AuthURI                 string `json:"auth_uri"`                    // Expects valid URL
	TokenURI                string `json:"token_uri"`                   // Expects valid URL
	AuthProviderX509CertURL string `json:"auth_provider_x509_cert_url"` // Expects valid URL
	ClientX509CertURL       string `json:"client_x509_cert_url"`        // Expects valid URL
	UniverseDomain          string `json:"universe_domain"`
}

func DeploymentSpec(deploymentSpec mdbv1.AtlasDeploymentSpec) error {
	var err error

	if allAreNil(deploymentSpec.AdvancedDeploymentSpec, deploymentSpec.ServerlessSpec, deploymentSpec.DeploymentSpec) {
		err = errors.Join(err, errors.New("expected exactly one of spec.deploymentSpec or spec.advancedDepploymentSpec or spec.serverlessSpec to be present, but none were"))
	}

	if moreThanOneIsNonNil(deploymentSpec.AdvancedDeploymentSpec, deploymentSpec.ServerlessSpec, deploymentSpec.DeploymentSpec) {
		err = errors.Join(err, errors.New("expected exactly one of spec.deploymentSpec, spec.advancedDepploymentSpec or spec.serverlessSpec, more than one were present"))
	}

	if deploymentSpec.DeploymentSpec != nil {
		if deploymentSpec.DeploymentSpec.ProviderSettings != nil && (deploymentSpec.DeploymentSpec.ProviderSettings.InstanceSizeName == "" && deploymentSpec.DeploymentSpec.ProviderSettings.ProviderName != "SERVERLESS") {
			err = errors.Join(err, errors.New("must specify instanceSizeName if provider name is not SERVERLESS"))
		}
		if deploymentSpec.DeploymentSpec.ProviderSettings != nil && (deploymentSpec.DeploymentSpec.ProviderSettings.InstanceSizeName != "" && deploymentSpec.DeploymentSpec.ProviderSettings.ProviderName == "SERVERLESS") {
			err = errors.Join(err, errors.New("must not specify instanceSizeName if provider name is SERVERLESS"))
		}
	}

	if deploymentSpec.AdvancedDeploymentSpec != nil {
		instanceSizeErr := instanceSizeForAdvancedDeployment(deploymentSpec.AdvancedDeploymentSpec.ReplicationSpecs)
		if instanceSizeErr != nil {
			err = errors.Join(err, instanceSizeErr)
		}

		autoscalingErr := autoscalingForAdvancedDeployment(deploymentSpec.AdvancedDeploymentSpec.ReplicationSpecs)
		if autoscalingErr != nil {
			err = errors.Join(err, autoscalingErr)
		}
	}

	return err
}

func Project(project *mdbv1.AtlasProject) error {
	if err := projectIPAccessList(project.Spec.ProjectIPAccessList); err != nil {
		return err
	}

	if err := projectCustomRoles(project.Spec.CustomRoles); err != nil {
		return err
	}

	if err := encryptionAtRest(project.Spec.EncryptionAtRest); err != nil {
		return err
	}

	return nil
}

func DatabaseUser(_ *mdbv1.AtlasDatabaseUser) error {
	return nil
}

func BackupSchedule(bSchedule *mdbv1.AtlasBackupSchedule, deployment *mdbv1.AtlasDeployment) error {
	var err error

	if bSchedule.Spec.Export == nil && bSchedule.Spec.AutoExportEnabled {
		err = errors.Join(err, errors.New("you must specify export policy when auto export is enabled"))
	}

	replicaSets := map[string]struct{}{}
	if deployment.Status.ReplicaSets != nil {
		for _, replicaSet := range deployment.Status.ReplicaSets {
			replicaSets[replicaSet.ID] = struct{}{}
		}
	}

	for position, copySetting := range bSchedule.Spec.CopySettings {
		if copySetting.RegionName == nil {
			err = errors.Join(err, fmt.Errorf("copy setting at position %d: you must set a region name", position))
		}

		if copySetting.ReplicationSpecID == nil {
			err = errors.Join(err, fmt.Errorf("copy setting at position %d: you must set a valid ReplicationSpecID", position))
		} else if _, ok := replicaSets[*copySetting.ReplicationSpecID]; !ok {
			err = errors.Join(err, fmt.Errorf("copy setting at position %d: referenced ReplicationSpecID is invalid", position))
		}

		if copySetting.ShouldCopyOplogs != nil && *copySetting.ShouldCopyOplogs {
			if deployment.Spec.AdvancedDeploymentSpec != nil &&
				(deployment.Spec.AdvancedDeploymentSpec.PitEnabled == nil ||
					!*deployment.Spec.AdvancedDeploymentSpec.PitEnabled) {
				err = errors.Join(err, fmt.Errorf("copy setting at position %d: you must enable pit before enable copyOplogs", position))
			}

			if deployment.Spec.DeploymentSpec != nil &&
				(deployment.Spec.DeploymentSpec.PitEnabled == nil ||
					!*deployment.Spec.DeploymentSpec.PitEnabled) {
				err = errors.Join(err, fmt.Errorf("copy setting at position %d: you must enable pit before enable copyOplogs", position))
			}
		}
	}

	return err
}

func getNonNilCount(values ...interface{}) int {
	nonNilCount := 0
	for _, v := range values {
		if !reflect.ValueOf(v).IsNil() {
			nonNilCount += 1
		}
	}
	return nonNilCount
}

// allAreNil returns true if all elements are nil.
func allAreNil(values ...interface{}) bool {
	return getNonNilCount(values...) == 0
}

// moreThanOneIsNil returns true if there are more than one non nil elements.
func moreThanOneIsNonNil(values ...interface{}) bool {
	return getNonNilCount(values...) > 1
}

func instanceSizeForAdvancedDeployment(replicationSpecs []*mdbv1.AdvancedReplicationSpec) error {
	var instanceSize string
	err := errors.New("instance size must be the same for all nodes in all regions and across all replication specs for advanced deployment ")

	isInstanceSizeEqual := func(nodeInstanceType string) bool {
		if instanceSize == "" {
			instanceSize = nodeInstanceType
		}

		return nodeInstanceType == instanceSize
	}

	for _, replicationSpec := range replicationSpecs {
		for _, regionSpec := range replicationSpec.RegionConfigs {
			if instanceSize == "" {
				instanceSize = regionSpec.ElectableSpecs.InstanceSize
			}

			if regionSpec.ElectableSpecs != nil && !isInstanceSizeEqual(regionSpec.ElectableSpecs.InstanceSize) {
				return err
			}

			if regionSpec.ReadOnlySpecs != nil && !isInstanceSizeEqual(regionSpec.ReadOnlySpecs.InstanceSize) {
				return err
			}

			if regionSpec.AnalyticsSpecs != nil && !isInstanceSizeEqual(regionSpec.AnalyticsSpecs.InstanceSize) {
				return err
			}
		}
	}

	return nil
}

func autoscalingForAdvancedDeployment(replicationSpecs []*mdbv1.AdvancedReplicationSpec) error {
	var autoscaling *mdbv1.AdvancedAutoScalingSpec
	first := true

	for _, replicationSpec := range replicationSpecs {
		for _, regionSpec := range replicationSpec.RegionConfigs {
			if first {
				autoscaling = regionSpec.AutoScaling
				first = false
			}

			if cmp.Diff(autoscaling, regionSpec.AutoScaling, cmpopts.EquateEmpty()) != "" {
				return errors.New("autoscaling must be the same for all regions and across all replication specs for advanced deployment ")
			}
		}
	}

	return nil
}

func projectIPAccessList(ipAccessList []project.IPAccessList) error {
	if len(ipAccessList) == 0 {
		return nil
	}

	var err error
	for _, item := range ipAccessList {
		if item.IPAddress == "" && item.CIDRBlock == "" && item.AwsSecurityGroup == "" {
			err = errors.Join(err, errors.New("invalid config! one of option must be configured"))
		}

		if item.CIDRBlock != "" {
			if item.AwsSecurityGroup != "" || item.IPAddress != "" {
				err = errors.Join(err, errors.New("don't set ipAddress or awsSecurityGroup when configuring cidrBlock"))
			}

			_, _, cidrErr := net.ParseCIDR(item.CIDRBlock)
			if cidrErr != nil {
				err = errors.Join(err, fmt.Errorf("invalid cidrBlock: %s", item.CIDRBlock))
			}
		}

		if item.IPAddress != "" {
			if item.AwsSecurityGroup != "" || item.CIDRBlock != "" {
				err = errors.Join(err, errors.New("don't set cidrBlock or awsSecurityGroup when configuring ipAddress"))
			}

			ip := net.ParseIP(item.IPAddress)
			if ip == nil {
				err = errors.Join(err, fmt.Errorf("invalid ipAddress: %s", item.IPAddress))
			}
		}

		if item.AwsSecurityGroup != "" {
			if item.IPAddress != "" || item.CIDRBlock != "" {
				err = errors.Join(err, errors.New("don't set cidrBlock or ipAddress when configuring awsSecurityGroup"))
			}

			reg := regexp.MustCompile("^([0-9]*/)?sg-([0-9]*)")
			if !reg.MatchString(item.AwsSecurityGroup) {
				err = errors.Join(err, fmt.Errorf("invalid awsSecurityGroup: %s", item.AwsSecurityGroup))
			}
		}

		if item.DeleteAfterDate != "" {
			_, delErr := timeutil.ParseISO8601(item.DeleteAfterDate)
			if delErr != nil {
				err = errors.Join(err, fmt.Errorf("invalid deleteAfterDate: %s. value should follow ISO8601 format", item.DeleteAfterDate))
			}
		}
	}

	return err
}

func projectCustomRoles(customRoles []mdbv1.CustomRole) error {
	if len(customRoles) == 0 {
		return nil
	}

	var err error
	customRolesMap := map[string]struct{}{}

	for _, customRole := range customRoles {
		if _, ok := customRolesMap[customRole.Name]; ok {
			err = errors.Join(err, fmt.Errorf("the custom role \"%s\" is duplicate. custom role name must be unique", customRole.Name))
		}

		customRolesMap[customRole.Name] = struct{}{}
	}

	return err
}

func encryptionAtRest(encryption *mdbv1.EncryptionAtRest) error {
	if encryption != nil &&
		encryption.GoogleCloudKms.Enabled != nil &&
		*encryption.GoogleCloudKms.Enabled {
		if encryption.GoogleCloudKms.ServiceAccountKey == "" {
			return fmt.Errorf("missing Google Service Account Key but GCP KMS is enabled")
		}
		if err := gceServiceAccountKey(encryption.GoogleCloudKms.ServiceAccountKey); err != nil {
			return fmt.Errorf("failed to validate Google Service Account Key: %w", err)
		}
	}
	return nil
}

func unfilter(key string) string {
	return strings.ReplaceAll(key, "\\\\n", "\\n")
}

func gceServiceAccountKey(key string) error {
	emptyKey := googleServiceAccountKey{}
	gceSAKey := googleServiceAccountKey{}
	if err := json.Unmarshal(([]byte)(unfilter(key)), &gceSAKey); err != nil {
		return fmt.Errorf("invalid service account key format: %w", err)
	}
	if emptyKey == gceSAKey {
		return fmt.Errorf("invalid empty service account key")
	}
	for _, rawURL := range []string{gceSAKey.AuthURI,
		gceSAKey.TokenURI,
		gceSAKey.ClientX509CertURL,
		gceSAKey.AuthProviderX509CertURL} {
		if _, err := url.ParseRequestURI(rawURL); err != nil {
			return fmt.Errorf("invalid URL address %q: %w", rawURL, err)
		}
	}
	block, _ := pem.Decode([]byte(gceSAKey.PrivateKey))
	if block == nil || !strings.HasSuffix(block.Type, "PRIVATE KEY") {
		return fmt.Errorf("failed to decode PEM block containing a private key")
	}

	err := assertParsePrivateKey(block.Bytes)
	if err != nil {
		return fmt.Errorf("failed to parse PEM private key: %w", err)
	}
	return nil
}

func assertParsePrivateKey(key []byte) error {
	_, err := x509.ParsePKCS1PrivateKey(key)
	if err != nil && strings.Contains(err.Error(), "ParsePKCS8PrivateKey") {
		_, err = x509.ParsePKCS8PrivateKey(key)
	}
	return err
}
