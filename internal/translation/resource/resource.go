// Package resource contains internal type for the Atlas Resource
package resource

// Resource is the internal representation of the Atlas Resource managed here
type Resource struct {
	// TODO: (DELETE ME) replace with the actual CRD struct and extra fields
	// akov2.AtlasResource
	// Some extra fields, like projectID or secrets,
	// anything needed but not directly attached in the CRD
}

// GetID sample
// TODO: (DELETE ME) replace with the actual CRD struct and extra fields
func (s *Resource) GetID() string {
	//return pointer.GetOrDefault(s.ID, "")
	panic("unimplemented")
}

// NewResource creates an internal type from the CRD struct and fields
func NewResource( /* resource *akov2.Resource, ...*/ ) *Resource {
	return &Resource{
		// AtlasResource:                resource,
		// ...
	}
}

// fromAtlas is an internal conversion function from atlas to internal type
func fromAtlas( /*resource admin.Resource*/ ) (*Resource, error) {
	/*
		...
		return &Resource{
			// ...
		}, errors.Join(errs...)
	*/
	panic("unimplemented")
}

// Normalize prepares Resource for internal type  comparisons to work flawlessly
func (s *Resource) Normalize() (*Resource, error) {
	panic("unimplemented")
}

// toAtlas is an internal conversion to get the atlas type equivalent value
func (s *Resource) toAtlas() /* *admin.Resource, */ error {
	/*
		...
		return &admin.Resource{
			// ...
		}, errors.Join(errs...)
	*/
	panic("unimplemented")
}
