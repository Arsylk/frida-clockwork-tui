package app

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type StateError struct {
	helper
	err error
}

func newStateError(application Application, err error) StateError {
	return StateError{
		helper: helper{
			Application: application,
		},
		err: err,
	}
}

func (s StateError) Init() tea.Cmd {
	return nil
}

func (s StateError) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	s.helper = s.helper.Update(msg)

	switch msg.(type) {
	case tea.KeyMsg:
		return s, tea.Quit
	}

	return s, nil
}

func (s StateError) View() string {
	return fmt.Sprintf("error: %s", s.err)
}
