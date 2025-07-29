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

type SDKInfo struct {
	version string
}

func NewTranslator(crd *apiextensionsv1.CustomResourceDefinition, crdVersion string, sdkVersion string, deps DependencyFinder) *Translator {
	return &Translator{
		crd:  CRDInfo{definition: crd, version: crdVersion},
		sdk:  SDKInfo{version: sdkVersion},
		deps: deps,
	}
}

func ToAPI[T any](t *Translator, target *T, source client.Object) error {
	targetType := reflect.TypeOf(target).Elem()
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
		return fmt.Errorf("failed to convert source value to unstructured: %w", err)
	}
	targetUnstructured := map[string]any{}
	value, err := accessField[map[string]any](unstructuredSrc, "spec", t.sdk.version)
	if err != nil {
		return fmt.Errorf("failed to access source spec value: %w", err)
	}

	if err := t.processMappings(value); err != nil {
		return fmt.Errorf("failed to process API mappings: %w", err)
	}

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

func (t *Translator) processMappings(spec map[string]any) error {
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
	return processProperties([]string{}, props, spec, t.deps)
}

func findEntryPathInTarget(targetType reflect.Type) []string {
	if targetType.String() == "admin.CreateAlertConfigurationApiParams" {
		return []string{"GroupAlertsConfig"}
	}
	return []string{}
}
