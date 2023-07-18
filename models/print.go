package models

import (
	"fmt"
	// "os"
	// "os/exec"
	// "runtime"
)

func clearScreen() {
	fmt.Print("\033[1000A")
	fmt.Print("\033[0J")
	// osName := runtime.GOOS
	// var arg string
	// if osName == "windows" {
	// 	arg = "cls"
	// } else {
	// 	arg = "clear"
	// }
	// cmd := exec.Command(arg)
	// cmd.Stdout = os.Stdout
	// cmd.Stderr = os.Stderr
	// cmd.Stdin = os.Stdin
	// cmd.Run()
}
