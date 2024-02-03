package component

import (
	"time"

	"github.com/Arsylk/frida-clockwork-tui/internal/pkg/style"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	spinnerChars = spinner.Dot
	spinnerStyle = lipgloss.NewStyle().
			Foreground(style.Yellow())
)

type SpinnerModel struct {
	spinner.Model
	IsEnabled bool
}

type stateMsg struct {
	tea.Msg
	isEnabled bool
}

func NewSpinnerModel() SpinnerModel {
	model := spinner.New(spinner.WithSpinner(spinnerChars), spinner.WithStyle(spinnerStyle))
	return SpinnerModel{
		Model:     model,
		IsEnabled: true,
	}
}

func (m SpinnerModel) Init() tea.Cmd {
	return m.Model.Tick
}

func (m SpinnerModel) Update(msg tea.Msg) (SpinnerModel, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case stateMsg:
		m.IsEnabled = msg.isEnabled
		return m, cmd
	case spinner.TickMsg:
		m.Model, cmd = m.Model.Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m SpinnerModel) View() string {
	if !m.IsEnabled {
		return ""
	}
	return m.Model.View()
}

func (m SpinnerModel) TickCmd() tea.Cmd {
	return tea.Tick(m.Spinner.FPS, func(t time.Time) tea.Msg {
		return m.Model.Tick()
	})
}
