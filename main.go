package main

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func PrintBox(app *App) {
	top := "┌───────────────────────────────────────────────────┐\n"
	bottom := "└───────────────────────────────────────────────────┘\n"
	cmd := fmt.Sprintf("│  process .......... %d", len(app.Program))
	fmt.Printf(top)
	fmt.Printf("│                    GoHotRelaod                    │\n")
	fmt.Printf("│                                                   │\n")
	fmt.Printf(cmd)
	fmt.Printf(strings.Repeat(" ", (54 - len(cmd))))
	fmt.Printf("│\n")
	fmt.Println(bottom)
}

func Handler(prog *Program) bool {
	files := make(map[string]time.Time)
	update := false

	err := filepath.Walk(prog.Path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && filepath.Ext(path) == prog.Extension {
			files[path] = info.ModTime()
			time1, ok1 := files[path]
			time2, ok2 := prog.Files[path]
			if ok1 && ok2 && !time1.Equal(time2) {
				fmt.Println(path, "Is Update")
				update = true
			}
		}
		return nil
	})
	prog.Files = files
	if err != nil {
		fmt.Println("Error:", err)
		return false
	}
	return update
}

func KillPid(prog *Program) error {
	err := prog.Process.Signal(os.Kill)
	if err != nil {
		fmt.Println("Kill process error:", err)
		return err
	}
	fmt.Println("Process are kill")
	return nil
}

func main() {
	app := New()
	if app.Errror() != nil {
		fmt.Println(app.Errror())

		os.Exit(-1)
	}
	app.CheckingParse()
	if app.Errror() != nil {
		app.Mu.Lock()
		fmt.Println(app.Errror())
		app.Mu.Unlock()
		os.Exit(-1)
	}
	fmt.Println(app.Program)
	PrintBox(app)
	err := app.Start().Errror()
	if err != nil {
		app.Mu.Lock()
		fmt.Println(err)
		app.Mu.Unlock()
		os.Exit(-1)
	}
}
