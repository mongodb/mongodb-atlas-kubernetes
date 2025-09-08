package translate

import (
	"errors"
	"fmt"
	"reflect"

	"gopkg.in/yaml.v3"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const SetFallbackNamespace = "."

type DependencyFinder interface {
	Find(name, namespace string) client.Object
}

type StaticDependencies struct {
	deps              map[string]client.Object
	fallbackNamespace string
}

type depsBuilder struct {
	StaticDependencies
}

func (db *depsBuilder) dependencies() []client.Object {
	return []client.Object{}
}

// NewStaticDependencies creates a static set of find-able Kubernetes client.Objects
func NewStaticDependencies(fallbackNamespace string, objs ...client.Object) StaticDependencies {
	deps := map[string]client.Object{}
	for _, obj := range objs {
		deps[client.ObjectKeyFromObject(obj).String()] = obj
	}
	return StaticDependencies{
		deps:              deps,
		fallbackNamespace: fallbackNamespace,
	}
}

// Find will reteieve the object with matching name and namespace if present in the static set
func (sd StaticDependencies) Find(name, namespace string) client.Object {
	ns := namespace
	if ns == SetFallbackNamespace {
		ns = sd.fallbackNamespace
	}
	return sd.deps[client.ObjectKey{Name: name, Namespace: ns}.String()]
}

// Translator allows to translate back and forth between a CRD schema version
// and SDK API structures of a certain version
type Translator struct {
	crd  CRDInfo
	sdk  SDKInfo
	deps DependencyFinder
}

// SDKInfo holds the SDK version information
type SDKInfo struct {
	version string
}

// NewTranslator creates a translator for a particular CRD and SDK version pairs,
// and with a particular set of known Kubernetes dependencies
func NewTranslator(crd *apiextensionsv1.CustomResourceDefinition, crdVersion string, sdkVersion string, deps DependencyFinder) *Translator {
	return &Translator{
		crd:  CRDInfo{definition: crd, version: crdVersion},
		sdk:  SDKInfo{version: sdkVersion},
		deps: deps,
	}
}

// FromAPI translaters a source API structure into a Kubernetes object
func FromAPI[S any, T any, P interface {
	*T
	client.Object
}](t *Translator, target P, source *S) ([]client.Object, error) {
	sourceUnstructured, err := toUnstructured(source)
	if err != nil {
		return nil, fmt.Errorf("failed to convert API source value to unstructured: %w", err)
	}

	targetUnstructured := map[string]any{}

	versionedSpec := map[string]any{}
	copyFields(versionedSpec, sourceUnstructured)
	createField(targetUnstructured, versionedSpec, "spec", t.sdk.version, "entry")

	versionedStatus := map[string]any{}
	copyFields(versionedStatus, sourceUnstructured)
	createField(targetUnstructured, versionedStatus, "status", t.sdk.version)

	extraObjects, err := t.processKubeMappings(targetUnstructured)
	if err != nil {
		return nil, fmt.Errorf("failed to process API mappings: %w", err)
	}
	if err := fromUnstructured(target, targetUnstructured); err != nil {
		return nil, fmt.Errorf("failed set structured kubernetes object from unstructured: %w", err)
	}
	return append([]client.Object{target}, extraObjects...), nil
}

// ToAPI translates a source Kubernetes spec into a target API structure
func ToAPI[T any](t *Translator, target *T, source client.Object) error {
	specVersion := selectVersion(&t.crd.definition.Spec, t.crd.version)
	kind := t.crd.definition.Spec.Names.Kind
	props, err := getOpenAPIProperties(kind, specVersion)
	if err != nil {
		return fmt.Errorf("failed to enumerate CRD schema properties: %w", err)
	}
	specProps, err := getSpecPropertiesFor(kind, props, "spec")
	if err != nil {
		return fmt.Errorf("failed to enumerate CRD spec properties: %w", err)
	}
	if _, ok := specProps[t.sdk.version]; !ok {
		return fmt.Errorf("failed to match the CRD spec version %q in schema", t.sdk.version)
	}
	unstructuredSrc, err := toUnstructured(source)
	if err != nil {
		return fmt.Errorf("failed to convert k8s source value to unstructured: %w", err)
	}
	targetUnstructured := map[string]any{}
	value, err := accessField[map[string]any](unstructuredSrc, "spec", t.sdk.version)
	if err != nil {
		return fmt.Errorf("failed to access source spec value: %w", err)
	}

	if err := t.processAPIMappings(value); err != nil {
		return fmt.Errorf("failed to process API mappings: %w", err)
	}

	targetType := reflect.TypeOf(target).Elem()
	if targetType.Kind() != reflect.Struct {
		return fmt.Errorf("target must be a struct but got %v", targetType.Kind())
	}

	rawEntry := value["entry"]
	if entry, ok := rawEntry.(map[string]any); ok {
		copyFields(targetUnstructured, skipKeys(value, "entry"))
		entryPathInTarget := findEntryPathInTarget(targetType)
		dst := targetUnstructured
		if len(entryPathInTarget) > 0 {
			newValue := map[string]any{}
			if err = createField(targetUnstructured, newValue, entryPathInTarget...); err != nil {
				return fmt.Errorf("failed to set target copy destination to path %v: %w", entryPathInTarget, err)
			}
			dst = newValue
		}
		copyFields(dst, entry)
	} else {
		copyFields(targetUnstructured, value)
	}
	delete(targetUnstructured, "groupref")
	if err := fromUnstructured(target, targetUnstructured); err != nil {
		return fmt.Errorf("failed to set structured value from unstructured: %w", err)
	}
	return nil
}

func (t *Translator) processKubeMappings(obj map[string]any) ([]client.Object, error) {
	mappingsYML := t.crd.definition.ObjectMeta.Annotations[APIMAppingsAnnotation]
	if mappingsYML == "" {
		return []client.Object{}, nil
	}
	mappings := map[string]any{}
	yaml.Unmarshal([]byte(mappingsYML), mappings)

	deps := depsBuilder{}
	if err := t.processKubeMappingsAt(obj, mappings, &deps, "spec"); err != nil {
		return nil, fmt.Errorf("failed to map properties of spec from API to Kubernetes: %w", err)
	}
	if err := t.processKubeMappingsAt(obj, mappings, &deps, "status"); err != nil {
		return nil, fmt.Errorf("failed to map properties of status from API to Kubernetes: %w", err)
	}
	return deps.dependencies(), nil
}

func (t *Translator) processKubeMappingsAt(obj, mappings map[string]any, deps *depsBuilder, fieldName string) error {
	props, err := accessField[map[string]any](mappings,
		"properties", fieldName, "properties", t.sdk.version, "properties")
	if errors.Is(err, ErrNotFound) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("failed to access the API mapping properties for the %s: %w", fieldName, err)
	}
	field, err := accessField[map[string]any](obj, fieldName, t.sdk.version)
	if err != nil {
		return fmt.Errorf("failed to access object's %s: %w", fieldName, err)
	}
	if err := processKubeProperties([]string{}, props, field, deps); err != nil {
		return fmt.Errorf("failed to process properties from API into %s: %w", fieldName, err)
	}
	return nil
}

func (t *Translator) processAPIMappings(spec map[string]any) error {
	mappingsYML := t.crd.definition.ObjectMeta.Annotations[APIMAppingsAnnotation]
	if mappingsYML == "" {
		return nil
	}
	mappings := map[string]any{}
	yaml.Unmarshal([]byte(mappingsYML), mappings)
	props, err := accessField[map[string]any](mappings,
		"properties", "spec", "properties", t.sdk.version, "properties")
	if errors.Is(err, ErrNotFound) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("failed to access the API mapping properties for the spec: %w", err)
	}
	return processAPIProperties([]string{}, props, spec, t.deps)
}

func findEntryPathInTarget(targetType reflect.Type) []string {
	if targetType.String() == "admin.CreateAlertConfigurationApiParams" {
		return []string{"GroupAlertsConfig"}
	}
	return []string{}
}
