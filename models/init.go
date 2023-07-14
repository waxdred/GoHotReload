package models

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
)

func New() *App {
	app := &App{}
	app.Mu.Lock()
	defer app.Mu.Unlock()
	clearScreen()
	fmt.Println("Parsing config ...")
	executablePath, err := os.Executable()
	if err != nil {
		fmt.Println("Erreur lors de la récupération du chemin du binaire :", err)
		return app
	}
	binaryPath := filepath.Dir(executablePath)
	dirfile, err := ioutil.ReadDir(binaryPath + "/config/")
	if err != nil {
		app.error = err
		return app
	}
	for _, file := range dirfile {
		extend := path.Ext(file.Name())
		if extend == ".json" {
			app.configFile = append(app.configFile, file.Name())
		}
	}
	if len(app.configFile) > 1 {
		app.OptionsView()
	}
	clearScreen()
	if app.ConfigSelect == "" && len(app.configFile) == 1 {
		app.ConfigSelect = app.configFile[0]
	}
	data, err := ioutil.ReadFile(binaryPath + "/config/" + app.ConfigSelect)
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
