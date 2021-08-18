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

type escapeElement string

func (escapeElement) Resolve(ctx *Context) (string, error) {
	return "", nil
}

func (e escapeElement) ToSource() string {
	return string(e)
}
