package models

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"time"
)

func (app *App) handlerProcess(prog *Program) bool {
	prog.Process, app.error = os.FindProcess(prog.Pid)
	if app.error != nil {
		return false
	}
	err := prog.Process.Signal(syscall.Signal(0))
	if err == nil {
		return true
	}

	if err.Error() == "os: process already finished" {
		return false
	}
	return false
}

func ExecPsMem(prog Program) string {
	cmd := exec.Command("ps", "-p", fmt.Sprint(prog.Pid), "-o", "%mem")
	output, err := cmd.Output()
	if err != nil {
		return "0"
	}
	lines := strings.Split(string(output), "\n")
	if len(lines) >= 1 {
		line := strings.Replace(lines[1], " ", "", 1)
		return line
	}
	return "0"
}

func (app *App) execPs(prog *Program) error {
	cmd := exec.Command("ps")
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Error cmd:", output)
		return err
	}
	lines := strings.Split(string(output), "\n")
	for _, line := range lines[1:] {
		fields := strings.Fields(line)
		if len(line) >= 4 {
			pids := fields[0]
			tty := fields[1]
			pid, err := strconv.Atoi(pids)
			if err != nil {
				return err
			}
			if pid == prog.Pid {
				prog.TTY = tty
			}
		}
	}
	return nil
}

func execCmd(prog *Program) error {
	// open tty fd
	go func() {
		// fmt.Println("Open tty")
		tty, err := os.OpenFile(fmt.Sprint("/dev/", prog.TTY), os.O_RDWR, 0)
		defer tty.Close()
		if err != nil {
			fmt.Println(err)
			return
		}
		parse := strings.Fields(prog.Cmd)
		if len(parse) == 0 {
			fmt.Println("empty command string")
			return
		}
		args := parse[1:]
		cmd := exec.Command(parse[0], args...)
		cmd.Dir = prog.Path

		// change stdin and out in the tty
		cmd.Stdin = tty
		cmd.Stdout = tty
		cmd.Stderr = tty
		cmd.SysProcAttr = &syscall.SysProcAttr{
			Setsid: true,
		}

		err = cmd.Start()
		if err != nil {
			fmt.Println("failed to execute command:", err)
			return
		}

		err = cmd.Wait()
		if err != nil {
			fmt.Println("Erreur lors de l'attente de la commande:", err)
			return
		}
	}()

	return nil
}

func killPid(prog *Program) error {
	err := prog.Process.Signal(os.Kill)
	if err != nil {
		fmt.Println("Kill process error:", err)
		return err
	}
	// fmt.Println("Process are kill")
	return nil
}

func (app *App) process(prog *Program) {
	// fmt.Println("Process", prog.Executable, "Running")
	prog.process = true
	for {
		update := Handler(prog)
		if !app.handlerProcess(prog) {
			fmt.Println("Pid stop by user")
			return
		}
		if update {
			prog.process = false
			killPid(prog)
			prog.restart = true
			time.Sleep(5 * time.Second)
			if err := execCmd(prog); err != nil {
				fmt.Println(err)
			}
			prog.restart = false
			break
		}
		time.Sleep(time.Duration(prog.Interval) * time.Second)
	}
}
