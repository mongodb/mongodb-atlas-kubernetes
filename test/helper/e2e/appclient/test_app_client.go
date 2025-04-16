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

package app

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

type App struct {
	uri string
}

func NewTestAppClient(port string) *App {
	return &App{
		uri: "http://localhost:" + port,
	}
}

func (app *App) Post(json string) error {
	res, err := http.Post(
		app.uri+"/mongo/",
		"application/json;",
		strings.NewReader(json),
	)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	return nil
}

func (app *App) Get(endpoint string) string {
	res, err := http.Get(app.uri + endpoint)
	if err != nil {
		fmt.Print(err)
		return err.Error()
	}
	defer res.Body.Close()
	data, _ := io.ReadAll(res.Body)
	fmt.Print(string(data))
	return string(data)
}
