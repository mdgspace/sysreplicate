package tui

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	choices  []string
	cursor   int
	selected int
}

func initialModel(display_list []string) model {
	return model{
		choices:  display_list,
		selected: 0,
	}
}

func PublicOptions(string_list []string) int {
	p := tea.NewProgram(initialModel(string_list))
	m, err := p.Run()
	if err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
	return m.(model).selected
}
