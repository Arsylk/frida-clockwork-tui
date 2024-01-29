package app

import (
	"fmt"

	"github.com/Arsylk/frida-clockwork-tui/internal/pkg/event"
	tea "github.com/charmbracelet/bubbletea"
)

type StateMenu struct {
	helper
	prev tea.Model
}

func newStateMenu(application Application, prev tea.Model) StateMenu {
	return StateMenu{
		helper: helper{
			Application: application,
		},
		prev: prev,
	}
}

func (s StateMenu) Init() tea.Cmd {
	return s.helper.LoadSession
}

func (s StateMenu) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	s.Log(stringMaxCut(fmt.Sprintf("state initial: %T\n", msg), 100))
	s.helper = s.helper.Update(msg)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		return s.transactionBack(s.prev)
	case event.ErrorMsg:
		return s.transactionError(msg.Err)
	}

	return s, nil
}

func (s StateMenu) View() string {
	return "Menu ..."
}
