package app

import (
	"cmp"
	"fmt"
	"slices"
	"strings"

	"github.com/Arsylk/frida-clockwork-tui/internal/pkg/component"
	"github.com/Arsylk/frida-clockwork-tui/internal/pkg/source"
	"github.com/Arsylk/frida-clockwork-tui/internal/pkg/style"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/koki-develop/go-fzf"
	ansi "github.com/leaanthony/go-ansi-parser"
	"github.com/mattn/go-runewidth"
	"github.com/muesli/termenv"
)

type StateFzf struct {
	helper

	items    *source.LogItems
	itemsLen int

	promptWidth    int
	cursorPosition int

	matches []source.EntryMatch
	choices []int

	windowWidth        int
	windowHeight       int
	windowYPosition    int
	mainViewWidth      int
	previewWindowWidth int

	styles style.StylesFzf

	input   textinput.Model
	spinner component.SpinnerModel
}

func newStateFzf(application Application, data *source.ParsedLogData) StateFzf {
	styles := style.DefaultStyles

	input := textinput.New()
	input.Prompt = "> "
	input.PromptStyle = styles.Prompt
	input.Focus()
	lipgloss.SetColorProfile(termenv.TrueColor)

	spinner := component.NewSpinnerModel()

	state := StateFzf{
		helper: helper{Application: application},

		matches: []source.EntryMatch{},

		cursorPosition: 0,
		choices:        []int{},

		promptWidth: lipgloss.Width(input.Prompt),

		// window
		windowWidth:     application.Width(),
		windowHeight:    application.Height(),
		windowYPosition: 0,

		// styles
		styles: styles,

		// components
		input:   input,
		spinner: spinner,
	}

	state.spinner.IsEnabled = true
	state.loadItems(data.GetItems())
	state.spinner.IsEnabled = false

	return state
}

func (m *StateFzf) loadItems(items *source.LogItems) {
	m.items = items
	if items != nil {
		m.itemsLen = len(*items)
	} else {
		m.itemsLen = 0
	}
	m.filter()
}

func (m StateFzf) Init() tea.Cmd {
	runewidth.DefaultCondition.EastAsianWidth = false

	cmds := []tea.Cmd{
		textinput.Blink,
		tea.EnterAltScreen,
		m.spinner.TickCmd(),
	}

	for i, cmd := range cmds {
		m.Log(stringMaxCut(fmt.Sprintf("state fzf: Init(%d): %T\n", i, cmd), 100))
	}

	return tea.Batch(cmds...)
}

// update
type forceReloadMsg struct{}

func (m StateFzf) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	m.Log(stringMaxCut(fmt.Sprintf("state fzf: Update: %T\n", msg), 100))
	m.helper = m.helper.Update(msg)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		case "ctrl+\\":
			return m.transactionMenu(m)
		case "tab":
			m.toggle()
			m.fixYPosition()
			m.fixCursor()
		case "up":
			m.cursorUp()
			m.fixYPosition()
			m.fixCursor()
		case "down":
			m.cursorDown()
			m.fixYPosition()
			m.fixCursor()
		}
	case tea.WindowSizeMsg:
		m.windowWidth = msg.Width
		m.windowHeight = msg.Height
		m.fixYPosition()
		m.fixCursor()
		m.fixWidth()
	case forceReloadMsg:
		m.forceReload()
		return m, nil
	}

	var cmd tea.Cmd
	var cmds []tea.Cmd
	previousInput := m.input.Value()

	{
		m.input, cmd = m.input.Update(msg)
		cmds = append(cmds, cmd)
		m.spinner, cmd = m.spinner.Update(msg)
		cmds = append(cmds, cmd)
	}

	if nextValue := m.input.Value(); previousInput != nextValue {
		if len(m.matches) > 0 || !strings.HasPrefix(nextValue, previousInput) {
			m.filter()
			m.fixYPosition()
			m.fixCursor()
		}
	}

	return m, tea.Batch(cmds...)
}

func (m *StateFzf) filter() {
	var matches []source.EntryMatch

	s := m.input.Value()
	if s == "" {
		matches = make([]source.EntryMatch, m.itemsLen)
		for i := 0; i < m.itemsLen; i += 1 {
			matches[i] = source.EntryMatch{
				// text:  (*m.items)[i].GetText(),
				Index: i,
			}
		}
	} else {
		results := fzf.Search(m.items, s, fzf.WithSearchCaseSensitive(false))
		matches = make([]source.EntryMatch, len(results))
		for i := 0; i < len(results); i += 1 {
			matches[i] = source.EntryMatch{
				// text:  (*m.items)[i].GetText(),
				Index: results[i].Index,
			}
		}
	}

	slices.SortFunc(matches, func(a, b source.EntryMatch) int {
		return cmp.Compare(a.Index, b.Index)
	})
	m.matches = matches
}

func (m *StateFzf) toggle() {
	if len(m.matches) == 0 {
		return
	}

	match := m.matches[m.cursorPosition]
	if intContains(m.choices, match.Index) {
		m.choices = intFilter(m.choices, func(i int) bool { return i != match.Index })
	} else {
		m.choices = append(m.choices, match.Index)
	}

	m.cursorDown()
}

func (m *StateFzf) cursorUp() {
	if m.cursorPosition > 0 {
		m.cursorPosition -= 1
	}
}

func (m *StateFzf) cursorDown() {
	if m.cursorPosition+1 < len(m.matches) {
		m.cursorPosition += 1
	}
}

func (m *StateFzf) fixCursor() {
	if m.cursorPosition < 0 {
		m.cursorPosition = 0
		return
	}

	if m.cursorPosition+1 > len(m.matches) {
		m.cursorPosition = max(len(m.matches)-1, 0)
		return
	}
}

func (m *StateFzf) fixYPosition() {
	inputHeight := m.inputHeight()
	footerHeight := m.footerHeight()
	displayHeight := m.windowHeight - inputHeight - footerHeight

	if displayHeight > len(m.matches) {
		m.windowYPosition = 0
		return
	}

	if m.cursorPosition < m.windowYPosition {
		m.windowYPosition = m.cursorPosition
		return
	}

	// m.Log(stringMaxCut(fmt.Sprintf("state fzf: fixYPosition(cursor: %d, windowY: %d, itemsHeight: %d): \n", m.cursorPosition, m.windowYPosition, displayHeight), 100))
	if m.cursorPosition >= m.windowYPosition+displayHeight {
		m.windowYPosition = min(m.windowYPosition+1, len(m.matches))
		return
	}
}

func (m *StateFzf) fixWidth() {
	m.mainViewWidth = m.windowWidth
	// todo handle preview window

	m.input.Width = m.mainViewWidth - m.promptWidth - 1
}

func (m *StateFzf) forceReload() {
	m.loadItems(m.items)
}

// view
func (m StateFzf) View() string {
	rows := make([]string, 3)

	windowStyle := lipgloss.NewStyle().
		Width(m.Width()).
		Height(m.Height()).
		AlignVertical(lipgloss.Top)
	rows[0] = m.inputView()
	rows[1] = m.itemsView()
	rows[2] = m.footerView()

	return windowStyle.Render(lipgloss.JoinVertical(lipgloss.Left, rows...))
}

func (m *StateFzf) inputView() string {
	rows := make([]string, 3)

	rows[0] = m.headerView()
	rows[1] = m.input.View()
	rows[2] = m.countView()

	return lipgloss.JoinVertical(lipgloss.Left, rows...)
}

func (m *StateFzf) headerView() string {
	return style.RenderLabelTop(m.Application.Path, m.Width())
}

func (m *StateFzf) countView() string {
	return style.RenderCountView(m.spinner.View(), len(m.matches), m.itemsLen, m.Width())
}

func (m *StateFzf) inputHeight() int {
	return lipgloss.Height(m.inputView())
}

func (m *StateFzf) footerView() string {
	return style.RenderFooter(fmt.Sprintf("cp: %d, %s", m.cursorPosition, m.choices), m.Width())
}

func (m *StateFzf) footerHeight() int {
	return lipgloss.Height(m.footerView())
}

func (m *StateFzf) itemsHeight() int {
	h := min(m.Height()-m.inputHeight()-m.footerHeight(), len(m.matches))
	return h
}

func (m *StateFzf) itemsView() string {
	// m.Log(stringMaxCut(fmt.Sprintf("state fzf: itemsView(%d, %d, %d): cursor: %d, windowY %d \n", m.inputHeight(), m.itemsHeight(), m.footerHeight(), m.cursorPosition, m.windowYPosition), 100))
	itemsHeight := m.itemsHeight()
	if itemsHeight < 1 {
		return ""
	}

	sliceEnd := min(itemsHeight+m.windowYPosition, len(m.matches)-1)
	sliceStart := min(m.windowYPosition, sliceEnd)
	displayItems := m.matches[sliceStart:sliceEnd]
	rows := make([]string, itemsHeight)

	for i, item := range displayItems {
		cursorLine := m.cursorPosition == (i + m.windowYPosition)
		rows[i] = m.itemView(item.Index, m.items.Get(item.Index).GetText(), cursorLine)
	}
	for i := len(displayItems); i < itemsHeight; i += 1 {
		rows[i] = m.styles.EllipsisStyle.Render("~")
	}

	return lipgloss.JoinVertical(lipgloss.Left, rows...)
}

func (m *StateFzf) itemView(index int, match string, cursorLine bool) string {
	var sb strings.Builder

	pad := len(fmt.Sprintf("%d", m.itemsLen))
	sb.WriteString(m.styles.EllipsisStyle.Render(
		fmt.Sprintf(" %*d │ ", pad, index),
	))

	if cursorLine {
		sb.WriteString(m.styles.Cursor.Render(">"))
	} else {
		sb.WriteString(m.styles.Cursor.Render(" "))
	}

	// todo options enable/disable
	if intContains(m.choices, index) {
		sb.WriteString(m.styles.SelectedPrefix.Render("*"))
	} else {
		sb.WriteString(m.styles.UnselectedPrefix.Render(" "))
	}

	maxWidth := (m.Width() - 2) - lipgloss.Width(sb.String())
	rightDots := false
	cut, _ := ansi.Truncate(match, maxWidth)
	sb.WriteString(cut)

	rightDots = (lipgloss.Width(sb.String()) > lipgloss.Width(cut))
	if rightDots {
		sb.WriteString(m.styles.EllipsisStyle.Render(".."))
	}

	return sb.String()
}
