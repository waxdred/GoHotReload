package models

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/mitchellh/go-ps"
)

const PROCESSLEN = 16

type App struct {
	Program      []Program
	Mu           sync.Mutex
	config       int
	configFile   []string
	ConfigSelect string
	model        *model
	error        error
}

func (app *App) Listen() *App {
	ticker := time.NewTicker(2 * time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				app.Mu.Lock()
				app.printBox()
				app.Mu.Unlock()
			}
		}
	}()
	return app
}

type Program struct {
	Pid        int
	Process    *os.Process
	Files      map[string]time.Time
	Path       string `json:"path"`
	Executable string `json:"executable"`
	Extension  string `json:"extension"`
	Cmd        string `json:"cmd"`
	Interval   int    `json:"interval"`
	TTY        string
	pid        chan bool
	check      bool
	process    bool
	restart    bool
	info       string
}

func (app *App) Start() *App {
	err := app.checkPath().error
	if err != nil {
		return app
	}

	var wg sync.WaitGroup
	for i := range app.Program {
		prog := &app.Program[i]
		prog.check = true
		wg.Add(1)
		go func() {
			pid := make(chan bool, 1)
			defer wg.Done()
			ticker := time.NewTicker(time.Duration(prog.Interval) * time.Second)
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
						prog.TTY = ""
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
							if len(prog.Executable) > PROCESSLEN {
								tmp = prog.Executable[:PROCESSLEN]
							} else {
								tmp = prog.Executable
							}
							if process.Executable() == tmp {
								prog.info = fmt.Sprintf("%s: PID found: %d\n", prog.Executable, process.Pid())
								prog.Pid = process.Pid()
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
