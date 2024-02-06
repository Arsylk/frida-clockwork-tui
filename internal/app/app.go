package app

import (
	"fmt"
	"os"

	"github.com/Arsylk/frida-clockwork-tui/internal/pkg/event"
	tea "github.com/charmbracelet/bubbletea"
)

type Application struct {
	Path           string
	LastWindowSize tea.WindowSizeMsg
	logFile        *os.File
	runner         *event.Runner
}

func (app Application) Width() int {
	w := app.LastWindowSize.Width
	return w
}

func (app Application) Height() int {
	h := app.LastWindowSize.Height
	return h
}

func (app Application) Log(format string, a ...any) {
	fmt.Fprintf(app.logFile, format, a...)
}

func (app Application) Start(program *tea.Program) {
	app.runner.Go(program)
}

func (app Application) Close() error {
	return app.logFile.Close()
}

func newApplication(path string, logPath string) Application {
	const (
		initialWidth  = 70
		initialHeight = 20
	)

	fp, _ := tea.LogToFile(logPath, "debug")

	runner := event.NewRunner()

	return Application{
		Path: path,
		LastWindowSize: tea.WindowSizeMsg{
			Width:  initialWidth,
			Height: initialHeight,
		},
		logFile: fp,
		runner:  runner,
	}
}

func NewModel(path string, logPath string) (tea.Model, Application) {
	app := newApplication(path, logPath)
	return newStateInitial(app), app
}
