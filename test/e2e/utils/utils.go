package utils

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
)

//LoadUserProjectConfig load configuration into object
func LoadUserProjectConfig(path string) *v1.AtlasProject {
	var config v1.AtlasProject
	ReadInYAMLFileAndConvert(path, &config)
	return &config
}

//LoadUserClusterConfig load configuration into object
func LoadUserClusterConfig(path string) *v1.AtlasCluster {
	var config v1.AtlasCluster
	ReadInYAMLFileAndConvert(path, &config)
	return &config
}

// ReadInYAMLFileAndConvert reads in the yaml file given by the path given
func ReadInYAMLFileAndConvert(pathToYamlFile string, cnfg interface{}) interface{} {
	// Read in the yaml file at the path given
	yamlFile, err := ioutil.ReadFile(pathToYamlFile)
	if err != nil {
		log.Printf("Error while parsing YAML file %v, error: %s", pathToYamlFile, err)
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
