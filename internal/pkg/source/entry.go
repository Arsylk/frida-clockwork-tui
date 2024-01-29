package source

import (
	"fmt"
	"reflect"
)

type ReflectEntry interface {
}

type IEntry interface {
	GetSession() *Session
	Raw() *string
	Tag() *string
	Render() string
	GetRendered() *string
}

type AnyEntry interface {
	Entry | JvmClass | JvmMethod | JvmCall | JvmReturn
}

type Entry struct {
	session  *Session
	raw      string
	tag      string
	rendered *string
}

type Entries []reflect.Value

func (e Entry) GetSession() *Session {
	return e.session
}
func (e Entry) Raw() *string {
	return &e.raw
}
func (e Entry) Tag() *string {
	return &e.tag
}
func (e Entry) GetRendered() *string {
	return e.rendered
}

func (e *Entries) Filter(predicate func(IEntry, int) bool) *Entries {
	length := len(*e)
	arr := make(Entries, 0, length)
	for i, rv := range *e {
		entry := e.Get(i)
		if predicate(entry, i) {
			arr = append(arr, rv)
		}
	}

	return &arr
}

// func (e *Entries) MapToContent() string {
// 	var sb strings.Builder
// 	length := len(*e)
// 	for i := 0; i < length; i += 1 {
// 		entry := e.Get(i)
// 		if entry.GetRendered() == nil {
// 			str := fmt.Sprintf(" %*d │ %s", length, i, RenderEntry(e.Get(i)))
// 			ptrRef := reflect.ValueOf(e.GetReflect(i))
// 			reflect.Indirect(ptrRef).FieldByName("rendered").SetString(str)

// 		}
// 		sb.WriteString(*entry.GetRendered())
// 		sb.WriteString("\n")
// 	}
// 	return sb.String()
// }

func (e *Entries) Get(i int) IEntry {
	rv := *e
	deref := rv[i]
	return deref.Interface().(IEntry)
}
func (e *Entries) GetReflect(i int) *reflect.Value {
	rv := *e
	return &rv[i]
}
func RenderEntry(e ReflectEntry) string {
	switch e := e.(type) {
	case interface{ Render() string }:
		return e.Render()
	default:
		return fmt.Sprintf("%T", e)
	}
}

func RawEntry(e ReflectEntry) string {
	switch e := e.(type) {
	case interface{ Raw() *string }:
		return *e.Raw()
	default:
		return fmt.Sprintf("%T", e)
	}
}
