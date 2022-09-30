package utils

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"os"
	"path/filepath"
	"time"

	"github.com/pborman/uuid"
	"github.com/sethvargo/go-password/password"
	yaml "gopkg.in/yaml.v3"

	v1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
)

// LoadUserProjectConfig load configuration into object
func LoadUserProjectConfig(path string) *v1.AtlasProject {
	var config v1.AtlasProject
	ReadInYAMLFileAndConvert(path, &config)
	return &config
}

func RandomName(base string) string {
	randomSuffix := uuid.New()[0:6]
	return fmt.Sprintf("%s-%s", base, randomSuffix)
}

func UserSecretPassword() string {
	return uuid.New()
}

func SaveToFile(path string, data []byte) error {
	err := os.MkdirAll(filepath.Dir(path), os.ModePerm)
	if err != nil {
		return err
	}
	err = os.WriteFile(path, data, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

func JSONToYAMLConvert(cnfg interface{}) []byte {
	var jsonI interface{}
	j, _ := json.Marshal(cnfg)
	err := yaml.Unmarshal(j, &jsonI)
	if err != nil {
		return nil
	}
	y, _ := yaml.Marshal(jsonI)
	return y
}

// ReadInYAMLFileAndConvert reads in the yaml file given by the path given
func ReadInYAMLFileAndConvert(pathToYamlFile string, cnfg interface{}) interface{} {
	// Read in the yaml file at the path given
	yamlFile, err := os.ReadFile(filepath.Clean(pathToYamlFile))
	if err != nil {
		log.Printf("Error while parsing YAML file %v, error: %s", filepath.Clean(pathToYamlFile), err)
	}

	// Map yamlFile to interface
	var body interface{}
	if err := yaml.Unmarshal(yamlFile, &body); err != nil {
		panic(err)
	}

	// Convert yaml to its json counterpart
	body = ConvertYAMLtoJSONHelper(body)

	// Generate json string from data structure
	jsonFormat, err := json.Marshal(body)
	if err != nil {
		panic(err)
	}

	if err := json.Unmarshal(jsonFormat, &cnfg); err != nil {
		panic(err)
	}

	return cnfg
}

// ConvertYAMLtoJSONHelper converts the yaml to json recursively
func ConvertYAMLtoJSONHelper(i interface{}) interface{} {
	switch item := i.(type) {
	case map[interface{}]interface{}:
		document := map[string]interface{}{}
		for k, v := range item {
			document[k.(string)] = ConvertYAMLtoJSONHelper(v)
		}
		return document
	case []interface{}:
		for i, arr := range item {
			item[i] = ConvertYAMLtoJSONHelper(arr)
		}
	}

	return i
}

func GenUniqID() string {
	return uuid.NewRandom().String()
}

func GenID() string {
	id, _ := password.Generate(10, 3, 0, true, true)
	return id
}

func CopyFile(source, target string) {
	data, _ := os.ReadFile(filepath.Clean(source))
	err := os.WriteFile(target, data, os.ModePerm)
	if err != nil {
		panic(err)
	}
}

func GenerateX509Cert() ([]byte, *rsa.PrivateKey, *rsa.PublicKey, error) {
	template := &x509.Certificate{
		IsCA:                  true,
		BasicConstraintsValid: true,
		SubjectKeyId:          []byte{1, 2, 3},
		SerialNumber:          big.NewInt(1234),
		Issuer: pkix.Name{
			CommonName: "x509-user",
		},
		NotBefore: time.Now(),
		NotAfter:  time.Now().AddDate(5, 5, 5),
		DNSNames:  []string{"x509-user"},
		// see http://golang.org/pkg/crypto/x509/#KeyUsage
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:    x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
	}

	privatekey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, nil, err
	}

	publickey := &privatekey.PublicKey

	var parent = template
	cert, err := x509.CreateCertificate(rand.Reader, template, parent, publickey, privatekey)
	if err != nil {
		return nil, nil, nil, err
	}

	return cert, privatekey, publickey, nil
}
