package app

import (
	"fmt"
	"strings"

	"github.com/Arsylk/frida-clockwork-tui/internal/pkg/component"
	"github.com/Arsylk/frida-clockwork-tui/internal/pkg/event"
	"github.com/Arsylk/frida-clockwork-tui/internal/pkg/source"
	"github.com/Arsylk/frida-clockwork-tui/internal/pkg/style"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	ansi "github.com/leaanthony/go-ansi-parser"
	"github.com/mattn/go-runewidth"
	"github.com/muesli/termenv"
)

type StateFzf struct {
	helper

	data     *source.FormatLogData
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

func newStateFzf(application Application, data *source.FormatLogData) StateFzf {
	styles := style.DefaultStyles

	input := textinput.New()
	input.Prompt = "> "
	input.PromptStyle = styles.Prompt
	input.Focus()
	lipgloss.SetColorProfile(termenv.TrueColor)

	spinner := component.NewSpinnerModel()

	state := StateFzf{
		helper: helper{Application: application},

		data:     data,
		itemsLen: data.Len(),

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

	return state
}

func (m StateFzf) Init() tea.Cmd {
	runewidth.DefaultCondition.EastAsianWidth = false

	cmds := []tea.Cmd{
		textinput.Blink,
		tea.EnterAltScreen,
		m.spinner.TickCmd(),
		m.PerformSearchCmd(""),
	}

	for i, cmd := range cmds {
		m.Log(stringMaxCut(fmt.Sprintf("state fzf: Init(%d): %T\n", i, cmd), 100))
	}

	return tea.Batch(cmds...)
}

// update
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
	case event.FormatLogDataMsg:
		m.data = msg.Data
		m.itemsLen = m.data.Len()
	case event.LoadingMsg:
		m.spinner.IsEnabled = msg.IsLoading
	case event.SearchFinishedMsg:
		m.matches = *msg.Matches
		m.fixYPosition()
		m.fixCursor()
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
			// m.fixYPosition()
			// m.fixCursor()
		}
	}

	return m, tea.Batch(cmds...)
}

func (m *StateFzf) filter() {
	s := m.input.Value()
	m.PerformSearchCmd(s)()
}

func (m *StateFzf) toggle() {
	if len(m.matches) == 0 {
		return
	}

	match := m.matches[m.cursorPosition]
	indices := m.data.GetEntryIndices(match.Entry)

	if intContains(m.choices, match.Index) {
		m.choices = intFilter(m.choices, func(i int) bool { return !intContains(indices, i) })
	} else {
		m.choices = append(m.choices, indices...)
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
	m.helper.LoadFileCmd()
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
	log := fmt.Sprintf("cp: %d, %d", m.cursorPosition, m.Height())
	log, _ = ansi.Truncate(log, m.Width())
	return style.RenderFooter(log, m.Width())
}

func (m *StateFzf) footerHeight() int {
	return lipgloss.Height(m.footerView())
}

func (m *StateFzf) itemsHeight() int {
	h := min(m.Height()-m.inputHeight()-m.footerHeight(), m.Height())
	return h
}

func (m *StateFzf) itemsView() string {
	// m.Log(stringMaxCut(fmt.Sprintf("state fzf: itemsView(%d, %d, %d): cursor: %d, windowY %d \n", m.inputHeight(), m.itemsHeight(), m.footerHeight(), m.cursorPosition, m.windowYPosition), 100))
	itemsHeight := m.itemsHeight()
	if itemsHeight < 1 || m.data == nil {
		return ""
	}

	sliceEnd := min(itemsHeight+m.windowYPosition, len(m.matches))
	sliceStart := min(m.windowYPosition, sliceEnd)
	displayItems := m.matches[sliceStart:sliceEnd]
	rows := make([]string, itemsHeight)

	for i, item := range displayItems {
		cursorLine := m.cursorPosition == (i + m.windowYPosition)
		rows[i] = m.itemView(item.Index, m.data.Get(item.Index), cursorLine)
	}
	for i := len(displayItems); i < itemsHeight; i += 1 {
		rows[i] = m.styles.EllipsisStyle.Render("~")
	}

	return lipgloss.JoinVertical(lipgloss.Left, rows...)
}

func (m *StateFzf) itemView(index int, match *source.FormatLogItem, cursorLine bool) string {
	var sb strings.Builder

	pad := len(fmt.Sprintf("%d", m.data.GetEntryIndex()))
	var prefix string
	if match.LineIndex == 0 {
		prefix = fmt.Sprintf(" %*d │ ", pad, m.data.Get(index).EntryIndex)
	} else {
		prefix = fmt.Sprintf(" %*s │ ", pad, "")
	}

	sb.WriteString(m.styles.EllipsisStyle.Render(prefix))

	if cursorLine {
		sb.WriteString(m.styles.Cursor.Render(">"))
	} else {
		sb.WriteString(m.styles.Cursor.Render(" "))
	}

	// todo options enable/
	if intContains(m.choices, index) {
		if match.LineIndex == 0 {
			sb.WriteString(m.styles.SelectedPrefix.Render("*"))
		} else {
			sb.WriteString(m.styles.SelectedPrefix.Render("+"))
		}
	} else {
		sb.WriteString(m.styles.UnselectedPrefix.Render(" "))
	}

	maxWidth := m.Width() - 4 - (pad + 4)

	text := match.GetText()
	rightDots := lipgloss.Width(*text) >= maxWidth

	cut, _ := ansi.Truncate(*text, maxWidth)
	sb.WriteString(cut)

	if rightDots {
		sb.WriteString(m.styles.EllipsisStyle.Render(".."))
	}

	return sb.String()
}
