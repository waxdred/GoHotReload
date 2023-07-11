package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"
)

type Config struct {
	Cmd       string
	Interval  string
	Extension string
	Path      string
}

func PrintBox(app *App) {
	top := "┌───────────────────────────────────────────────────┐\n"
	bottom := "└───────────────────────────────────────────────────┘\n"
	cmd := fmt.Sprintf("│  cmd .......... %s", app.Cmd)
	extension := fmt.Sprintf("│  extension .... %s", app.Extension)
	path := fmt.Sprintf("│  path ......... %s", app.Path)
	pid := fmt.Sprint("│  pid ........... ", app.Pid)
	fmt.Printf(top)
	fmt.Printf("│                    GoHotRelaod                    │\n")
	fmt.Printf("│                                                   │\n")
	fmt.Printf(cmd)
	fmt.Printf(strings.Repeat(" ", (54 - len(cmd))))
	fmt.Printf("│\n")
	fmt.Printf(path)
	fmt.Printf(strings.Repeat(" ", (54 - len(path))))
	fmt.Printf("│\n")
	fmt.Printf(pid)
	fmt.Printf(strings.Repeat(" ", (54 - len(extension))))
	fmt.Printf("│\n")
	fmt.Printf(extension)
	fmt.Printf(strings.Repeat(" ", (54 - len(extension))))
	fmt.Printf("│\n")

	fmt.Println(bottom)
}

func Handler(app *App) bool {
	files := make(map[string]time.Time)
	update := false

	err := filepath.Walk(app.Path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && filepath.Ext(path) == app.Extension {
			files[path] = info.ModTime()
			time1, ok1 := files[path]
			time2, ok2 := app.Files[path]
			if ok1 && ok2 && !time1.Equal(time2) {
				fmt.Println(path, "Is Update")
				update = true
			}
		}
		return nil
	})
	app.Files = files
	if err != nil {
		fmt.Println("Error:", err)
		return false
	}
	return update
}

func KillPid(app *App) error {
	err := app.Process.Signal(os.Kill)
	if err != nil {
		fmt.Println("Kill process error:", err)
		return err
	}
	fmt.Println("Process are kill")
	return nil
}

func process(app *App) {
	sigs := make(chan os.Signal, 1)
	ticker := time.NewTicker(app.Interval)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)

	go func() {
		for {
			select {
			case <-sigs:
				fmt.Println("Thank see you next time...")
				return
			case <-ticker.C:
				fmt.Println("Running")
				update := Handler(app)
				if update {
					KillPid(app)
					return
				}
				fmt.Println("update:", update)
			}
		}
	}()
	<-sigs
}

//go:embed config.json
var configFile []byte

func ParseConfig(app *App) ([]Config, error) {
	var (
		config []Config
		err    error
	)
	if err = json.Unmarshal(configFile, &config); err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return config, nil
}

func main() {
	var (
		app    = New()
		err    error
		config []Config
	)
	config, err = ParseConfig(app)
	log.Println(config)
	if err != nil {
		fmt.Println("Please check your config")
		os.Exit(-1)
	}

	PrintBox(app)
	ok := strings.HasPrefix(app.Extension, ".")
	if !ok {
		fmt.Println("Please use correct extenstion:")
		os.Exit(-1)
	}
	if app.Cmd == "" {
		fmt.Println("Please add command with flag -cmd")
		os.Exit(-1)
	}
	fmt.Println("info", app.Cmd)
	err = app.Start().Errror()
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	process(app)
}
