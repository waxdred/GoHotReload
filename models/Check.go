package models

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

func (app *App) checkPath() *App {
	_, err := os.Stat(app.Program.Config.Path)
	if err != nil {
		app.error = err
		return app
	}
	if os.IsNotExist(err) {
		app.error = err
		return app
	}
	return app
}

func (app *App) CheckingParse() *App {
	if app.Program.Config == nil {
		os.Exit(0)
	}
	if app.Program.Config.Cmd[0] == "" {
		app.error = errors.New("Please add command")
		return app
	} else if app.Program.Config.Interval == 0 {
		app.Program.Config.Interval = 4
	} else if app.Program.Config.Extension == "" {
		app.Program.Config.Extension = ".go"
	} else if app.Program.Config.Path == "" {
		app.Program.Config.Path = "./"
	}
	ok := strings.HasPrefix(app.Program.Config.Extension, ".")
	if !ok {
		app.error = errors.New("Please use correct extenstion:")
		return app
	}
	ok = strings.HasPrefix(app.Program.Config.Path, "~")
	if ok {
		homeDire, err := os.UserHomeDir()
		if err != nil {
			app.error = err
			return app
		}
		app.Program.Config.Path = strings.TrimPrefix(app.Program.Config.Path, "~")
		app.Program.Config.Path = fmt.Sprint(homeDire, app.Program.Config.Path)
	}
	return app
}
