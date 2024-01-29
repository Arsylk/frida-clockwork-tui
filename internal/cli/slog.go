package main

import (
	"fmt"
	"os"
	"reflect"
)

type Base struct {
	Tag string
}

type Item1 struct {
	Base
	name string
}

type Item2 struct {
	Base
	id int
}

func (i Item1) Test() int {
	return 1
}
func (i Item2) Test() int {
	return 2
}

type Testable struct{ Base }

func main() {
	item1 := Item1{
		Base: Base{Tag: ""},
		name: "",
	}
	item2 := Item2{
		Base: Base{Tag: ""},
		id:   0,
	}

	slice := []reflect.Value{reflect.ValueOf(item1), reflect.ValueOf(item2)}

	for _, item := range slice {
		field := item.FieldByIndex(0).Interface()
		fmt.Fprintf(os.Stdout, "%T\n", field)
		a := item.Interface()
		switch a := a.(type) {
		case Testable:
			fmt.Fprintf(os.Stdout, "test %s\n", a.Tag)
		case Item1:
			fmt.Fprintf(os.Stdout, "item1 %s\n", a.Tag)
		case Item2:
			fmt.Fprintf(os.Stdout, "item12 %s\n", a.Tag)
		}
	}
}
