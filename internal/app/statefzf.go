package app

import (
	"fmt"
	"strings"
	"time"

	"github.com/Arsylk/frida-clockwork-tui/internal/pkg/fzfimpl"
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
	session  *source.Session
	items    *fzfimpl.Items
	itemsLen int

	promptWidth    int
	cursorPosition int

	matches fzf.Matches
	choices []int

	windowWidth        int
	windowHeight       int
	windowYPosition    int
	mainViewWidth      int
	previewWindowWidth int

	styles style.StylesFzf

	input textinput.Model
}

func newStateFzf(application Application, session *source.Session) StateFzf {
	entries := session.Entries()
	items, err := fzfimpl.NewColdItems(entries, func(i int) string {
		return source.RawEntry(entries.Get(i))
	})
	if err != nil {
		panic(err)
	}

	styles := style.DefaultStyles

	input := textinput.New()
	input.Prompt = "> "
	input.PromptStyle = styles.Prompt
	input.Focus()
	lipgloss.SetColorProfile(termenv.TrueColor)

	state := StateFzf{
		helper:  helper{Application: application},
		session: session,

		cursorPosition: 0,
		matches:        fzf.Matches{},
		choices:        []int{},

		promptWidth: lipgloss.Width(input.Prompt),

		// window
		windowWidth:     application.Width(),
		windowHeight:    application.Height(),
		windowYPosition: 0,

		// styles
		styles: styles,

		// components
		input: input,
	}

	state.loadItems(items)

	return state
}

func (m *StateFzf) loadItems(items *fzfimpl.Items) {
	m.items = items
	m.itemsLen = items.Len()
	m.filter()
}

func (m StateFzf) Init() tea.Cmd {
	runewidth.DefaultCondition.EastAsianWidth = false

	cmds := []tea.Cmd{
		textinput.Blink,
		tea.EnterAltScreen,
	}

	for i, cmd := range cmds {
		m.Log(stringMaxCut(fmt.Sprintf("state fzf: Init(%d): %T\n", i, cmd), 100))
	}

	return tea.Batch(cmds...)
}

// update
type watchReloadMsg struct{}
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
	case watchReloadMsg:
		return m, m.watchReload()
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
	s := m.input.Value()
	if s == "" {
		var matches fzf.Matches
		for i := 0; i < m.items.Len(); i += 1 {
			matches = append(matches, fzf.Match{
				Str:   m.items.ItemString(i),
				Index: i,
			})
		}
		m.matches = matches
	}
	m.matches = fzf.Search(m.items, s, fzf.WithSearchCaseSensitive(false))
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

func (m *StateFzf) watchReload() tea.Cmd {
	return tea.Tick(30*time.Microsecond, func(_ time.Time) tea.Msg {
		if m.itemsLen != m.items.Len() {
			m.loadItems(m.items)
		}
		return watchReloadMsg{}
	})
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
	return style.RenderCountView(len(m.matches), m.itemsLen, m.Width())
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

	matches := m.matches[m.windowYPosition : itemsHeight+m.windowYPosition]
	rows := make([]string, len(matches))

	for i, match := range matches {
		cursorLine := m.cursorPosition == (i + m.windowYPosition)
		entry := m.session.Entries().Get(match.Index)
		text := entry.Render()
		rows[i] = m.itemView(match.Index, text, cursorLine)
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

	maxWidth := m.Width() - lipgloss.Width(sb.String())
	rightDots := false
	cut, _ := ansi.Truncate(match, maxWidth)
	sb.WriteString(cut)

	if rightDots {
		sb.WriteString(m.styles.EllipsisStyle.Render(".."))
	}

	return sb.String()
}
