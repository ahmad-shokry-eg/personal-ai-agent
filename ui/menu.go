package ui

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

type ActionType int

const (
	ActionBack ActionType = iota
	ActionDebug
	ActionPush
	ActionBuild
	ActionQuit
)

type MenuResult struct {
	Action   ActionType
	BuildCmd string
}

type model struct {
	choices  []string
	cursor   int
	submenu  bool
	subopts  []string
	subcur   int
	selected MenuResult
	quitting bool
}

func initialModel() model {
	return model{
		choices: []string{
			"Back",
			"Debug last command",
			"Push code",
			"Build the code",
		},
		subopts: []string{
			"go run .",
			"docker-compose up --build -d",
			"npm run build && npm run preview",
			"flutter build apk",
		},
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			if m.submenu && msg.String() == "esc" {
				m.submenu = false
				return m, nil
			}
			m.selected.Action = ActionQuit
			m.quitting = true
			return m, tea.Quit

		case "up", "k":
			if m.submenu {
				if m.subcur > 0 {
					m.subcur--
				}
			} else {
				if m.cursor > 0 {
					m.cursor--
				}
			}

		case "down", "j":
			if m.submenu {
				if m.subcur < len(m.subopts)-1 {
					m.subcur++
				}
			} else {
				if m.cursor < len(m.choices)-1 {
					m.cursor++
				}
			}

		case "enter", " ":
			if m.submenu {
				m.selected.Action = ActionBuild
				m.selected.BuildCmd = m.subopts[m.subcur]
				m.quitting = true
				return m, tea.Quit
			}

			switch m.cursor {
			case 0:
				m.selected.Action = ActionBack
				m.quitting = true
				return m, tea.Quit
			case 1:
				m.selected.Action = ActionDebug
				m.quitting = true
				return m, tea.Quit
			case 2:
				m.selected.Action = ActionPush
				m.quitting = true
				return m, tea.Quit
			case 3:
				m.submenu = true
			}
		}
	}
	return m, nil
}

func (m model) View() string {
	if m.quitting {
		return ""
	}

	s := "Select an action:\n\n"
	
	if m.submenu {
		s = "Select build command (ESC to go back):\n\n"
		for i, choice := range m.subopts {
			cursor := " "
			if m.subcur == i {
				cursor = ">"
			}
			s += fmt.Sprintf("%s %s\n", cursor, choice)
		}
	} else {
		for i, choice := range m.choices {
			cursor := " "
			if m.cursor == i {
				cursor = ">"
			}
			s += fmt.Sprintf("%s %s\n", cursor, choice)
		}
	}

	s += "\nPress j/k or up/down to move, enter to select, esc/q to quit.\n"
	return s
}

// ShowMenu runs the Bubble Tea UI and returns the user's selection
func ShowMenu() MenuResult {
	p := tea.NewProgram(initialModel())
	m, err := p.Run()
	if err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}

	if finalModel, ok := m.(model); ok {
		return finalModel.selected
	}

	return MenuResult{Action: ActionQuit}
}
