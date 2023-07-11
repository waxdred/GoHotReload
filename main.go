package main

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func PrintBox(app *App) {
	// top := "┌───────────────────────────────────────────────────┐\n"
	// bottom := "└───────────────────────────────────────────────────┘\n"
	// cmd := fmt.Sprintf("│  cmd .......... %s", app.Cmd)
	// extension := fmt.Sprintf("│  extension .... %s", app.Extension)
	// path := fmt.Sprintf("│  path ......... %s", app.Path)
	// pid := fmt.Sprint("│  pid ........... ", app.Pid)
	// fmt.Printf(top)
	// fmt.Printf("│                    GoHotRelaod                    │\n")
	// fmt.Printf("│                                                   │\n")
	// fmt.Printf(cmd)
	// fmt.Printf(strings.Repeat(" ", (54 - len(cmd))))
	// fmt.Printf("│\n")
	// fmt.Printf(path)
	// fmt.Printf(strings.Repeat(" ", (54 - len(path))))
	// fmt.Printf("│\n")
	// fmt.Printf(pid)
	// fmt.Printf(strings.Repeat(" ", (54 - len(extension))))
	// fmt.Printf("│\n")
	// fmt.Printf(extension)
	// fmt.Printf(strings.Repeat(" ", (54 - len(extension))))
	// fmt.Printf("│\n")

	// fmt.Println(bottom)
}

func Handler(prog *Program) bool {
	files := make(map[string]time.Time)
	update := false

	err := filepath.Walk(prog.Path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// fmt.Println("foo", prog.Extension, info.IsDir(), filepath.Ext(path), prog.Path)
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
	// PrintBox(app)
	err := app.Start().Errror()
	if err != nil {
		app.Mu.Lock()
		fmt.Println(err)
		app.Mu.Unlock()
		os.Exit(-1)
	}
}
