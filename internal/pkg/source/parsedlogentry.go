package source

import (
	"github.com/valyala/fastjson"
)

type LogEntry interface {
	GetRaw() *string
	GetTag() *string
}

type LogDataRenderer interface {
	Render(data *ParsedLogData) *[]string
}

type GenericEntry struct {
	raw string
	tag string
}

func (e GenericEntry) GetRaw() *string {
	return &e.raw
}

func (e GenericEntry) GetTag() *string {
	return &e.tag
}

func (data *ParsedLogData) JsonToParsedLogEntry(line string) (LogEntry, error) {
	var p fastjson.Parser
	value, err := p.Parse(line)
	if err != nil {
		return nil, err
	}

	tag := string(value.GetStringBytes("t"))
	generic := GenericEntry{
		raw: line,
		tag: tag,
	}

	switch tag {
	case "jvmclass":
		return JvmClass{
			LogEntry: &generic,
			cn:       string(value.GetStringBytes("cn")),
			id:       JvmId(value.GetStringBytes("id")),
		}, nil
	case "jvmmethod":
		avalues := value.GetArray("a")
		var a = make([]string, len(avalues))
		for i, avalue := range avalues {
			a[i] = string(avalue.GetStringBytes())
		}
		return JvmMethod{
			LogEntry: &generic,
			mn:       string(value.GetStringBytes("mn")),
			a:        a,
			r:        string(value.GetStringBytes("r")),
			id:       JvmId(value.GetStringBytes("id")),
		}, nil
	case "jvmcall":
		avalues := value.GetArray("av")
		var av = make([]string, len(avalues))
		for i, avalue := range avalues {
			av[i] = string(avalue.GetStringBytes())
		}
		stvalues := value.GetArray("st")
		var st = make([]string, len(stvalues))
		for i, stvalue := range stvalues {
			st[i] = string(stvalue.GetStringBytes())
		}
		return JvmCall{
			LogEntry: &generic,
			cn:       string(value.GetStringBytes("cn")),
			mn:       string(value.GetStringBytes("mn")),
			av:       av,
			rv:       string(value.GetStringBytes("rv")),
			id:       JvmId(value.GetStringBytes("id")),
			st:       st,
		}, nil
	case "jvmreturn":
		return JvmReturn{
			LogEntry: &generic,
			id:       JvmId(value.GetStringBytes("id")),
			rv:       string(value.GetStringBytes("rv")),
		}, nil
	}

	return nil, nil
}
