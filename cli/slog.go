package main

import (
	"fmt"
	"os"

	"github.com/Arsylk/frida-clockwork-tui/internal/app"
	"github.com/Arsylk/frida-clockwork-tui/internal/pkg/source"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	var pathArg string
	if len(os.Args) > 1 {
		pathArg = os.Args[1]
	} else {
		// TODO Debug only
		pathArg = "/home/arsylk/session-j.txt"
	}
	if _, err := os.Stat(pathArg); err != nil {
		fatalf("error reading file \"%s\": %v\n", pathArg, err)
	}
	pathLog := "/home/arsylk/frida-clockwork-tui/debug.log"

	data, _ := source.LoadFromFile(pathArg)
	fmt.Fprintf(os.Stdout, "%T%s", data, pathLog)
	appModel, application := app.NewModel(pathArg, pathLog)

	defer application.Close()
	program := tea.NewProgram(
		appModel,
		tea.WithAltScreen(),
	)
	go application.Start(program)

	if _, err := program.Run(); err != nil {
		fatalf("program error: %v\n", err)
	}
}

func fatalf(message string, args ...any) {
	_, _ = fmt.Fprintf(os.Stderr, message, args...)
	os.Exit(1)
}
