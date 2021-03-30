package utils

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	yaml "gopkg.in/yaml.v3"

	"github.com/pborman/uuid"

	v1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
)

// LoadUserProjectConfig load configuration into object
func LoadUserProjectConfig(path string) *v1.AtlasProject {
	var config v1.AtlasProject
	ReadInYAMLFileAndConvert(path, &config)
	return &config
}

func SaveToFile(path string, data []byte) error {
	err := os.MkdirAll(filepath.Dir(path), os.ModePerm)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(path, data, os.ModePerm)
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
	yamlFile, err := ioutil.ReadFile(filepath.Clean(pathToYamlFile))
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

func CopyFile(source, target string) {
	data, _ := ioutil.ReadFile(filepath.Clean(source))
	err := ioutil.WriteFile(target, data, os.ModePerm)
	if err != nil {
		panic(err)
	}
}
