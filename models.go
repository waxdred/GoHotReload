package main

import (
	"fmt"
	"os"
	"path"
	"time"

	"github.com/mitchellh/go-ps"
)

type App struct {
	Pid        int
	Executable string
	Process    *os.Process
	Path       string
	Extension  string
	Cmd        string
	Interval   time.Duration
	Files      map[string]time.Time
	error      error
}

func New() *App {
	return &App{}
}

func (app *App) Errror() error {
	return app.error
}

func (app *App) checkPath() *App {
	_, err := os.Stat(app.Path)
	if err != nil {
		app.error = err
	}
	if os.IsNotExist(err) {
		fmt.Println("Path not found")
		app.error = err
	}
	return app
}

func (app *App) handlerProcess() *App {
	app.Process, app.error = os.FindProcess(app.Pid)
	fmt.Println(app.Process)
	if app.error != nil {
		fmt.Println("Can't found Process")
		return app
	}
	return app
}

func (app *App) Start() *App {
	err := app.checkPath().error
	if err != nil {
		return app
	}
	app.Executable = path.Base(app.Path)
	pid := make(chan bool)
	ticker := time.NewTicker(app.Interval)

	go func() {
		for {
			select {
			case <-pid:
				return
			case <-ticker.C:
				fmt.Println("Checking for PID...")
				processes, err := ps.Processes()
				if err != nil {
					fmt.Println("Error:", err)
					os.Exit(1)
				}
				for _, process := range processes {
					if process.Executable() == app.Executable {
						fmt.Printf("PID found: %d\n", process.Pid())
						app.Pid = process.Pid()
						pid <- true
						return
					}
				}
			}
		}
	}()
	<-pid
	app.handlerProcess()
	return app
}
