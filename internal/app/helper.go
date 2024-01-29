package app

import (
	"strings"

	"github.com/Arsylk/frida-clockwork-tui/internal/pkg/event"
	"github.com/Arsylk/frida-clockwork-tui/internal/pkg/source"
	tea "github.com/charmbracelet/bubbletea"
)

type helper struct {
	Application
}

// command & messages
func (h helper) LoadSession() tea.Msg {
	session, err := source.LoadFile(h.Path)
	if err != nil {
		return event.ErrorMsg{Err: err}
	}
	return event.OnLoadSessionMsg{Session: session}
}

func FilterLogAndPrepare(e *source.Entries, query string) func() tea.Msg {
	lowerQuery := strings.ToLower(query)
	return func() tea.Msg {
		filtered := e.Filter(func(entry source.IEntry, i int) bool {
			q := entry.Raw()
			return strings.Contains(*q, lowerQuery)
		})
		// content := filtered.MapToContent()
		all := len(*e)
		found := len(*filtered)

		return event.PreparedContentMsg{
			Content: nil,
			All:     all,
			Found:   found,
		}
	}
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
	return initializeModel(newStateLogView(h.Application, entires))
}

func (h helper) transactionFzf(sessions *source.Session) (tea.Model, tea.Cmd) {
	return initializeModel(newStateFzf(h.Application, sessions))
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
