package v1

import (
	"reflect"
	"sort"
	"testing"

	"github.com/fatih/structtag"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/atlas/mongodbatlas"

	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/provider"
)

var (
	excludedClusterFieldsOurs   = map[string]bool{}
	excludedClusterFieldsTheirs = map[string]bool{}
)

func init() {
	excludedClusterFieldsOurs["projectRef"] = true

	excludedClusterFieldsTheirs["backupEnabled"] = true
	excludedClusterFieldsTheirs["id"] = true
	excludedClusterFieldsTheirs["groupId"] = true

	// Deprecated
	excludedClusterFieldsTheirs["replicationSpec"] = true
	excludedClusterFieldsTheirs["replicationFactor"] = true

	// These fields are shown in the status
	excludedClusterFieldsTheirs["mongoDBVersion"] = true
	excludedClusterFieldsTheirs["mongoURI"] = true
	excludedClusterFieldsTheirs["mongoURIUpdated"] = true
	excludedClusterFieldsTheirs["mongoURIWithOptions"] = true
	excludedClusterFieldsTheirs["connectionStrings"] = true
	excludedClusterFieldsTheirs["srvAddress"] = true
	excludedClusterFieldsTheirs["stateName"] = true
}

func TestCompatibility(t *testing.T) {
	compareStruct(AtlasClusterSpec{}, mongodbatlas.Cluster{}, t)
}

// TestEnums verifies that replacing the strings with "enum" in Atlas Operator works correctly and is (de)serialized
// into the correct Atlas Cluster
func TestEnums(t *testing.T) {
	atlasCluster := mongodbatlas.Cluster{
		ProviderSettings: &mongodbatlas.ProviderSettings{
			ProviderName: "AWS",
		},
		ClusterType: "GEOSHARDED",
	}
	operatorCluster := AtlasClusterSpec{
		ProviderSettings: &ProviderSettingsSpec{
			ProviderName: provider.ProviderAWS,
		},
		ClusterType: TypeGeoSharded,
	}
	transformedCluster, err := operatorCluster.Cluster()
	assert.NoError(t, err)
	assert.Equal(t, atlasCluster, *transformedCluster)
}

func compareStruct(ours interface{}, their interface{}, t *testing.T) {
	ourFields := getAllFieldsSorted(ours, excludedClusterFieldsOurs)
	theirFields := getAllFieldsSorted(their, excludedClusterFieldsTheirs)

	// Comparing the fields in sorted order first
	ourStructName := reflect.ValueOf(ours).Type().Name()
	theirStructName := reflect.ValueOf(their).Type().Name()
	assert.Equal(t, ourFields, theirFields, "The fields for structs [ours: %s, theirs: %s] don't match!", ourStructName, theirStructName)

	// Then recurse into the fields of type struct
	structFieldsTags := getAllStructFieldTags(ours, excludedClusterFieldsOurs)
	for _, field := range structFieldsTags {
		ourStructField := findFieldValueByTag(ours, field)
		theirStructField := findFieldValueByTag(their, field)

		compareStruct(ourStructField, theirStructField, t)
	}
}

func findFieldValueByTag(theStruct interface{}, tag string) interface{} {
	o := reflect.ValueOf(theStruct)
	for i := 0; i < o.NumField(); i++ {
		theTag := parseJSONName(o.Type().Field(i).Tag)
		if theTag == tag {
			v := reflect.New(o.Type().Field(i).Type.Elem()).Elem().Interface()
			return v
		}
	}
	panic("Field with tag not found")
}

func getAllStructFieldTags(theStruct interface{}, excludedFields map[string]bool) []string {
	o := reflect.ValueOf(theStruct)
	var res []string
	for i := 0; i < o.NumField(); i++ {
		theTag := parseJSONName(o.Type().Field(i).Tag)
		ft := o.Field(i).Type()
		if ft.Kind() == reflect.Ptr {
			ft = ft.Elem()
		}
		if _, ok := excludedFields[theTag]; !ok && ft.Kind() == reflect.Struct {
			res = append(res, theTag)
		}
	}
	return res
}

func getAllFieldsSorted(theStruct interface{}, excluded map[string]bool) []string {
	var res []string
	o := reflect.ValueOf(theStruct)
	for i := 0; i < o.NumField(); i++ {
		theTag := parseJSONName(o.Type().Field(i).Tag)
		if _, ok := excluded[theTag]; !ok {
			res = append(res, theTag)
		}
	}
	sort.Strings(res)
	return res
}

func parseJSONName(t reflect.StructTag) string {
	tags, err := structtag.Parse(string(t))
	if err != nil {
		panic(err)
	}
	jsonTag, err := tags.Get("json")
	if err != nil {
		panic(err)
	}
	return jsonTag.Name
}
