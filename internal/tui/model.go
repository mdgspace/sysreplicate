package tui

import (
	"fmt"

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

func PublicOptions(string_list []string) (int, error) {
	p := tea.NewProgram(initialModel(string_list))
	m, err := p.Run()
	if err != nil {
		return 0, fmt.Errorf("failed to run TUI: %w", err)
	}
	return m.(model).selected, nil
}
