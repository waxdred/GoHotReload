package models

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func clearScreen() {
	cmd := exec.Command("clear") // Use "cls" instead of "clear" on Windows
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func printLine(title, value string) {
	if len(value) >= 30 {
		value = value[:30] + " ..."
	}
	line := fmt.Sprintf("│ %s: %s", title, value)
	fmt.Printf(line)
	fmt.Printf("%s│\n", strings.Repeat(" ", (54-len(line))))
}

func printMem(title, value string) {
	if len(value) >= 30 {
		value = value[:30] + " ..."
	}
	line := fmt.Sprintf("│ %s: %s", title, value)
	fmt.Printf(line)
	fmt.Printf("%s│\n", strings.Repeat(" ", (55-len(line))))
}

func printCheck(check, process, restart bool) {
	checkStr := fmt.Sprintf("%-5v", check)
	processStr := fmt.Sprintf("%-7v", process)
	restartStr := fmt.Sprintf("%-8v", restart)

	fmt.Printf("├───────────────┼────────────────┬──────────────────┤\n")
	fmt.Printf("│ Check  %v  │ Process %v│ Restart  %v│\n", checkStr, processStr, restartStr)
	fmt.Printf("├───────────────┴────────────────┴──────────────────┤\n")
}

func printHandler(prog Program, idx int) {
	fmt.Printf("│                                                   │\n")
	fmt.Printf("├───────────────┐                                   │\n")
	fmt.Printf("│ Handler: %d    │                                   │\n", idx)
	fmt.Printf("├───────────────┤                                   │\n")
	fmt.Printf("│ Status        │                                   │\n")
	printCheck(prog.check, prog.process, prog.restart)
	printLine("Executable", prog.Executable)
	printLine("Path      ", prog.Path)
	printLine("Cmd       ", prog.Cmd)
	printLine("Extension ", prog.Extension)
	printLine("Pid       ", fmt.Sprint(prog.Pid))
	printLine("TTY       ", prog.TTY)
	printMem("Mem       ", ExecPsMem(prog)+"%%")
	fmt.Printf("├───────────────────────────────────────────────────┤\n")
}

func (app *App) printBox() {
	clearScreen()
	top := "┌───────────────────────────────────────────────────┐\n"
	bottom := "└───────────────────────────────────────────────────┘\n"
	handler := fmt.Sprintf("│ handler .......... %d", len(app.Program))
	fmt.Printf(top)
	fmt.Printf("│                    GoHotRelaod                    │\n")
	fmt.Printf("│                                                   │\n")
	fmt.Printf(handler)
	fmt.Printf(strings.Repeat(" ", (54 - len(handler))))
	fmt.Printf("│\n")
	for i, prog := range app.Program {
		printHandler(prog, i)
	}
	fmt.Println(bottom)
}
