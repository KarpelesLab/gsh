package gsh

import (
	"bufio"
	"bytes"
	"io"
	"log"
)

type parser struct {
	session   *Session
	buf       bytes.Buffer
	mode      []int
	filename  string
	line, col int
	bio       *bufio.Reader
}

type command struct {
	filename string
	line     int
}

func runScript(ctx *Context, r io.Reader, filename string) error {
	p := &parser{session: ctx.Session, line: 1}

	// if v is already a bufio.Reader, use it as is
	switch v := r.(type) {
	case *bufio.Reader:
		p.bio = v
	default:
		p.bio = bufio.NewReader(v)
	}

	return p.run()
}

func (p *parser) run() error {
	for {
		tok, err := p.readToken()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		log.Printf("tok = %+v", tok)
		return nil // XXX
	}
	return nil
}

func (p *parser) readToken() (Token, error) {
	t := Token{}
	buf := &bytes.Buffer{}
	line, col := p.line, p.col

	for {
		r, _, err := p.readRune()
		if err != nil {
			if err == io.EOF {
				if len(t) == 0 {
					return nil, io.EOF
				}
			}
			return nil, err
		}

		if isShellSafe(r) {
			// that's a standard bash char
			buf.WriteRune(r)
			continue
		}

		if buf.Len() > 0 {
			// flush buf if not empty
			s := stringElement{value: buf.String(), filename: p.filename, line: line, col: col}
			t = append(t, s)
			// reset line, col
			line, col = p.line, p.col
		}

		switch r {
		case ' ', '\t', '\r', '\n':
			if len(t) > 0 {
				return t, nil
			}
			// haven't reached start of token yet, keep reading (and set line, col forward)
			line, col = p.line, p.col
			continue
		case '\'':
			s, err := p.readSingleQuote(line, col)
			if err != nil {
				return nil, err
			}
			t = append(t, s)
		}
	}
}

func (p *parser) readSingleQuote(line, col int) (*stringElement, error) {
	// in single quotes, everything until the next singlequote is part of the string
	buf := &bytes.Buffer{}
	for {
		b, err := p.readByte()
		if err != nil {
			if err == io.EOF {
				return nil, io.ErrUnexpectedEOF
			}
			return nil, err
		}
		if b != '\'' {
			buf.WriteByte(b)
			continue
		}
		// end of string
		return &stringElement{value: buf.String(), filename: p.filename, line: line, col: col}, nil
	}
}

func (p *parser) readByte() (byte, error) {
	b, err := p.bio.ReadByte()
	if err != nil {
		return b, err
	}
	if b == '\n' {
		p.line += 1
		p.col = 0
	} else {
		p.col += 1
	}
	return b, nil
}

func (p *parser) readRune() (rune, int, error) {
	r, l, err := p.bio.ReadRune()
	if err != nil {
		return r, l, err
	}
	if r == '\n' {
		p.line += 1
		p.col = 0
	} else {
		p.col += 1
	}
	return r, l, nil
}
