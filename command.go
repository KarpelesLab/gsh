package gsh

import "bytes"

type command struct {
	tokens []Token
}

func (c *command) Argv(ctx *Context) ([]string, error) {
	argv := make([]string, len(c.tokens))

	for argn, argtok := range c.tokens {
		buf := &bytes.Buffer{}
		for _, elem := range argtok {
			v, err := elem.Resolve(ctx)
			if err != nil {
				return nil, err
			}
			buf.WriteString(v)
		}
		argv[argn] = buf.String()
	}
	return argv, nil
}

func (c *command) ToSource() string {
	buf := &bytes.Buffer{}

	for _, argtok := range c.tokens {
		for _, elem := range argtok {
			buf.WriteString(elem.ToSource())
		}
	}

	return buf.String()
}
