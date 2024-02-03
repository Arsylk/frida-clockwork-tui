package source

import (
	"reflect"
)

type IEntry interface {
	Raw() *string
	Tag() *string
	Render() string
	GetRendered() *string
}

type AnyEntry interface {
	Entry | JvmClass | JvmMethod | JvmCall | JvmReturn
}

type Entry struct {
	raw      string
	tag      string
	rendered *string
}

type Entries []reflect.Value

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
