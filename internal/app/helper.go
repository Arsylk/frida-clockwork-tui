package app

import (
	"github.com/Arsylk/frida-clockwork-tui/internal/pkg/source"
	tea "github.com/charmbracelet/bubbletea"
)

type helper struct {
	Application
}

// command & messages
func (h helper) PerformSearchCmd(query string) tea.Cmd {
	return func() tea.Msg {
		h.runner.SetSearchEvent(query)
		return nil
	}
}

func (h helper) LoadFileCmd() tea.Msg {
	h.runner.SetReadEvent(h.Path)
	return nil
}

// quick access
func (h helper) HandleQuit(msg tea.KeyMsg) (helper, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "esc":
		return h, tea.Quit
	}
	return h, nil
}

func (h helper) Update(msg tea.Msg) helper {
	if msg, ok := msg.(tea.WindowSizeMsg); ok {
		h.LastWindowSize = msg
	}

	return h
}

func (h helper) Width() int {
	return h.LastWindowSize.Width
}

func (h helper) Height() int {
	return h.LastWindowSize.Height
}

// change state screens
func (h helper) transactionLogView(entires *source.Entries) (tea.Model, tea.Cmd) {
	// return initializeModel(newStateLogView(h.Application, entires))
	return nil, nil
}

func (h helper) transactionFzf(data *source.FormatLogData) (tea.Model, tea.Cmd) {
	return initializeModel(newStateFzf(h.Application, data))
}

func (h helper) transactionError(err error) (tea.Model, tea.Cmd) {
	return initializeModel(newStateError(h.Application, err))
}

func (h helper) transactionMenu(prev tea.Model) (tea.Model, tea.Cmd) {
	return initializeModel(newStateMenu(h.Application, prev))
}

func (h helper) transactionBack(prev interface{ tea.Model }) (tea.Model, tea.Cmd) {
	return prev, nil
}
