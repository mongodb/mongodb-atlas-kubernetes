package search

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/contract"
)

const (
	TestName = "search"
)

var (
	// WipeResources removes resources at test cleanup time, no reuse will be possible afterwards
	WipeResources = contract.BoolEnv("WIPE_RESOURCES", false)
)

var resources *contract.Resources

func TestMain(m *testing.M) {
	ctx := context.Background()
	beforeAll(ctx)
	code := m.Run()
	afterAll(ctx)
	os.Exit(code)
}

func beforeAll(ctx context.Context) {
	log.Printf("WipeResources set to %v", WipeResources)
	resources = contract.MustDeployTestResources(ctx,
		TestName,
		WipeResources,
		contract.DefaultProject(TestName),
		contract.WithServerless(contract.DefaultServerless(TestName)))
	log.Printf("Resources ready:\n%v", resources)
}

func afterAll(ctx context.Context) {
	resources.MustRecycle(ctx, WipeResources)
}

func TestCreateSearchIndex(t *testing.T) {
	fmt.Printf("yay!\n")
}
