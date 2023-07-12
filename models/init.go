package models

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

func New() *App {
	app := &App{}
	app.Mu.Lock()
	defer app.Mu.Unlock()
	fmt.Println("Parsing config ...")
	data, err := ioutil.ReadFile("./config.json")
	if err != nil {
		app.error = err
		return app
	}
	err = json.Unmarshal(data, &app.Program)
	if err != nil {
		app.error = err
		return app
	}
	return app
}

func (app *App) Errror() error {
	return app.error
}
