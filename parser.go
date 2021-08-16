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

func runScript(ctx *Context, r io.Reader, filename string) error {
	return ctx.Session.newParser(r, filename).run()
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

func (p *parser) readCommand() (*command, error) {
	cmd := &command{}

	for {
		tok, err := p.readToken()
		if err != nil {
			if err == io.EOF {
				if len(cmd.tokens) > 0 {
					return cmd, nil
				}
			}
			return nil, err
		}
		if len(tok) == 0 {
			continue
		}
		switch tok[len(tok)-1].(type) {
		case newlineElement, semicolonElement:
			if len(tok) > 1 {
				// append if the end element was not alone, but skip end element
				cmd.tokens = append(cmd.tokens, tok[:len(tok)-1])
			}
			return cmd, nil
		default:
			cmd.tokens = append(cmd.tokens, tok)
		}
	}
}

func (p *parser) readToken() (Token, error) {
	t := Token{}
	buf := &bytes.Buffer{}
	line, col := p.line, p.col

	for {
		r, _, err := p.readRune()
		if err != nil {
			if err == io.EOF {
				if buf.Len() > 0 {
					s := stringElement{value: buf.String(), filename: p.filename, line: line, col: col}
					t = append(t, s)
				}
				if len(t) > 0 {
					return t, nil
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
			// reset
			buf = &bytes.Buffer{}
			line, col = p.line, p.col
		}

		switch r {
		case ' ', '\t', '\r':
			if len(t) > 0 {
				return t, nil
			}
			// haven't reached start of token yet, keep reading (and set line, col forward)
			line, col = p.line, p.col
		case '\n':
			t = append(t, newlineElement{})
			return t, nil
		case ';':
			t = append(t, semicolonElement{})
			return t, nil
		case '\\':
			// https://www.gnu.org/software/bash/manual/html_node/Escape-Character.html
			// It preserves the literal value of the next character that follows, with the exception of newline
			r, _, err := p.readRune()
			if err != nil {
				if err == io.EOF {
					err = io.ErrUnexpectedEOF
				}
				return nil, err
			}
			if r == '\r' {
				r, _, err = p.readRune()
				if err != nil {
					if err == io.EOF {
						err = io.ErrUnexpectedEOF
					}
					return nil, err
				}
			}
			if r != '\n' {
				buf.WriteRune(r)
			}
		case '\'':
			s, err := p.readSingleQuote(line, col)
			if err != nil {
				return nil, err
			}
			t = append(t, s)
		}
	}
}

func (p *parser) readSingleQuote(line, col int) (stringElement, error) {
	// https://www.gnu.org/software/bash/manual/html_node/Single-Quotes.html
	// in single quotes, everything until the next singlequote is part of the string
	buf := &bytes.Buffer{}
	for {
		b, err := p.readByte()
		if err != nil {
			if err == io.EOF {
				return stringElement{}, io.ErrUnexpectedEOF
			}
			return stringElement{}, err
		}
		if b != '\'' {
			buf.WriteByte(b)
			continue
		}
		// end of string
		return stringElement{value: buf.String(), filename: p.filename, line: line, col: col}, nil
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
