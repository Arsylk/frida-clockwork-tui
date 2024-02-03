package app

import (
	"fmt"

	"github.com/Arsylk/frida-clockwork-tui/internal/pkg/component"
	"github.com/Arsylk/frida-clockwork-tui/internal/pkg/event"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type StateMenu struct {
	helper
	prev tea.Model

	spinner component.SpinnerModel
}

func newStateMenu(application Application, prev tea.Model) StateMenu {
	spin := component.NewSpinnerModel()
	spin.Spinner = spinner.Dot
	spin.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	return StateMenu{
		helper: helper{
			Application: application,
		},
		prev:    prev,
		spinner: spin,
	}
}

func (s StateMenu) Init() tea.Cmd {
	return s.spinner.Tick
}

func (s StateMenu) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	s.Log(stringMaxCut(fmt.Sprintf("state menu: %T\n", msg), 100))
	s.helper = s.helper.Update(msg)

	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case spinner.TickMsg:
		s.spinner, cmd = s.spinner.Update(msg)
		cmds = append(cmds, cmd)
	case tea.KeyMsg:
		return s.transactionBack(s.prev)
	case event.ErrorMsg:
		return s.transactionError(msg.Err)
	}

	return s, tea.Batch(cmds...)
}

func (s StateMenu) View() string {
	return lipgloss.JoinHorizontal(lipgloss.Center, "Menu ...", s.spinner.View())
}
