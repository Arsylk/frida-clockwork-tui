package app

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func initializeModel[T tea.Model](m T) (T, tea.Cmd) {
	return m, m.Init()
}

func runUpdate[T any](m T, cmd tea.Cmd) func(cmds []tea.Cmd) (T, []tea.Cmd) {
	return func(cmds []tea.Cmd) (T, []tea.Cmd) {
		if cmd != nil {
			cmds = append(cmds, cmd)
		}

		return m, cmds
	}
}

func runRender(parts ...string) string {
	all := ""

	length := len(parts)
	for i := 0; i < length; i += 1 {
		all = fmt.Sprintf("%s%s", all, parts[i])
		if i < length-1 {
			all = fmt.Sprintf("%s\n", all)
		}
	}

	return all
}
func intContains(is []int, i int) bool {
	for _, l := range is {
		if l == i {
			return true
		}
	}
	return false
}

func intFilter(is []int, f func(i int) bool) []int {
	var rtn []int
	for _, i := range is {
		if f(i) {
			rtn = append(rtn, i)
		}
	}
	return rtn
}

func stringMaxCut(text string, max int) string {
	runes := []rune(text)
	for lipgloss.Width(text) > max {
		text = string(runes[:len(runes)-1])
	}
	return string(runes)
}
