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
	var err error
	prog.Process, err = os.FindProcess(prog.Pid)
	if err != nil {
		fmt.Println(app.error)
		return false
	}
	err = prog.Process.Signal(syscall.Signal(0))
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
		for _, cmds := range prog.Config.Cmd {
			parse := strings.Fields(cmds)
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
			stderr.Close()
			stdout.Close()
			// cmd.ProcessState
			if err != nil {
				prog.info = fmt.Sprint("failed to waiting command:", err)
				return
			}
		}
	}()
	prog.restart = true
	prog.process = false
	prog.check = false
	return nil
}

func killPidStop(prog *Program) error {
	prog.info = fmt.Sprintf("Kill pid %d", prog.Pid)
	if prog.Process != nil {
		err := prog.Process.Kill()
		if err != nil {
			fmt.Println(err)
			return err
		}
		fmt.Println("Wait kill...")
		prog.Process.Wait()
		prog.Pid = 0
	}
	return nil
}

func killPid(prog *Program) error {
	prog.info = fmt.Sprintf("Kill pid %d", prog.Pid)
	if prog.Process != nil {
		err := prog.Process.Signal(syscall.SIGHUP)
		if err != nil {
			fmt.Println(prog.info)
			return err
		}
		prog.Process.Wait()
		prog.Chan.kill <- true
		prog.Pid = 0
	}
	return nil
}

func (app *App) process(prog *Program) {
	prog.process = true
	prog.info = "Programm Running"
	// app.Chan.Call <- true
	for {
		if err := execCmd(prog); err != nil {
			prog.info = fmt.Sprint(err)
		}
	}
}
