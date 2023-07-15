package models

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/charmbracelet/lipgloss"
	"gopkg.in/yaml.v2"
)

func New() *App {
	app := &App{}
	app.Mu.Lock()
	defer app.Mu.Unlock()
	clearScreen()
	executablePath, err := os.Executable()
	if err != nil {
		fmt.Println("Error retrieving executable path:", err)
		return app
	}
	binaryPath := filepath.Dir(executablePath)
	yamlFile, err := ioutil.ReadFile(binaryPath + "/config/config.yml")
	if err != nil {
		log.Fatalf("Error reading YAML file: %v", err)
	}
	var config Configs
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		log.Fatalf("Error deserializing YAML file: %v", err)
	}
	app.Config = config.Configs
	if len(app.Config) > 1 {
		fmt.Println(Style.Foreground(lipgloss.Color(blue)).Margin(1, 1).Render("Select your config"))
		app.OptionsView()
		for i := range app.Config {
			conf := app.Config[i]
			if conf.Name == app.ConfigSelect {
				fmt.Println("ok", i)
				app.Program = append(app.Program, *NewProg(&conf))
			}
		}
	} else if len(app.Config) != 0 {
		app.Program = append(app.Program, *NewProg(&app.Config[0]))
	} else {
		log.Fatal("Error")
	}
	clearScreen()
	return app
}

func (app *App) Errror() error {
	return app.error
}
