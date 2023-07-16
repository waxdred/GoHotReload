package models

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/mitchellh/go-ps"
)

const PROCESSLEN = 16

type ChanProg struct {
	reload chan bool
	pid    chan bool
	stdout chan string
	stderr chan string
}

type ChanApp struct {
	Call   chan bool
	Notify chan bool
}

type App struct {
	Program      []Program
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
	Files   map[string]time.Time
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
	return prog
}

func (app *App) Listen() *App {
	// ticker := time.NewTicker(2 * time.Second)
	go func() {
		for {
			select {
			case <-app.Chan.Call:
				app.Mu.Lock()
				if len(app.Program) > 0 {
					app.printBox(&app.Program[0])
				}
				app.Mu.Unlock()
			case err := <-app.Program[0].Chan.stderr:
				// TODO storage string in var for use in viewport
				fmt.Println(err)
			case out := <-app.Program[0].Chan.stdout:
				// TODO storage string in var for use in viewport
				fmt.Println("output:", out)
			}
		}
	}()
	return app
}

func (app *App) Start() *App {
	err := app.checkPath().error
	if err != nil {
		return app
	}

	var wg sync.WaitGroup
	go HandlerSig(app)
	app.Chan.Call <- true

	for i := range app.Program {
		prog := &app.Program[i]
		prog.Chan.reload = make(chan bool)
		prog.Chan.stderr = make(chan string)
		prog.Chan.stdout = make(chan string)
		defer close(prog.Chan.reload)
		defer close(prog.Chan.stdout)
		defer close(prog.Chan.stderr)
		prog.check = true
		if err := execCmd(prog); err != nil {
			return app
		}
		wg.Add(1)
		go func() {
			pid := make(chan bool, 1)
			defer wg.Done()
			ticker := time.NewTicker(time.Duration(prog.Config.Interval) * time.Second)
			routine := false
			for {
				prog.info = "Search Programm..."
				select {
				case <-pid:
					if routine {
						if app.handlerProcess(prog) {
							app.process(prog)
						}
						pid <- false
						prog.info = ""
						prog.Pid = 0
						routine = false
						prog.check = true
					}
				case <-ticker.C:
					if !routine {
						processes, err := ps.Processes()
						if err != nil {
							fmt.Println("Error:", err)
							os.Exit(1)
						}
						for _, process := range processes {
							tmp := ""
							if len(prog.Config.Executable) > PROCESSLEN {
								tmp = prog.Config.Executable[:PROCESSLEN]
							} else {
								tmp = prog.Config.Executable
							}
							if process.Executable() == tmp {
								prog.info = fmt.Sprintf("%s: PID found: %d\n", prog.Config.Executable, process.Pid())
								prog.Pid = process.Pid()
								app.Chan.Call <- true
								app.execPs(prog)
								routine = true
								prog.check = false
								pid <- true
								break
							}
						}
					}
				}
			}
		}()
	}
	wg.Wait()
	return app
}
