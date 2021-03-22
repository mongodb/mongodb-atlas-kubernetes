package app

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type App struct {
	uri string
}

func NewApp(uri string) *App {
	return &App{
		uri: uri,
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
	data, _ := ioutil.ReadAll(res.Body)
	fmt.Print(string(data))
	return string(data)
}
