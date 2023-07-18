package models

import (
	"os"

	"github.com/charmbracelet/lipgloss"
	"golang.org/x/crypto/ssh/terminal"
)

const (
	// In real life situations we'd adjust the document to fit the width we've
	// detected. In the case of this example we're hardcoding the width, and
	// later using the detected width only to truncate in order to avoid jaggy
	// wrapping.
	columnWidth = 30
)

var (
	fd                 = int(os.Stdout.Fd())
	Wwidth, Wheight, _ = terminal.GetSize(fd)
	WeightSet          = 23
	WeightChat         = Wwidth - WeightSet - 8
	width              = Wwidth - 10
	height             = Wheight / 3

	heightPrompt  = 1
	heightSetting = 6
	heightSession = Wheight - heightSetting - 7
	heightChat    = Wheight - heightPrompt - 7
)

func tabBorderWithBottom(left, middle, right string) lipgloss.Border {
	border := lipgloss.RoundedBorder()
	border.BottomLeft = left
	border.Bottom = middle
	border.BottomRight = right
	return border
}

var (
	blue       = "#1F6FEB"
	purple     = "#6D00E8"
	red        = "#F60000"
	blueSelect = "#8D908B"
	orange     = "#E7220D"
	green      = "#23A559"
	grey       = "#8D908B"
	yellow     = "#FF8C00"
	pink       = "#FF4D86"
	pinkDark   = "#B341E7"
	darkGrey   = "#343433"
	greyHelper = "#515150"
	subtle     = lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}
	highlight  = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	special    = lipgloss.AdaptiveColor{Light: "#43BF6D", Dark: "#73F59F"}

	inactiveTabBorder = tabBorderWithBottom("┴", "─", "┴")
	activeTabBorder   = tabBorderWithBottom("┘", " ", "└")
	underTabBorder    = tabBorderWithBottom("┘", " ", "└")
	docStyle          = lipgloss.NewStyle()
	highlightColor    = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	inactiveTabStyle  = lipgloss.NewStyle().Border(inactiveTabBorder, true).BorderForeground(highlightColor)
	activeTabStyle    = inactiveTabStyle.Copy().Border(activeTabBorder, true)
	windowStyle       = lipgloss.NewStyle().
				BorderForeground(highlightColor).
				Align(lipgloss.Left).
				Border(lipgloss.NormalBorder()).
				UnsetBorderTop()

	divider = lipgloss.NewStyle().
		SetString("•").
		Foreground(subtle).
		String()

	url = lipgloss.NewStyle().Foreground(special).Render

	// Tabs.

	tabBorder = lipgloss.Border{
		Top:         "─",
		Bottom:      "─",
		Left:        "│",
		Right:       "│",
		TopLeft:     "╭",
		TopRight:    "╮",
		BottomLeft:  "┴",
		BottomRight: "┴",
	}
	tabGapB = lipgloss.Border{
		Top:         "─",
		Bottom:      "─",
		Left:        "│",
		Right:       "│",
		TopLeft:     "╭",
		TopRight:    "╮",
		BottomLeft:  "└",
		BottomRight: "┘",
	}

	styleBorder = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			Width(WeightChat).MarginLeft(1).Height(heightChat).
			BorderForeground(lipgloss.Color(blue)).
			BorderTop(false)

	tab = lipgloss.NewStyle().
		Border(tabBorder, true).
		BorderForeground(highlight)

	gap = lipgloss.NewStyle().
		Border(tabGapB, true).
		BorderForeground(highlight)
	activeTab = tab.Copy().Border(activeTabBorder, true)

	tabGap = gap.Copy().
		BorderTop(false).
		BorderLeft(false).
		BorderRight(false)
	tabSep = gap.Copy().
		BorderTop(false).
		BorderBottom(false)

	// Title.

	titleStyle = lipgloss.NewStyle().
			MarginLeft(1).
			MarginRight(5).
			Italic(true).
			Foreground(lipgloss.Color("#FFF7DB")).
			SetString("Lip Gloss")

	descStyle = lipgloss.NewStyle().MarginTop(1)
	Style     = lipgloss.NewStyle()

	infoStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderTop(true).
			BorderForeground(subtle)

	// Dialog.

	dialogBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#874BFD")).
			BorderTop(true).
			BorderLeft(true).
			BorderRight(true).
			BorderBottom(true)

	buttonStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFF7DB")).
			Background(lipgloss.Color("#888B7E")).
			MarginTop(1)

	activeButtonStyle = buttonStyle.Copy().
				Foreground(lipgloss.Color("#FFF7DB")).
				Background(lipgloss.Color("#F25D94")).
				MarginRight(2).
				Underline(true)

	// List.

	list = lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), false, true, false, false).
		BorderForeground(subtle).
		MarginRight(2).
		Height(8).
		Width(columnWidth + 1)

	listHeader = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderBottom(true).
			BorderForeground(subtle).
			MarginRight(2).
			Render

	listItem = lipgloss.NewStyle().PaddingLeft(2).Render

	checkMark = lipgloss.NewStyle().SetString("✓").
			Foreground(special).
			PaddingRight(1).
			String()

	listDone = func(s string) string {
		return checkMark + lipgloss.NewStyle().
			Strikethrough(true).
			Foreground(lipgloss.AdaptiveColor{Light: "#969B86", Dark: "#696969"}).
			Render(s)
	}

	// Paragraphs/History.

	historyStyle = lipgloss.NewStyle().
			Align(lipgloss.Left).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(highlight).
			Margin(1, 3, 0, 0).
			Padding(1, 2).
			Height(19).
			Width(columnWidth)

	// Status Bar.

	statusNugget = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFDF5")).
			Padding(0, 1)

	statusBarStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#343433", Dark: "#C1C6B2"}).
			Background(lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#353533"})

	statusStyle = lipgloss.NewStyle().
			Inherit(statusBarStyle).
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#FF5F87")).
			Padding(0, 1).
			MarginRight(1)

	encodingStyle = statusNugget.Copy().
			Background(lipgloss.Color("#A550DF")).
			Align(lipgloss.Right)

	statusText = lipgloss.NewStyle().Inherit(statusBarStyle)

	fishCakeStyle = statusNugget.Copy().Background(lipgloss.Color("#6124DF"))

	// Page.
)
