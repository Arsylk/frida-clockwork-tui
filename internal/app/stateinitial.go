package app

import (
	"fmt"

	"github.com/Arsylk/frida-clockwork-tui/internal/pkg/event"
	tea "github.com/charmbracelet/bubbletea"
)

type StateInitial struct {
	helper
}

func newStateInitial(application Application) StateInitial {
	return StateInitial{
		helper: helper{
			Application: application,
		},
	}
}

func (s StateInitial) Init() tea.Cmd {
	return s.helper.LoadFileCmd
}

func (s StateInitial) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	s.Log(stringMaxCut(fmt.Sprintf("state initial: %T\n", msg), 100))
	s.helper = s.helper.Update(msg)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		return nil, tea.Quit
	case event.ErrorMsg:
		return s.transactionError(msg.Err)
	case event.OnLoadLogData:
		return s.transactionFzf(msg.Data)
	case event.LoadedEntriesMsg:
		return s.transactionLogView(msg.Entries)
	}

	return s, nil
}

func (s StateInitial) View() string {
	return "Loading ..."
}
