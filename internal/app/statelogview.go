package app

import (
	"github.com/Arsylk/frida-clockwork-tui/internal/pkg/event"
	"github.com/Arsylk/frida-clockwork-tui/internal/pkg/source"
	"github.com/Arsylk/frida-clockwork-tui/internal/pkg/style"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type StateLogView struct {
	helper
	entries *source.Entries

	input    textinput.Model
	viewport logViewportModel
}

func newStateLogView(application Application, entries *source.Entries) StateLogView {
	input := textinput.New()
	input.Focus()

	viewport := newLogViewportModel(application, entries)

	return StateLogView{
		helper:  helper{Application: application},
		entries: entries,

		input:    input,
		viewport: viewport,
	}
}

func (s StateLogView) Init() tea.Cmd {
	return tea.Batch(s.viewport.Init())
}

func (s StateLogView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	s.helper = s.helper.Update(msg)

	var cmds []tea.Cmd

	prevQuery := s.input.Value()
	switch msg := msg.(type) {
	case event.ErrorMsg:
		return s.transactionError(msg.Err)
	case tea.KeyMsg:
		s.helper, cmds = runUpdate(s.HandleQuit(msg))(cmds)
		s.input, cmds = runUpdate(s.input.Update(msg))(cmds)
		newQuery := s.input.Value()
		if prevQuery != newQuery {
			cmds = append(cmds, FilterLogAndPrepare(s.entries, newQuery))
		}
	}

	s.viewport, cmds = runUpdate(s.viewport.Update(msg))(cmds)

	return s, tea.Batch(cmds...)
}

func (s StateLogView) View() string {
	return runRender(
		style.RenderLabelTop(s.Path, s.Width()),
		s.input.View(),
		style.RenderFooter("", s.Width()),
		s.viewport.View(),
		style.RenderFooter("", s.Width()),
	)
}
