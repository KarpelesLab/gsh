package gsh

import "strings"

type Token []Element

type Element interface {
	Resolve(ctx *Context) (string, error)
	ToSource() string
}

type stringElement struct {
	value     string
	filename  string
	line, col int
}

func (s stringElement) Resolve(ctx *Context) (string, error) {
	return s.value, nil
}

func (s stringElement) ToSource() string {
	return Quote(s.value)
}

type varElement struct {
	// ...
	v string
}

func (v *varElement) Resolve(ctx *Context) (string, error) {
	// TODO resolve this
	return "TODO", nil
}

func (v *varElement) ToSource() string {
	// something like that
	return "${" + v.v + "}"
}

type shellCallElement struct {
	cmd []string
}

func (v *shellCallElement) Resolve(ctx *Context) (string, error) {
	// TODO call method
	return "TODO", nil
}

func (v *shellCallElement) ToSource() string {
	return "$(" + strings.Join(v.cmd, " ") + ")"
}

type newlineElement struct{}

func (newlineElement) Resolve(ctx *Context) (string, error) {
	return "", nil
}

func (newlineElement) ToSource() string {
	return "\n"
}

type semicolonElement struct{}

func (semicolonElement) Resolve(ctx *Context) (string, error) {
	return "", nil
}

func (semicolonElement) ToSource() string {
	return ";"
}
