package models

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
)

func HandlerSig(app *App) {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	<-sig
	fmt.Println("\nClose Program...")
	for _, prog := range app.Program {
		killPid(&prog)
	}
	fmt.Println("Thank see you next time...")
	os.Exit(1)
}

func (app *App) Handler(prog *Program) bool {
	files := make(map[string]time.Time)
	update := false

	err := filepath.Walk(prog.Config.Path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && filepath.Ext(path) == prog.Config.Extension {
			files[path] = info.ModTime()
			time1, ok1 := files[path]
			time2, ok2 := prog.Files[path]
			if ok1 && ok2 && !time1.Equal(time2) {
				prog.info = fmt.Sprint(path, "Is Update")
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
