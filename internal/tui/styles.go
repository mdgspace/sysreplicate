package tui

import (
	"github.com/charmbracelet/lipgloss"
)


func getCustomStyle(width int) lipgloss.Style {
	style := lipgloss.NewStyle().
    Bold(true).
    Foreground(lipgloss.Color("#FAFAFA")).
    Background(lipgloss.Color("#7D56F4")).
    PaddingTop(2).
	PaddingBottom(2).
    PaddingLeft(4).
    Width(width)
	return style
}

var (
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
)

// func main() {
// 	fmt.Println(style.Render("Hello, kitty"))
// }