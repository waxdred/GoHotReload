package models

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/charmbracelet/lipgloss"
	"gopkg.in/yaml.v2"
)

func New() *App {
	app := &App{}
	app.Mu.Lock()
	app.Chan.signalChan = make(chan os.Signal, 1)
	signal.Notify(app.Chan.signalChan, syscall.SIGHUP)
	defer app.Mu.Unlock()
	app.Chan.Call = make(chan bool, 1)
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
	for i := range app.Config {
		conf := app.Config[i]
		if conf.Name == app.ConfigSelect {
			app.Program = *NewProg(&conf)
		}
	}
	clearScreen()
	return app
}

func (app *App) Errror() error {
	return app.error
}
