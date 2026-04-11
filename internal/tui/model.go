package tui
import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	choices  []string         // items on list
	cursor   int              // which item our cursor is pointing at
	selected int // which items are selected
}

func initialModel(display_list []string) model {
	return model{
		choices: display_list,
		selected: 0, 
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

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
				m.selected = m.cursor+1
				return m, tea.Quit
			}
		}

		// Return the updated model to the Bubble Tea runtime for processing.
		// Note that we're not returning a command.
		return m, nil
}

func (m model) View() string {
	var s string

	s += getCustomStyle(80).Render("=== SysReplicate - Distro Hopping Tool ===")
	s += "\n\n"
	for i, choice := range m.choices {
		cursor := " "
		if m.cursor == i {
			cursor = cursorStyle.Render(">")
		}

		checked := " "
		if ok := m.selected; ok==i+1 {
			checked = checkedStyle.Render("x")
		}

		line := fmt.Sprintf("%s [%s] ", cursor, checked)
		s += itemStyle.Render(line)+ textStyle.Render(choice) + "\n"
	}

	s += "\n" + footerStyle.Render("Press q to quit.")

	return s
}

func PublicPrint(s []string) { // to use outside this package
	for _, str := range s {
		fmt.Println(textStyle.Render(str))
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

// currently how this works:
//When user boots, bubble tea starts, when he selects an option, bubble tea quits and returns the selected option to Run function in run.go
//and idk if I need to replace fmt everywhere, but for now, just made changes to run.go