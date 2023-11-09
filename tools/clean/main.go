package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"tools/clean/atlas"
	"tools/clean/provider"

	"github.com/jedib0t/go-pretty/v6/text"
)

func main() {
	ctx := context.Background()
	awsCleaner := provider.NewAWSCleaner()

	gcpCleaner, err := provider.NewGCPCleaner(ctx)
	if err != nil {
		fmt.Println(text.FgRed.Sprintf(err.Error()))

		return
	}

	azureCleaner, err := provider.NewAzureCleaner()
	if err != nil {
		fmt.Println(text.FgRed.Sprintf(err.Error()))

		return
	}

	c, err := atlas.NewCleaner(awsCleaner, gcpCleaner, azureCleaner)
	if err != nil {
		fmt.Println(text.FgRed.Sprintf(err.Error()))

		return
	}

	lifetime, err := strconv.Atoi(os.Getenv("PROJECT_LIFETIME"))
	if err != nil {
		fmt.Println(text.FgRed.Sprintf(err.Error()))

		return
	}

	err = c.Clean(ctx, lifetime)
	if err != nil {
		fmt.Println(text.FgRed.Sprintf(err.Error()))
	}
}
