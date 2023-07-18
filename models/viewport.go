package models

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/nsf/termbox-go"
)

type viewPort struct {
	app         *App
	Tabs        []string
	Stdout      []string
	Stdin       []string
	General     []string
	viewGeneral view
	viewStdout  view
	viewSterr   view
	activeTab   int
	done        chan bool
}

func (v *viewPort) listen(app *App) {
	err := termbox.Init()
	if err != nil {
		log.Fatal(err)
	}
	termbox.SetInputMode(termbox.InputEsc)

	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyArrowUp:
				if v.activeTab == 0 {
					v.viewGeneral.LineUp(1)
				} else if v.activeTab == 1 {
					v.viewStdout.LineUp(1)
				} else if v.activeTab == 2 {
					v.viewSterr.LineUp(1)
				}
			case termbox.KeyArrowDown:
				if v.activeTab == 0 {
					v.viewGeneral.LineDown(1)
				} else if v.activeTab == 1 {
					v.viewStdout.LineDown(1)
				} else if v.activeTab == 2 {
					v.viewSterr.LineDown(1)
				}
			case termbox.KeyCtrlC:
				v.done <- true
			case termbox.KeyTab:
				v.activeTab = min(v.activeTab+1, len(v.Tabs)-2)
				app.Chan.Call <- true
				break
			case termbox.KeyArrowRight:
				v.activeTab = min(v.activeTab+1, len(v.Tabs)-2)
				app.Chan.Call <- true
				break
			case termbox.KeyArrowLeft:
				v.activeTab = max(v.activeTab-1, 0)
				app.Chan.Call <- true
				break
			default:
				switch ev.Ch {
				case 'k':
					if v.activeTab == 0 {
						v.viewGeneral.LineUp(1)
					} else if v.activeTab == 1 {
						v.viewStdout.LineUp(1)
					} else if v.activeTab == 2 {
						v.viewSterr.LineUp(1)
					}
				case 'j':
					if v.activeTab == 0 {
						v.viewGeneral.LineDown(1)
					} else if v.activeTab == 1 {
						v.viewStdout.LineDown(1)
					} else if v.activeTab == 2 {
						v.viewSterr.LineDown(1)
					}
				case 'l', 'n':
					v.activeTab = min(v.activeTab+1, len(v.Tabs)-2)
					app.Chan.Call <- true
					break
				case 'h', 'p':
					v.activeTab = max(v.activeTab-1, 0)
					app.Chan.Call <- true
					break
				}
			}
		}
	}
}

func (v *viewPort) ViewLeft(prog *Program) string {
	doc := strings.Builder{}
	var renderedTabs []string
	tabs := []string{"GohotReaload", "handler", strings.Repeat(" ", (width - (18 + 13)))}
	var style lipgloss.Style

	for i, t := range tabs {
		if i == 0 {
			style = activeTabStyle.Copy()
		} else if i == 1 {
			style = inactiveTabStyle.Copy()
		} else {
			style = lipgloss.NewStyle().Border(underTabBorder).BorderForeground(highlightColor).Padding(0, 1)
		}
		border, _, _, _, _ := style.GetBorder()
		if i == 0 {
			border.BottomLeft = "│"
		} else if i == 1 {
			border.BottomRight = "┴"
		} else {
			border.BottomRight = "┐"
			border.TopRight = ""
			border.TopLeft = ""
			border.Top = ""
			border.Bottom = "─"
			border.BottomLeft = "─"
			border.Left = ""
			border.Right = ""
		}
		style = style.Border(border)
		renderedTabs = append(renderedTabs, style.Render(t))
	}

	row := lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)
	doc.WriteString(row)
	doc.WriteString("\n")
	title := infoStyle.Render(" From waxdred and G33KM44N38" + divider + url("https://github.com/waxdred/GoHotReload"))
	check := tab.Render("status: ", Style.Foreground(lipgloss.Color(orange)).Render("Checking"))
	if prog.restart {
		check = tab.Render("status: ", Style.Foreground(lipgloss.Color(red)).Render("Restart"))
	} else if prog.process {
		check = tab.Render("status: ", Style.Foreground(lipgloss.Color(green)).Render("Running"))
	} else if prog.check {
		check = tab.Render("status: ", Style.Foreground(lipgloss.Color(orange)).Render("Checking"))
	}
	rows := lipgloss.JoinHorizontal(
		lipgloss.Top,
		check,
		lipgloss.NewStyle().Border(lipgloss.Border{
			Bottom:      "─",
			BottomRight: "─",
			BottomLeft:  "─",
		}).BorderForeground(highlightColor).Padding(0, 1).Render(strings.Repeat(" ", lipgloss.Width(row)-lipgloss.Width(check)-8)),
	)
	gap := tabGap.Render(strings.Repeat(" ", max(0, width-lipgloss.Width(row)-2)))
	rows = lipgloss.JoinHorizontal(lipgloss.Bottom, rows, gap)
	exec := Style.Render("Executable:", Style.Foreground(lipgloss.Color(orange)).Render(prog.Config.Executable))
	path := Style.Render("Path      :", Style.Foreground(lipgloss.Color(orange)).Render(prog.Config.Path))
	Cmd := Style.Render("Cmd       :", Style.Foreground(lipgloss.Color(orange)).Render(prog.Config.Cmd[0]))
	Extension := Style.Render("Extension :", Style.Foreground(lipgloss.Color(orange)).Render(prog.Config.Extension))
	Pid := Style.Render("Pid       :", Style.Foreground(lipgloss.Color(orange)).Render(fmt.Sprint(prog.Pid)))
	Mem := Style.Render(
		"Men       :",
		Style.Foreground(lipgloss.Color(orange)).Render(fmt.Sprint(ExecPsMem(*prog), "%%")),
	)
	info := Style.Render("info      :", Style.Foreground(lipgloss.Color(green)).Render(prog.info))

	doc.WriteString(
		windowStyle.Width((lipgloss.Width(row) - windowStyle.GetHorizontalFrameSize())).
			Render(rows, "\n", title, "\n", exec, "\n", path, "\n", Cmd, "\n", Extension, "\n", Pid, "\n", Mem, "\n", info, "\n"),
	)

	return docStyle.Render(doc.String())
}

func (v *viewPort) ViewSdtout(input *Input) string {
	doc := strings.Builder{}

	var renderedTabs []string

	for i, t := range v.Tabs {
		var style lipgloss.Style
		isFirst, isLast, isActive := i == 0, i == len(v.Tabs)-1, i == v.activeTab
		if isActive {
			style = activeTabStyle.Copy()
		} else {
			style = inactiveTabStyle.Copy()
		}
		border, _, _, _, _ := style.GetBorder()
		if isFirst && isActive {
			border.BottomLeft = "╵"
		} else if isFirst && !isActive {
			border.BottomLeft = "└"
		} else if isLast && isActive {
			border.BottomRight = "│"
		} else if isLast && !isActive {
			border.BottomRight = "┘"
		}
		style = style.Border(border)
		renderedTabs = append(renderedTabs, style.Render(t))
	}

	row := lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)
	doc.WriteString(row)
	doc.WriteString("\n")
	return doc.String()
}

func (app *App) RunViewPort() {
	w := app.view
	go w.listen(app)
	go func() {
		check := false
		for {
			select {
			case <-w.done:
				termbox.Close()
				fmt.Println("\nClose Program...")
				if app.handlerProcess(&app.Program) {
					killPidStop(&app.Program)
				}
				fmt.Println("Thank see you next time...")
				os.Exit(0)
			case <-app.Chan.Call:
				app.Mu.Lock()
				clearScreen()
				if app.Program.Pid == 0 && !check {
					w.viewGeneral.ClearLine()
					w.viewStdout.ClearLine()
					w.viewSterr.ClearLine()
					app.Input.global = []string{}
					app.Input.stderr = []string{}
					app.Input.stdout = []string{}
					check = true
				} else if app.Program.Pid != 0 {
					check = false
				}
				fmt.Println(w.ViewLeft(&app.Program))
				fmt.Println(w.ViewSdtout(app.Input))
				if w.activeTab == 0 {
					w.viewGeneral.printVisbleLines()
				} else if w.activeTab == 1 {
					w.viewStdout.printVisbleLines()
				} else if w.activeTab == 2 {
					w.viewSterr.printVisbleLines()
				}
				app.Mu.Unlock()
			case err := <-app.Program.Chan.stderr:
				app.Input.stderr = append(app.Input.stderr, err+"\n")
				app.Input.global = append(app.Input.global, err+"\n")
				w.viewGeneral.Update(app.Input.global)
				w.viewSterr.Update(app.Input.stderr)
				if w.activeTab == 0 {
					w.viewGeneral.printVisbleLines()
				} else if w.activeTab == 2 {
					w.viewSterr.printVisbleLines()
				}
			case out := <-app.Program.Chan.stdout:
				app.Input.stdout = append(app.Input.stdout, out+"\n")
				app.Input.global = append(app.Input.global, out+"\n")
				w.viewGeneral.Update(app.Input.global)
				w.viewStdout.Update(app.Input.stdout)
				if w.activeTab == 0 {
					w.viewGeneral.printVisbleLines()
				} else if w.activeTab == 1 {
					w.viewStdout.printVisbleLines()
				}
			}
		}
	}()
}
