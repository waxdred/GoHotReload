package models

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func New() *App {
	app := &App{}
	app.Mu.Lock()
	defer app.Mu.Unlock()
	fmt.Println("Parsing config ...")
	executablePath, err := os.Executable()
	if err != nil {
		fmt.Println("Erreur lors de la récupération du chemin du binaire :", err)
		return app
	}
	binaryPath := filepath.Dir(executablePath)
	data, err := ioutil.ReadFile(binaryPath + "/config/config.json")
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
