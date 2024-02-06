module github.com/Arsylk/frida-clockwork-tui

go 1.21.6

require (
	// github.com/Arsylk/frida-clockwork-tui/styles
	github.com/charmbracelet/bubbletea v0.25.0
	github.com/muesli/termenv v0.15.2
)

require (
	github.com/charmbracelet/bubbles v0.17.1
	github.com/charmbracelet/lipgloss v0.9.1
	github.com/junegunn/fzf v0.0.0-20240201091300-3c0a6304756e
	github.com/koki-develop/go-fzf v0.15.0
	github.com/leaanthony/go-ansi-parser v1.6.1
	github.com/mattn/go-runewidth v0.0.15
	github.com/valyala/fastjson v1.6.4
)

require (
	github.com/atotto/clipboard v0.1.4 // indirect
	github.com/aymanbagabas/go-osc52/v2 v2.0.1 // indirect
	github.com/containerd/console v1.0.4-0.20230313162750-1ae8d489ac81 // indirect
	github.com/gdamore/encoding v1.0.0 // indirect
	github.com/gdamore/tcell/v2 v2.7.0 // indirect
	github.com/lucasb-eyer/go-colorful v1.2.0 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mattn/go-localereader v0.0.1 // indirect
	github.com/mattn/go-shellwords v1.0.12 // indirect
	github.com/muesli/ansi v0.0.0-20230316100256-276c6243b2f6 // indirect
	github.com/muesli/cancelreader v0.2.2 // indirect
	github.com/muesli/reflow v0.3.0 // indirect
	github.com/rivo/uniseg v0.4.6 // indirect
	github.com/saracen/walker v0.1.3 // indirect
	golang.org/x/sync v0.6.0 // indirect
	golang.org/x/sys v0.16.0 // indirect
	golang.org/x/term v0.16.0 // indirect
	golang.org/x/text v0.14.0 // indirect
)

replace github.com/junegunn/fzf => ./go/pkg/mod/github.com/junegunn/fzf@custom
