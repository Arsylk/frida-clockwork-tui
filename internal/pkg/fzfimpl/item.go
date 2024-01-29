package fzfimpl

import (
	"fmt"
	"reflect"
	"regexp"
)

type Items struct {
	items    reflect.Value
	itemFunc func(i int) string
}

var (
	linesPattern = regexp.MustCompile(`\n+`)
)

// is the items struct even necessary and/or useful ?
func newItems(rv reflect.Value, itemFunc func(i int) string) (*Items, error) {
	return &Items{
		items:    rv,
		itemFunc: itemFunc,
	}, nil
}

// not really sure why
func NewHotItems(items interface{}, itemFunc func(i int) string) (*Items, error) {
	rv := reflect.ValueOf(items)

	if !(rv.Kind() == reflect.Ptr && reflect.Indirect(rv).Kind() == reflect.Array) {
		return nil, fmt.Errorf("items must be a pointer to an array, but got %T", items)
	}

	return newItems(rv, itemFunc)
}

// not really sure why, even less sure why now
func NewColdItems(items interface{}, itemFunc func(i int) string) (*Items, error) {
	rv := reflect.ValueOf(items)

	// if !(rv.Kind() == reflect.Ptr && reflect.Indirect(rv).Kind() == reflect.Slice) {
	// 	return nil, fmt.Errorf("items must be a pointer to an array, but got %T", items)
	// }

	return newItems(rv, itemFunc)
}

func (is Items) ItemString(i int) string {
	return linesPattern.ReplaceAllString(is.itemFunc(i), " ")
}

func (is Items) Len() int {
	if is.items.Kind() == reflect.Ptr {
		return reflect.Indirect(is.items).Len()
	} else {
		return is.items.Len()
	}
}
