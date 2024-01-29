package app

import (
	"github.com/Arsylk/frida-clockwork-tui/internal/pkg/event"
	"github.com/Arsylk/frida-clockwork-tui/internal/pkg/source"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	useHighPerformanceRenderer = true
	vMarginTop                 = 3
	vMarginBottom              = 1
)

type logViewportModel struct {
	helper
	entries *source.Entries

	viewport viewport.Model
}

func newLogViewportModel(application Application, entries *source.Entries) logViewportModel {
	viewport := viewport.New(application.LastWindowSize.Width, application.LastWindowSize.Height-vMarginTop-vMarginBottom)
	viewport.HighPerformanceRendering = useHighPerformanceRenderer
	viewport.YPosition = vMarginTop + 1
	viewport.KeyMap.Up.Unbind()
	viewport.KeyMap.Down.Unbind()
	viewport.KeyMap.Up.SetKeys("up")
	viewport.KeyMap.Down.SetKeys("down")
	viewport.KeyMap.HalfPageUp.Unbind()
	viewport.KeyMap.HalfPageDown.Unbind()
	viewport.KeyMap.PageUp.Unbind()
	viewport.KeyMap.PageDown.Unbind()
	viewport.KeyMap.PageUp.SetKeys("pgup")
	viewport.KeyMap.PageDown.SetKeys("pgdown")

	return logViewportModel{
		helper:   helper{Application: application},
		entries:  entries,
		viewport: viewport,
	}
}

func (m logViewportModel) Init() tea.Cmd {
	return func() tea.Msg {
		// content := m.entries.MapToContent()
		all := len(*m.entries)
		found := all

		return event.PreparedContentMsg{
			Content: nil,
			All:     all,
			Found:   found,
		}
	}
}

func (m logViewportModel) Update(msg tea.Msg) (logViewportModel, tea.Cmd) {
	m.helper = m.helper.Update(msg)

	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m, cmds = runUpdate(m.resizeToWindowSize(msg.Width, msg.Height))(cmds)
	case event.PreparedContentMsg:
		m, cmds = runUpdate(m.updateContentAndCount(*msg.Content, msg.All, msg.Found))(cmds)
	}

	m.viewport, cmds = runUpdate(m.viewport.Update(msg))(cmds)

	return m, tea.Batch(cmds...)
}

func (m logViewportModel) View() string {
	return m.viewport.View()
}

func (m logViewportModel) resizeToWindowSize(width int, height int) (logViewportModel, tea.Cmd) {
	m.viewport.Width = width
	m.viewport.Height = height - vMarginTop - vMarginBottom
	m.viewport.YPosition = vMarginTop

	if useHighPerformanceRenderer {
		return m, viewport.Sync(m.viewport)
	}
	return m, nil
}

func (m logViewportModel) updateContentAndCount(content string, all int, found int) (logViewportModel, tea.Cmd) {
	m.viewport.SetContent(content)

	if useHighPerformanceRenderer {
		return m, viewport.Sync(m.viewport)
	}
	return m, nil
}
