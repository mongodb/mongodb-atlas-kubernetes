package atlas

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/atlas/mongodbatlas"
)

func Test_TraversePages(t *testing.T) {
	t.Run("Project Found (single page)", func(t *testing.T) {
		found := false
		iterations := 0
		pagesScanned := 0
		err := TraversePages(organizationPages(1, &pagesScanned), projectFound("test3", &found, &iterations))
		assert.True(t, found)
		assert.NoError(t, err)
		assert.Equal(t, 4, iterations)
		assert.Equal(t, 1, pagesScanned)
	})
	t.Run("Project Not Found (single page)", func(t *testing.T) {
		found := false
		iterations := 0
		pagesScanned := 0
		err := TraversePages(organizationPages(1, &pagesScanned), projectFound("fake", &found, &iterations))
		assert.False(t, found)
		assert.NoError(t, err)
		assert.Equal(t, 500, iterations)
		assert.Equal(t, 1, pagesScanned)
	})
	t.Run("Project Found (multiple pages)", func(t *testing.T) {
		found := false
		iterations := 0
		pagesScanned := 0
		err := TraversePages(organizationPages(3, &pagesScanned), projectFound("test600", &found, &iterations))
		assert.True(t, found)
		assert.NoError(t, err)
		assert.Equal(t, 601, iterations)
		assert.Equal(t, 2, pagesScanned)
	})
	t.Run("Project Not Found (multiple pages)", func(t *testing.T) {
		found := false
		iterations := 0
		pagesScanned := 0
		err := TraversePages(organizationPages(3, &pagesScanned), projectFound("fake", &found, &iterations))
		assert.False(t, found)
		assert.NoError(t, err)
		assert.Equal(t, 1500, iterations)
		assert.Equal(t, 3, pagesScanned)
	})
	t.Run("Error happened", func(t *testing.T) {
		err := TraversePages(func(pageNum int) (Paginated, error) { return nil, errors.New("Error!") }, nil)
		assert.Error(t, err)
	})
}

/*func TestPagination_MultiplePages(t *testing.T) {
	found, err := TraversePages(multipleOrganizationsPage, func(obj interface{}) bool { return obj.(*Organization).Name == "test1220" })
	assert.True(t, found)
	assert.NoError(t, err)
	assert.Equal(t, 3, numberOfPagesTraversed)

	found, err = TraversePages(multipleOrganizationsPage, func(obj interface{}) bool { return obj.(*Organization).Name == "test1400" })
	assert.False(t, found)
	assert.NoError(t, err)
}

func TestPagination_Error(t *testing.T) {
	_, err := TraversePages(func(pageNum int) (Paginated, error) { return nil, errors.New("Error!") },
		func(obj interface{}) bool { return obj.(*Organization).Name == "test1220" })
	assert.Errorf(t, err, "Error!")
}*/

func organizationPages(totalPages int, pagesScanned *int) func(pageNum int) (Paginated, error) {
	return func(pageNum int) (Paginated, error) {
		*pagesScanned++
		links := []*mongodbatlas.Link{{Rel: "next"}}
		if pageNum == totalPages {
			links = []*mongodbatlas.Link{}
		}
		return &atlasPaginated{
			links:    links,
			entities: generateProjects((pageNum-1)*500, 500),
		}, nil
	}
}

var singleOrganizationsPage = func(pageNum int) (Paginated, error) {
	if pageNum == 1 {
		// Note, that we don't specify 'next' attribute, so no extra pages will be requested
		return &atlasPaginated{
			links:    []*mongodbatlas.Link{},
			entities: generateProjects(0, 100),
		}, nil
	}
	return nil, errors.New("Not found!")
}

var numberOfPagesTraversed = 0

var multipleOrganizationsPage = func(pageNum int) (Paginated, error) {
	numberOfPagesTraversed++
	// page 1
	if pageNum == 1 {
		return &atlasPaginated{
			links:    []*mongodbatlas.Link{{Rel: "next"}},
			entities: generateProjects(0, 500),
		}, nil
	} else if pageNum == 2 {
		return &atlasPaginated{
			links:    []*mongodbatlas.Link{{Rel: "next"}},
			entities: generateProjects(500, 1000),
		}, nil
	} else if pageNum == 3 {
		return &atlasPaginated{
			links:    []*mongodbatlas.Link{},
			entities: generateProjects(1000, 300),
		}, nil
	}
	return nil, errors.New("Not found!")
}

func generateProjects(startFrom, count int) []interface{} {
	ans := make([]interface{}, count)
	c := startFrom
	for i := 0; i < count; i++ {
		ans[i] = &mongodbatlas.Project{ID: fmt.Sprintf("id%d", c), Name: fmt.Sprintf("test%d", c)}
		c++
	}
	return ans
}

func projectFound(name string, found *bool, iterations *int) func(obj interface{}) bool {
	return func(obj interface{}) bool {
		*iterations++
		if obj.(*mongodbatlas.Project).Name == name {
			*found = true
			return true
		}
		return false
	}
}
