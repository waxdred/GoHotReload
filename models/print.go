package models

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/charmbracelet/lipgloss"
	// "github.com/charmbracelet/lipgloss"
)

func clearScreen() {
	cmd := exec.Command("clear") // Use "cls" instead of "clear" on Windows
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func (app *App) printBox() {
	clearScreen()
	doc := strings.Builder{}
	handler := fmt.Sprintf("handler: %d", len(app.Program))
	row := lipgloss.JoinHorizontal(
		lipgloss.Top,
		activeTab.Render("GoHotReaload"),
		tab.Render(handler),
	)
	gap := tabGap.Render(strings.Repeat(" ", max(0, width-lipgloss.Width(row)-2)))
	row = lipgloss.JoinHorizontal(lipgloss.Bottom, row, gap)
	doc.WriteString(row + "\n")
	// Title
	{
		var title strings.Builder
		desc := lipgloss.JoinVertical(lipgloss.Left,
			infoStyle.Render("From waxdred and G33KM44N38"+divider+url("https://github.com/waxdred/GoHotReload")),
		)

		row := lipgloss.JoinHorizontal(lipgloss.Top, title.String(), desc)
		doc.WriteString(row + "\n")
	}

	// Dialog
	{
		for i, prog := range app.Program {
			check := tab.Render("status: ", Style.Foreground(lipgloss.Color(orange)).Render("Checking"))
			if prog.restart {
				check = tab.Render("status: ", Style.Foreground(lipgloss.Color(red)).Render("Restart"))
			} else if prog.process {
				check = tab.Render("status: ", Style.Foreground(lipgloss.Color(green)).Render("Running"))
			} else if prog.check {
				check = tab.Render("status: ", Style.Foreground(lipgloss.Color(orange)).Render("Checking"))
			}
			row := lipgloss.JoinHorizontal(
				lipgloss.Top,
				activeTab.Render("handler: ", fmt.Sprint(i)),
				check,
			)
			gap := tabGap.Render(strings.Repeat(" ", max(0, width-lipgloss.Width(row)-2)))
			row = lipgloss.JoinHorizontal(lipgloss.Bottom, row, gap)
			exec := Style.Render("Executable:", Style.Foreground(lipgloss.Color(orange)).Render(prog.Config.Executable))
			path := Style.Render("Path      :", Style.Foreground(lipgloss.Color(orange)).Render(prog.Config.Path))
			Cmd := Style.Render("Cmd       :", Style.Foreground(lipgloss.Color(orange)).Render(prog.Config.Cmd))
			Extension := Style.Render("Extension :", Style.Foreground(lipgloss.Color(orange)).Render(prog.Config.Extension))
			Pid := Style.Render("Pid       :", Style.Foreground(lipgloss.Color(orange)).Render(fmt.Sprint(prog.Pid)))
			TTY := Style.Render("TTY       :", Style.Foreground(lipgloss.Color(orange)).Render(fmt.Sprint(prog.TTY)))
			Mem := Style.Render(
				"Men       :",
				Style.Foreground(lipgloss.Color(orange)).Render(fmt.Sprint(ExecPsMem(prog), "%%")),
			)
			info := Style.Render("info      :", Style.Foreground(lipgloss.Color(green)).Render(prog.info))
			doc.WriteString(
				row + "\n" + exec + "\n" + path + "\n" + Cmd + "\n" + Extension + "\n" + Pid + "\n" + TTY + "\n" + Mem + "\n" + info + "\n\n",
			)
		}
	}

	fmt.Printf(doc.String())
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
