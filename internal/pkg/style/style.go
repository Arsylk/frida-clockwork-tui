package style

import (
	"fmt"
	"strings"

	lg "github.com/charmbracelet/lipgloss"
)

type StylesFzf struct {
	Prompt           lg.Style
	InputPlaceholder lg.Style
	InputText        lg.Style
	Cursor           lg.Style
	CursorLine       lg.Style
	SelectedPrefix   lg.Style
	UnselectedPrefix lg.Style
	Matches          lg.Style
	EllipsisStyle    lg.Style
}

var (
	DefaultColor  = "#00ADD8"
	DefaultStyles = StylesFzf{
		Prompt:           lg.NewStyle(),
		InputPlaceholder: lg.NewStyle().Faint(true),
		InputText:        lg.NewStyle(),
		Cursor:           lg.NewStyle().Foreground(lg.Color(DefaultColor)),
		CursorLine:       lg.NewStyle().Bold(true),
		Matches:          lg.NewStyle().Foreground(lg.Color(DefaultColor)),
		SelectedPrefix:   lg.NewStyle().Foreground(lg.Color(DefaultColor)),
		UnselectedPrefix: lg.NewStyle().Faint(true),
		EllipsisStyle:    lg.NewStyle().Faint(true),
	}
)

var (
	Style = lg.NewStyle()
)

func RenderLabelTop(name string, width int) string {
	base := lg.NewStyle().
		Width(width).
		Align(lg.Center).
		Foreground(mauve).
		Render(name)
	return strings.ReplaceAll(base, " ", "─")
}

func RenderPrompt() string {
	return lg.NewStyle().
		Foreground(mauve).
		Render("> ")
}

func RenderCountView(spinnerText string, match int, all int, width int) string {
	var info string = Style.Render(fmt.Sprintf("%d/%d", match, all))
	var spinner string
	if lg.Width(spinnerText) == 0 {
		spinner = Style.Render(" ")
	} else {
		spinner = Style.Render(spinnerText)
	}

	usedSpace := lg.Width(info) + lg.Width(spinner) + 1
	line := strings.Repeat("─", max(0, width-usedSpace))
	line = Style.Render(line)

	return fmt.Sprintf("%s %s%s", line, spinner, info)
}

func RenderFooter(text string, width int) string {
	render := lg.NewStyle().
		Render(text)
	line := lg.NewStyle().
		Render(strings.Repeat("─", max(0, width-lg.Width(render))))
	return lg.NewStyle().
		Render(fmt.Sprintf("%s%s", render, line))
}

const rosewater = lg.Color("#f5e0dc")
const flamingo = lg.Color("#f2cdcd")
const pink = lg.Color("#f5c2e7")
const mauve = lg.Color("#cba6f7")
const red = lg.Color("#f38ba8")
const maroon = lg.Color("#eba0ac")
const peach = lg.Color("#fab387")
const yellow = lg.Color("#f9e2af")
const green = lg.Color("#a6e3a1")
const teal = lg.Color("#94e2d5")
const sky = lg.Color("#89dceb")
const sapphire = lg.Color("#74c7ec")
const blue = lg.Color("#89b4fa")
const lavender = lg.Color("#b4befe")
const text = lg.Color("#cdd6f4")
const subtext1 = lg.Color("#bac2de")
const subtext0 = lg.Color("#a6adc8")
const overlay2 = lg.Color("#9399b2")
const overlay1 = lg.Color("#7f849c")
const overlay0 = lg.Color("#6c7086")
const surface2 = lg.Color("#585b70")
const surface1 = lg.Color("#45475a")
const surface0 = lg.Color("#313244")
const base = lg.Color("#1e1e2e")
const mantle = lg.Color("#181825")
const crust = lg.Color("#11111b")

func Black() lg.Color         { return lg.Color("#a6adc8") }
func Blue() lg.Color          { return lg.Color("#89b4fa") }
func BrightBlack() lg.Color   { return lg.Color("#585b70") }
func BrightBlue() lg.Color    { return lg.Color("#89b4fa") }
func BrightCyan() lg.Color    { return lg.Color("#89dceb") }
func BrightGreen() lg.Color   { return lg.Color("#a6e3a1") }
func BrightMagenta() lg.Color { return lg.Color("#f5c2e7") }
func BrightRed() lg.Color     { return lg.Color("#f38ba8") }
func BrightWhite() lg.Color   { return lg.Color("#45475a") }
func BrightYellow() lg.Color  { return lg.Color("#f9e2af") }
func Cyan() lg.Color          { return lg.Color("#89dceb") }
func Green() lg.Color         { return lg.Color("#a6e3a1") }
func Magenta() lg.Color       { return lg.Color("#f5c2e7") }
func Red() lg.Color           { return lg.Color("#f38ba8") }
func White() lg.Color         { return lg.Color("#bac2de") }
func Yellow() lg.Color        { return lg.Color("#f9e2af") }
func Border() lg.Color        { return lg.Color("#585b70") }
