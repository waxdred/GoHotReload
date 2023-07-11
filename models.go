package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"path"
	"strings"
	"sync"
	"syscall"

	// "path"
	"time"

	"github.com/mitchellh/go-ps"
)

type App struct {
	Program []Program
	Mu      sync.Mutex
	error   error
}

type Program struct {
	Pid        int
	Process    *os.Process
	Files      map[string]time.Time
	Path       string `json:"path"`
	Executable string
	Extension  string `json:"extension"`
	Cmd        string `json:"cmd"`
	Interval   int    `json:"interval"`
	pid        chan bool
}

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

func (app *App) handlerProcess(prog *Program) *App {
	prog.Process, app.error = os.FindProcess(prog.Pid)
	if app.error != nil {
		fmt.Println("Can't found Process")
		os.Exit(-1)
	}
	return app
}

func (app *App) getExectutable() *App {
	for i := range app.Program {
		prog := &app.Program[i]
		if prog.Path == "." || prog.Path == "./" {
			pwd, _ := os.Getwd()
			prog.Executable = path.Base(pwd)
		} else {
			prog.Executable = path.Base(prog.Path)
		}
	}
	return app
}

func (app *App) process(prog *Program) {
	sigs := make(chan os.Signal, 1)
	ticker := time.NewTicker(time.Duration(prog.Interval) * time.Second)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)

	go func() {
		for {
			select {
			case <-sigs:
				fmt.Println("Thank see you next time...")
				return
			case <-ticker.C:
				fmt.Println("Running")
				update := Handler(prog)
				if update {
					KillPid(prog)
					return
				}
				fmt.Println("update:", update)
			}
		}
	}()
	<-sigs
}

func (app *App) Start() *App {
	err := app.checkPath().error
	if err != nil {
		return app
	}
	app.getExectutable()

	var wg sync.WaitGroup
	for i := range app.Program {
		prog := &app.Program[i]
		prog.pid = make(chan bool)
		wg.Add(1)
		go func() {
			defer wg.Done()
			ticker := time.NewTicker(time.Duration(prog.Interval) * time.Second)
			routine := false
			for {
				select {
				case <-time.After(30 * time.Second):
					fmt.Println("Timeout reached for", prog.Executable)
					return
				case <-prog.pid:
					fmt.Println("Start process")
					app.handlerProcess(prog)
					app.process(prog)
					fmt.Println("routine done need run again")
					routine = false
				case <-ticker.C:
					if !routine {
						fmt.Println("Checking", prog.Executable, "for PID...")
						processes, err := ps.Processes()
						if err != nil {
							fmt.Println("Error:", err)
							os.Exit(1)
						}
						for _, process := range processes {
							if process.Executable() == prog.Executable {
								fmt.Printf("%s: PID found: %d\n", prog.Executable, process.Pid())
								prog.Pid = process.Pid()
								routine = true
								break
							}
						}
						if routine {
							prog.pid <- true
						}
					}
				}
			}
		}()
	}

	wg.Wait()
	return app
}
