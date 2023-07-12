package models

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

func (app *App) checkPath() *App {
	for _, prog := range app.Program {
		_, err := os.Stat(prog.Path)
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
		if prog.Cmd == "" {
			app.error = errors.New("Please add command")
			break
		} else if prog.Interval == 0 {
			prog.Interval = 4
		} else if prog.Extension == "" {
			prog.Extension = ".go"
		} else if prog.Path == "" {
			prog.Path = "./"
		}
		ok := strings.HasPrefix(prog.Extension, ".")
		if !ok {
			app.error = errors.New("Please use correct extenstion:")
			break
		}
		ok = strings.HasPrefix(prog.Path, "~")
		if ok {
			homeDire, err := os.UserHomeDir()
			if err != nil {
				app.error = err
				break
			}
			prog.Path = strings.TrimPrefix(prog.Path, "~")
			prog.Path = fmt.Sprint(homeDire, prog.Path)
		}
	}
	return app
}
