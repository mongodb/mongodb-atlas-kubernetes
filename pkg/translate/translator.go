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

type DependencyRepo interface {
	MainObject() client.Object
	Find(name, namespace string) client.Object
	Add(obj client.Object)
	Added() []client.Object
}

type Dependencies struct {
	mainObj client.Object
	deps    map[string]client.Object
	added   []client.Object
}

// NewDependencies creates a set of Kubernetes client.Objects
func NewDependencies(mainObj client.Object, objs ...client.Object) *Dependencies {
	deps := map[string]client.Object{}
	for _, obj := range objs {
		deps[client.ObjectKeyFromObject(obj).String()] = obj
	}
	return &Dependencies{
		mainObj: mainObj,
		deps:    deps,
		added:   []client.Object{},
	}
}


// MainObject retried the main object for this dependecny repository
func (d *Dependencies) MainObject() client.Object {
	return d.mainObj
}

// Find looks for an object withing the dependencies by name and namespace
func (d *Dependencies) Find(name, namespace string) client.Object {
	ns := namespace
	if ns == SetFallbackNamespace {
		ns = d.mainObj.GetNamespace()
	}
	return d.deps[client.ObjectKey{Name: name, Namespace: ns}.String()]
}

// Add appends an object to the added list and records it in the general set
func (d *Dependencies) Add(obj client.Object) {
	if obj.GetNamespace() == SetFallbackNamespace {
		obj.SetNamespace(d.mainObj.GetNamespace())
	}
	d.deps[client.ObjectKeyFromObject(obj).String()] = obj
	for i := range d.added {
		if d.added[i].GetName() == obj.GetName() && d.added[i].GetNamespace() == obj.GetNamespace() {
			d.added[i] = obj
			return
		}
	}
	d.added = append(d.added, obj)
}

// Added dumps an array of all dependencies added to the set after creation
func (d *Dependencies) Added() []client.Object {
	return d.added
}

// Translator allows to translate back and forth between a CRD schema version
// and SDK API structures of a certain version
type Translator struct {
	crd  CRDInfo
	sdk  SDKInfo
	deps DependencyRepo
}

// SDKInfo holds the SDK version information
type SDKInfo struct {
	version string
}

// NewTranslator creates a translator for a particular CRD and SDK version pairs,
// and with a particular set of known Kubernetes dependencies
func NewTranslator(crd *apiextensionsv1.CustomResourceDefinition, crdVersion string, sdkVersion string, deps DependencyRepo) *Translator {
	return &Translator{
		crd:  CRDInfo{definition: crd, version: crdVersion},
		sdk:  SDKInfo{version: sdkVersion},
		deps: deps,
	}
}

// PtrClientObj is a pointer type implementing client.Object
type PtrClientObj[T any] interface {
	*T
	client.Object
}

// FromAPI translaters a source API structure into a Kubernetes object
func FromAPI[S any, T any, P PtrClientObj[T]](t *Translator, target P, source *S) ([]client.Object, error) {
	sourceUnstructured, err := toUnstructured(source)
	if err != nil {
		return nil, fmt.Errorf("failed to convert API source value to unstructured: %w", err)
	}

	targetUnstructured := map[string]any{}

	versionedSpec := map[string]any{}
	copyFields(versionedSpec, sourceUnstructured)
	if err := createField(targetUnstructured, versionedSpec, "spec", t.sdk.version); err != nil {
		return nil, fmt.Errorf("failed to create versioned spec object in unstructured target: %w", err)
	}
	versionedSpecEntry := map[string]any{}
	copyFields(versionedSpecEntry, sourceUnstructured)
	versionedSpec["entry"] = versionedSpecEntry

	versionedStatus := map[string]any{}
	copyFields(versionedStatus, sourceUnstructured)
	if err := createField(targetUnstructured, versionedStatus, "status", t.sdk.version); err != nil {
		return nil, fmt.Errorf("failed to create versioned status object in unsstructured target: %w", err)
	}

	extraObjects, err := t.expandMappings(targetUnstructured)
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

	if err := t.collapseMappings(value); err != nil {
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

func (t *Translator) expandMappings(obj map[string]any) ([]client.Object, error) {
	mappingsYML := t.crd.definition.ObjectMeta.Annotations[APIMAppingsAnnotation]
	if mappingsYML == "" {
		return []client.Object{}, nil
	}
	mappings := map[string]any{}
	yaml.Unmarshal([]byte(mappingsYML), mappings)

	if err := t.expandMappingsAt(obj, mappings, "spec", t.sdk.version); err != nil {
		return nil, fmt.Errorf("failed to map properties of spec from API to Kubernetes: %w", err)
	}
	if err := t.expandMappingsAt(obj, mappings, "spec", t.sdk.version, "entry"); err != nil {
		return nil, fmt.Errorf("failed to map properties of spec from API to Kubernetes: %w", err)
	}
	if err := t.expandMappingsAt(obj, mappings, "status", t.sdk.version); err != nil {
		return nil, fmt.Errorf("failed to map properties of status from API to Kubernetes: %w", err)
	}
	return t.deps.Added(), nil
}

func (t *Translator) expandMappingsAt(obj, mappings map[string]any, fields ...string) error {
	expandedPath := []string{"properties"}
	for _, field := range fields {
		expandedPath = append(expandedPath, field, "properties")
	}
	props, err := accessField[map[string]any](mappings, expandedPath...)
	if errors.Is(err, ErrNotFound) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("failed to access the API mapping properties for %v: %w", expandedPath, err)
	}
	field, err := accessField[map[string]any](obj, fields...)
	if err != nil {
		return fmt.Errorf("failed to access object's %v: %w", fields, err)
	}
	mapper := Mapper{deps: t.deps, expand: true}
	if err := mapper.mapProperties([]string{}, props, field); err != nil {
		return fmt.Errorf("failed to process properties from API into %v: %w", fields, err)
	}
	return nil
}

func (t *Translator) collapseMappings(spec map[string]any) error {
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
	mapper := Mapper{deps: t.deps, expand: false}
	return mapper.mapProperties([]string{}, props, spec)
}

func findEntryPathInTarget(targetType reflect.Type) []string {
	if targetType.String() == "admin.CreateAlertConfigurationApiParams" {
		return []string{"GroupAlertsConfig"}
	}
	return []string{}
}
