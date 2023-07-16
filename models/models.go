package models

import (
	"fmt"
	"os"
	"sync"
	// "time"

	// "github.com/mitchellh/go-ps"
	"github.com/waxdred/GoHotReload/watcher"
)

const PROCESSLEN = 16

type ChanProg struct {
	kill   chan bool
	stdout chan string
	stderr chan string
}

type ChanApp struct {
	Call       chan bool
	signalChan chan os.Signal
}

type App struct {
	Program      Program
	Mu           sync.Mutex
	ConfigSelect string
	Chan         ChanApp
	config       int
	Config       []Config `yaml:"configs"`
	model        *model
	error        error
}

type Configs struct {
	Configs []Config `yaml:"configs"`
}

type Config struct {
	Name       string   `yaml:"name"`
	Cmd        []string `yaml:"cmd"`
	Executable string   `yaml:"executable"`
	Extension  string   `yaml:"extension"`
	Interval   int      `yaml:"interval"`
	Path       string   `yaml:"path"`
}

type Program struct {
	Pid     int
	Process *os.Process
	Config  *Config `yaml:"configs"`
	TTY     string
	check   bool
	process bool
	restart bool
	info    string
	Chan    ChanProg
}

func NewProg(config *Config) *Program {
	prog := &Program{
		Config: config,
	}
	prog.Chan.kill = make(chan bool, 1)
	return prog
}

func (app *App) Listen() *App {
	// ticker := time.NewTicker(2 * time.Second)
	go func() {
		for {
			select {
			case <-app.Chan.Call:
				app.Mu.Lock()
				app.printBox(&app.Program)
				app.Mu.Unlock()
			case err := <-app.Program.Chan.stderr:
				// TODO storage string in var for use in viewport
				fmt.Println(err)
			case out := <-app.Program.Chan.stdout:
				// TODO storage string in var for use in viewport
				fmt.Println("output:", out)
			}
		}
	}()
	return app
}

func (app *App) Start() *App {
	app.Listen()
	err := app.checkPath().error
	watcher := watcher.NewWatcher().
		NewPath(app.Program.Config.Path).
		NewExtension(app.Program.Config.Extension).
		NewExecutable(app.Program.Config.Executable)
	if err != nil {
		return app
	}
	if err != nil {
		app.error = err
		return app
	}
	var wg sync.WaitGroup
	go HandlerSig(app)
	app.Chan.Call <- true

	app.Program.Chan.stderr = make(chan string)
	app.Program.Chan.stdout = make(chan string)
	defer close(app.Program.Chan.stdout)
	defer close(app.Program.Chan.stderr)
	app.Program.check = true
	if err := execCmd(&app.Program); err != nil {
		return app
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		app.Program.info = "Search Programm..."
		for {
			select {
			case Event := <-watcher.Event:
				fmt.Print(Event.Name)
				app.Program.info = "Kill program"
				app.Program.restart = true
				app.Program.process = false
				app.Program.check = false
				if app.handlerProcess(&app.Program) {
					killPid(&app.Program)
				} else {
					if err := execCmd(&app.Program); err != nil {
						return
					}
				}
				app.Program.Pid = 0
				app.Chan.Call <- true
			case err := <-watcher.Errors:
				fmt.Println("Error:", err)
			case <-app.Program.Chan.kill:
				app.Program.info = "Restart program"
				if err := execCmd(&app.Program); err != nil {
					return
				}
				app.Program.restart = false
				app.Program.process = true
				app.Program.check = false
				app.Chan.Call <- true
			case pid := <-watcher.EventPid:
				if app.Program.Pid == 0 {
					app.Program.info = fmt.Sprintf("%s: PID found: %d\n", app.Program.Config.Executable, pid.Pid)
					app.Program.Pid = pid.Pid
					app.Program.process = true
					app.Program.restart = false
					app.Program.check = false
					app.Chan.Call <- true
				}
			default:
				break
			}
		}
	}()
	wg.Wait()
	return app
}
