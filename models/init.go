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
	app.Call = make(chan bool)
	clearScreen()
	executablePath, err := os.Executable()
	if err != nil {
		fmt.Println("Erreur lors de la récupération du chemin du binaire :", err)
		return app
	}
	binaryPath := filepath.Dir(executablePath)
	yamlFile, err := ioutil.ReadFile(binaryPath + "/config/config.yml")
	if err != nil {
		log.Fatalf("Erreur lors de la lecture du fichier YAML : %v", err)
	}
	var config Configs
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		log.Fatalf("Erreur lors de la désérialisation du fichier YAML : %v", err)
	}
	app.Config = config.Configs
	if len(app.Config) > 1 {
		fmt.Println(Style.Foreground(lipgloss.Color(blue)).Margin(1, 1).Render("Select you config"))
		app.OptionsView()
	}
	for i := range app.Config {
		conf := app.Config[i]
		if conf.Name == app.ConfigSelect {
			fmt.Println("ok", i)
			app.Program = append(app.Program, *NewProg(&conf))
		}
	}
	clearScreen()
	return app
}

func (app *App) Errror() error {
	return app.error
}
