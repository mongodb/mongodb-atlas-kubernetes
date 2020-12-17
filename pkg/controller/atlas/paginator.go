package atlas

import (
	"reflect"

	"go.mongodb.org/atlas/mongodbatlas"
)

// Paginated is the general interface for a single page returned by Atlas api.
type Paginated interface {
	// HasNext returns if there are more pages
	HasNext() bool

	// Results returns the list of entities on a single page
	Results() []interface{}
}

type atlasPaginated struct {
	links    []*mongodbatlas.Link
	entities []interface{}
}

func NewAtlasPaginated(response *mongodbatlas.Response, entities interface{}) Paginated {
	return &atlasPaginated{links: response.Links, entities: toSlice(entities)}
}

func (p atlasPaginated) Results() []interface{} {
	return p.entities
}

// HasNext return true if there is next page (see 'ApiBaseResource.handlePaginationInternal` in mms code)
func (p atlasPaginated) HasNext() bool {
	for _, l := range p.links {
		if l.Rel == "next" {
			return true
		}
	}
	return false
}

func DefaultListOptions(pageNum int) *mongodbatlas.ListOptions {
	return &mongodbatlas.ListOptions{
		PageNum:      pageNum,
		ItemsPerPage: 500,
	}
}

// PageReader is the function that reads a single page by its number
type PageReader func(pageNum int) (Paginated, error)

// PageItemPredicate is the function that processes single item on the page and returns true if no further processing
// needs to be done (usually it's the search logic)
type PageItemPredicate func(entity interface{}) bool

// TraversePages reads page after page using 'reader' and applies the 'predicate' for each item on the page.
// Stops traversal when the 'predicate' returns true.
func TraversePages(reader PageReader, predicate PageItemPredicate) error {
	// Let's be safe and not get into infinite loop in case something goes wrong with the links
	for i := 1; i <= 500; i++ {
		paginated, e := reader(i)
		if e != nil {
			return e
		}
		for _, entity := range paginated.Results() {
			if predicate(entity) {
				return nil
			}
		}
		if !paginated.HasNext() {
			return nil
		}
	}
	return nil
}

func toSlice(data interface{}) []interface{} {
	value := reflect.ValueOf(data)

	result := make([]interface{}, value.Len())
	for i := 0; i < value.Len(); i++ {
		result[i] = value.Index(i).Interface()
	}
	return result
}
