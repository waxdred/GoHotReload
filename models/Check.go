package models

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

func (app *App) checkPath() *App {
	for _, prog := range app.Program {
		_, err := os.Stat(prog.Config.Path)
		if err != nil {
			app.error = err
			break
		}
		if os.IsNotExist(err) {
			app.error = err
			break
		}
	}
	return app
}

func (app *App) CheckingParse() *App {
	for i := range app.Program {
		prog := &app.Program[i]
		if prog.Config.Cmd == "" {
			app.error = errors.New("Please add command")
			break
		} else if prog.Config.Interval == 0 {
			prog.Config.Interval = 4
		} else if prog.Config.Extension == "" {
			prog.Config.Extension = ".go"
		} else if prog.Config.Path == "" {
			prog.Config.Path = "./"
		}
		ok := strings.HasPrefix(prog.Config.Extension, ".")
		if !ok {
			app.error = errors.New("Please use correct extenstion:")
			break
		}
		ok = strings.HasPrefix(prog.Config.Path, "~")
		if ok {
			homeDire, err := os.UserHomeDir()
			if err != nil {
				app.error = err
				break
			}
			prog.Config.Path = strings.TrimPrefix(prog.Config.Path, "~")
			prog.Config.Path = fmt.Sprint(homeDire, prog.Config.Path)
		}
	}
	return app
}
