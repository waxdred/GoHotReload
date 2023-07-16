package models

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	// "time"
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
	go func() {
		parse := strings.Fields(prog.Config.Cmd[0])
		if len(parse) == 0 {
			prog.info = "empty command string"
			return
		}
		args := parse[1:]
		cmd := exec.Command(parse[0], args...)
		cmd.Dir = prog.Config.Path

		stdout, err := cmd.StdoutPipe()
		if err != nil {
			prog.info = fmt.Sprint("failed to get stdout pipe:", err)
			return
		}
		stderr, err := cmd.StderrPipe()
		if err != nil {
			prog.info = fmt.Sprint("failed to get stderr pipe:", err)
			return
		}

		err = cmd.Start()
		if err != nil {
			prog.info = fmt.Sprint("failed to execute command:", err)
			return
		}

		go func() {
			scanner := bufio.NewScanner(stdout)
			for scanner.Scan() {
				prog.Chan.stdout <- scanner.Text()
			}
		}()

		go func() {
			scanner := bufio.NewScanner(stderr)
			for scanner.Scan() {
				prog.Chan.stderr <- scanner.Text()
			}
		}()

		err = cmd.Wait()
		if err != nil {
			prog.info = fmt.Sprint("failed to waiting command:", err)
			return
		}
	}()

	return nil
}

func killPid(prog *Program) error {
	prog.info = fmt.Sprintf("Kill pid %d", prog.Chan.pid)
	if prog.Process != nil {
		err := prog.Process.Signal(syscall.SIGHUP)
		if err != nil {
			prog.info = fmt.Sprintf("Kill process error: %v", err)
			return err
		}
	}
	return nil
}

func (app *App) notify(prog *Program) {
}

func (app *App) process(prog *Program) {
	prog.process = true
	prog.info = "Programm Running"
	app.Chan.Call <- true
	for {
		update := app.Handler(prog)
		if !app.handlerProcess(prog) {
			prog.info = "Pid stop by user"
			prog.check = true
			prog.process = false
			app.Chan.Call <- true
			return
		}

		if update {
			prog.process = false
			killPid(prog)
			app.Chan.Call <- true
			prog.restart = true
			prog.info = "Restart program"
			app.Chan.Call <- true
			if err := execCmd(prog); err != nil {
				prog.info = fmt.Sprint(err)
			}
			prog.restart = false
			break
		}
		// time.Sleep(time.Duration(prog.Config.Interval) * time.Second)
	}
}
