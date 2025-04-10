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
		fmt.Println(text.FgRed.Sprintf("%v", err.Error()))

		return
	}

	azureCleaner, err := provider.NewAzureCleaner()
	if err != nil {
		fmt.Println(text.FgRed.Sprintf("%v", err.Error()))

		return
	}

	c, err := atlas.NewCleaner(awsCleaner, gcpCleaner, azureCleaner)
	if err != nil {
		fmt.Println(text.FgRed.Sprintf("%v", err.Error()))

		return
	}

	lifetimeHours, err := strconv.Atoi(os.Getenv("PROJECT_LIFETIME"))
	if err != nil {
		err = fmt.Errorf("error parsing PROJECT_LIFETIME environment variable: %w", err)
		fmt.Println(text.FgRed.Sprintf("%v", err.Error()))

		return
	}

	err = c.Clean(ctx, lifetimeHours)
	if err != nil {
		fmt.Println(text.FgRed.Sprintf("%v", err.Error()))
	}
}
