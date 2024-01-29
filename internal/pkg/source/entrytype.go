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
		Entry
		cn string
		id JvmId
	}
	JvmMethod struct {
		Entry
		mn string
		a  []string
		r  string
		id JvmId
	}
	JvmCall struct {
		Entry
		cn string
		mn string
		av []string
		rv string
		id JvmId
		st []string
	}
	JvmReturn struct {
		Entry
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

func (e JvmClass) Render() string {
	return fmt.Sprintf("Hooking %s", renderClassName(e.cn))
}

func (e JvmMethod) Render() string {
	args := make([]string, len(e.a))
	for i, arg := range e.a {
		args[i] = renderClassName(arg)
	}
	return fmt.Sprintf("  >%s%s%s%s: %s", renderMethodName(e.mn), _bracket.Render("("), strings.Join(args, ", "), _bracket.Render(")"), renderClassName(e.r))
}

func (e JvmCall) Render() string {
	args := make([]string, len(e.av))
	for i, arg := range e.av {
		args[i] = renderValue(arg, e.GetSession().GetArgType(e.id, i))
	}
	if e.mn == "$init" {
		return fmt.Sprintf("call new %s%s%s%s", renderClassName(e.cn), _bracket.Render("("), strings.Join(args, ", "), _bracket.Render(")"))
	}
	return fmt.Sprintf("call %s::%s%s%s%s: %s", renderClassName(e.cn), renderMethodName(e.mn), _bracket.Render("("), strings.Join(args, ", "), _bracket.Render(")"), renderClassName(e.rv))
}

func (e JvmReturn) Render() string {
	value := renderValue(e.rv, e.GetSession().GetReturnType(e.id))
	return fmt.Sprintf("return %s", value)
}
