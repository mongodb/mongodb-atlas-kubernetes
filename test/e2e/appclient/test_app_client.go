//go:build e2e

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
		return ""
	}
	defer res.Body.Close()
	data, _ := io.ReadAll(res.Body)
	fmt.Print(string(data))
	return string(data)
}
