package tui

import (
    tea "github.com/charmbracelet/bubbletea"
)

// The Init method is called when a new bubbletea program is started. It returns the startup command, which in this case is nil since we don't need to perform any initial actions.
func (m model) Init() tea.Cmd {
    return nil
}

// The Update method is called whenever a message is received. It is responsible for updating our model (state)
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "ctrl+c", "q":
            m.selected = 6
            return m, tea.Quit
        case "up", "k":
            if m.cursor > 0 {
                m.cursor--
            }
        case "down", "j":
            if m.cursor < len(m.choices)-1 {
                m.cursor++
            }
        case "enter", " ":
            m.selected = m.cursor + 1
            return m, tea.Quit
        }
    }
    return m, nil
}

// The View method is responsible for rendering the UI based on the current state of the model.
func (m model) View() string {
	// The View method returns a single string that represents the entire UI. We build this string by concatenating different parts of the UI together.
    var s string

    s += RenderPrimaryBanner("=== SysReplicate - Distro Hopping Tool ===")
    s += "\n\n"

    for i, choice := range m.choices {
        s += RenderListItem(i, m.cursor, m.selected, choice)
    }

    s += "\n" + RenderFooter("Press q to quit.")
    return s
}