// Copyright 2025 MongoDB Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/jedib0t/go-pretty/v6/text"
	"tools/clean/atlas"
	"tools/clean/provider"
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
