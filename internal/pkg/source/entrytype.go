package source

import (
	"fmt"
	"strings"

	"github.com/Arsylk/frida-clockwork-tui/internal/pkg/style"
	"github.com/charmbracelet/lipgloss"
)

type JvmId string

type (
	JvmClass struct {
		LogEntry
		cn string
		id JvmId
	}
	JvmMethod struct {
		LogEntry
		mn string
		a  []string
		r  string
		id JvmId
	}
	JvmCall struct {
		LogEntry
		cn string
		mn string
		av []string
		rv string
		id JvmId
		st []string
	}
	JvmReturn struct {
		LogEntry
		id JvmId
		rv string
	}
)

var (
	_className  = lipgloss.NewStyle().Foreground(style.Cyan())
	_methodName = lipgloss.NewStyle().Foreground(style.Green())
	_bracket    = lipgloss.NewStyle().Foreground(style.Blue())
	_string     = lipgloss.NewStyle().Foreground(style.Yellow())
	_number     = lipgloss.NewStyle().Foreground(style.Magenta())
	_error      = lipgloss.NewStyle().Foreground(style.Red())
	_stacktrace = lipgloss.NewStyle().Foreground(style.BrightMagenta()).Faint(true)
	_keyword    = lipgloss.NewStyle().Faint(true)
	_unknown    = lipgloss.NewStyle().Foreground(style.White())
)

func renderClassName(cn string) string {
	parts := strings.Split(cn, ".")
	newparts := make([]string, len(parts))
	for i, part := range parts {
		newparts[i] = _className.Render(part)
	}
	return strings.Join(newparts, ".")
}

func renderMethodName(mn string) string {
	return _methodName.Render(mn)
}

func renderValue(v string, tp *string) string {
	if tp == nil {
		return v
	}

	switch *tp {
	case "int", "float", "double", "long", "byte", "boolean":
		return _number.Render(v)
	case "char":
		return _string.Render(fmt.Sprintf("'%s'", v))
	case "java.lang.String", "java.lang.CharSequence":
		return _string.Render(fmt.Sprintf("\"%s\"", v))
	}
	return fmt.Sprintf("%s%s%s %s", _bracket.Render("("), _error.Render(*tp), _bracket.Render(")"), _unknown.Render(v))
}

func renderStacktrace(st string) string {
	return _stacktrace.Render(fmt.Sprintf("  at %s", st))
}

func (e JvmClass) Render(data *ParsedLogData) *[]string {
	return &[]string{fmt.Sprintf("Hooking %s", renderClassName(e.cn))}
}

func (e JvmMethod) Render(data *ParsedLogData) *[]string {
	args := make([]string, len(e.a))
	for i, arg := range e.a {
		args[i] = renderClassName(arg)
	}
	return &[]string{fmt.Sprintf("  >%s%s%s%s: %s", renderMethodName(e.mn), _bracket.Render("("), strings.Join(args, ", "), _bracket.Render(")"), renderClassName(e.r))}
}

func (e JvmCall) Render(data *ParsedLogData) *[]string {
	var sb strings.Builder
	var sbArgs = func(sb *strings.Builder) {
		for i, arg := range e.av {
			str := renderValue(arg, data.GetArgType(e.id, i))
			(*sb).WriteString(str)
			if i < len(e.av)-1 {
				(*sb).WriteString(", ")
			}
		}
	}

	sb.WriteString(_keyword.Render("call"))
	sb.WriteString(" ")
	if e.mn == "$init" {
		sb.WriteString("new ")
		sb.WriteString(renderClassName(e.cn))
		sb.WriteString(_bracket.Render("("))
		sbArgs(&sb)
		sb.WriteString(_bracket.Render(")"))
	} else {
		sb.WriteString(renderClassName(e.cn))
		sb.WriteString("::")
		sb.WriteString(renderMethodName(e.mn))
		sb.WriteString(_bracket.Render("("))
		sbArgs(&sb)
		sb.WriteString(_bracket.Render(")"))
		sb.WriteString(": ")
		sb.WriteString(renderClassName(e.rv))
	}

	var result []string = make([]string, len(e.st))
	result[0] = sb.String()
	for i := 0; i < len(e.st)-1; i += 1 {
		result[i+1] = renderStacktrace(e.st[i])
	}

	return &result
}

func (e JvmReturn) Render(data *ParsedLogData) *[]string {
	value := renderValue(e.rv, data.GetReturnType(e.id))
	return &[]string{fmt.Sprintf("%s %s", _keyword.Render("return"), value)}
}
