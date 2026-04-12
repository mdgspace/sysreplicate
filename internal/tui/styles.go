package tui

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	primaryBanner lipgloss.Style = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			PaddingTop(2).
			PaddingBottom(2).
			PaddingLeft(4).
			Width(80)

	secondaryBanner = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#626262")).
			PaddingTop(1).
			PaddingBottom(1).
			PaddingLeft(4).
			Width(80)

	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA"))

	cursorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7D56F4")).
			Bold(true)

	textStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#f1fab7"))

	checkedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#04B575"))

	itemStyle = lipgloss.NewStyle().
			PaddingLeft(1)

	footerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#626262"))

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF0000"))
	
	warningStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFA500"))
)
