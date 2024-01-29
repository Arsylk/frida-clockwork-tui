package source

import (
	"bufio"
	"fmt"
	"os"
	"reflect"

	"github.com/valyala/fastjson"
)

func LoadFile(path string) (_ *Session, err error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("os: %w", err)
	}
	defer file.Close()

	session := Session{
		classes: make(map[JvmId]*JvmClass),
		methods: make(map[JvmId]*JvmMethod),
		entries: make(Entries, 0),
	}

	var sc fastjson.Scanner
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		sc.Init(line)
		if sc.Next() {
			value := sc.Value()
			tag := string(value.GetStringBytes("t"))

			switch tag {
			case "jvmclass":
				entry := JvmClass{
					Entry: Entry{
						raw:     line,
						tag:     tag,
						session: &session,
					},
					cn: string(value.GetStringBytes("cn")),
					id: JvmId(value.GetStringBytes("id")),
				}
				session.classes[entry.id] = &entry
				session.entries = append(session.entries, reflect.ValueOf(entry))
			case "jvmmethod":
				avalues := value.GetArray("a")
				var a = make([]string, len(avalues))
				for i, avalue := range avalues {
					a[i] = string(avalue.GetStringBytes())
				}
				entry := JvmMethod{
					Entry: Entry{
						raw:     line,
						tag:     tag,
						session: &session,
					},
					mn: string(value.GetStringBytes("mn")),
					a:  a,
					r:  string(value.GetStringBytes("r")),
					id: JvmId(value.GetStringBytes("id")),
				}
				session.methods[entry.id] = &entry
				session.entries = append(session.entries, reflect.ValueOf(entry))
			case "jvmcall":
				avalues := value.GetArray("av")
				var av = make([]string, len(avalues))
				for i, avalue := range avalues {
					av[i] = string(avalue.GetStringBytes())
				}
				stvalues := value.GetArray("av")
				var st = make([]string, len(stvalues))
				for i, stvalue := range stvalues {
					st[i] = string(stvalue.GetStringBytes())
				}
				entry := JvmCall{
					Entry: Entry{
						raw:     line,
						tag:     tag,
						session: &session,
					},
					cn: string(value.GetStringBytes("cn")),
					mn: string(value.GetStringBytes("mn")),
					av: av,
					rv: string(value.GetStringBytes("rv")),
					id: JvmId(value.GetStringBytes("id")),
					st: st,
				}
				session.entries = append(session.entries, reflect.ValueOf(entry))
			case "":
				entry := JvmReturn{
					Entry: Entry{
						raw:     line,
						tag:     tag,
						session: &session,
					},
					id: JvmId(value.GetStringBytes("id")),
					rv: string(value.GetStringBytes("rv")),
				}
				session.entries = append(session.entries, reflect.ValueOf(entry))
			}
		}
	}

	return &session, err
}
