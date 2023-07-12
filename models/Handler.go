package models

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
)

func HandlerSig() {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	<-sig
	fmt.Println("Thank see you next time...")
	// TODO store fd open for close here correctly
	os.Exit(1)
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
