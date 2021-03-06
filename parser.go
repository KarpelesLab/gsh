package gsh

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"unicode/utf8"
)

type parser struct {
	session   *Session
	mode      []int
	filename  string
	line, col int
	bio       *bufio.Reader
}

type pbuffer struct {
	bytes.Buffer
	p         *parser
	line, col int
}

const DefaultEscape = ";\n"

func runScript(ctx *Context, r io.Reader, filename string) error {
	return ctx.Session.newParser(r, filename).run()
}

func (p *parser) run() error {
	for {
		cmd, _, err := p.readCommand(DefaultEscape)
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		log.Printf("tok = %+v len=%d", cmd.tokens, len(cmd.tokens))
		return nil // XXX
	}
	return nil
}

func (p *parser) buf() *pbuffer {
	return &pbuffer{p: p, line: p.line, col: p.col}
}

func (b *pbuffer) reset() {
	b.line, b.col = b.p.line, b.p.col
	b.Buffer.Reset()
}

func (b *pbuffer) value() stringElement {
	v := stringElement{value: b.String(), filename: b.p.filename, line: b.line, col: b.col}
	b.reset()
	return v
}

func (b *pbuffer) app(t *Token) {
	if b.Len() == 0 {
		return
	}
	*t = append(*t, b.value())
}

func (p *parser) readCommand(escape string) (*command, escapeElement, error) {
	cmd := &command{}

	for {
		tok, err := p.readToken(escape)
		if err != nil {
			if err == io.EOF {
				if len(cmd.tokens) > 0 {
					return cmd, "", nil
				}
			}
			return nil, "", err
		}
		if len(tok) == 0 {
			continue
		}
		switch esc := tok[len(tok)-1].(type) {
		case escapeElement:
			if len(tok) > 1 {
				// append if the end element was not alone, but skip end element
				cmd.tokens = append(cmd.tokens, tok[:len(tok)-1])
			}
			// so which end element did we get?
			switch esc {
			case "|":
				// TODO pipe
				panic("TODO handle pipe")
			default:
				if len(cmd.tokens) > 0 {
					return cmd, esc, nil
				}
			}
		default:
			cmd.tokens = append(cmd.tokens, tok)
		}
	}
}

func (p *parser) readToken(escape string) (Token, error) {
	t := Token{}
	buf := p.buf()
	escapeAs, ok := makeASCIISet(escape)
	if !ok {
		return nil, errors.New("invalid escape char requested")
	}

	for {
		r, _, err := p.readRune()
		if err != nil {
			if err == io.EOF {
				if buf.Len() > 0 {
					t = append(t, buf.value())
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
			buf.app(&t)
		}

		if r < utf8.RuneSelf && escapeAs.contains(byte(r)) {
			t = append(t, escapeElement(string([]rune{r})))
			return t, nil
		}

		switch r {
		case ' ', '\t', '\r':
			if len(t) > 0 {
				return t, nil
			}
			// haven't reached start of token yet, keep reading (and set line, col forward)
			buf.reset()
		case '\n', ';':
			t = append(t, escapeElement(string([]rune{r})))
			return t, nil
		case '|':
			if len(t) == 0 {
				return nil, fmt.Errorf("%w `|`", ErrSyntaxUnexpectedToken)
			}
			v, err := p.bio.Peek(1)
			if err != nil || len(v) == 0 {
				// failed, at least return this one
				t = append(t, escapeElement("|"))
				return t, nil
			}
			if v[1] == '|' {
				t = append(t, escapeElement("||"))
				return t, nil
			}
			t = append(t, escapeElement("|"))
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
			err := p.readSingleQuote(buf)
			if err != nil {
				return nil, err
			}
			t = append(t, buf.value())
		case '"':
			s, err := p.readDoubleQuote()
			if err != nil {
				return nil, err
			}
			t = append(t, s...)
		case '$':
			// ok what is next?
			e, err := p.readVarCall()
			if err != nil {
				if err == io.ErrUnexpectedEOF {
					buf.WriteByte('$')
					t = append(t, buf.value())
					return t, nil
				}
				return nil, err
			}
			t = append(t, e)
		case '#':
			// comment
			err := p.skipLine()
			if err != nil {
				if err == io.EOF {
					return t, nil
				}
				return nil, err
			}
			if len(t) > 0 {
				return t, nil
			}
		case '(':
			if len(t) == 0 {
				// a single ( acts like a { when not preceded by a function name
			}
			// stringElement
			// function declaration (needs something before)
			log.Printf("v = %+v", t)
			panic(fmt.Sprintf("hi"))
		}
	}
}

func (p *parser) readSingleQuote(buf *pbuffer) error {
	// https://www.gnu.org/software/bash/manual/html_node/Single-Quotes.html
	// in single quotes, everything until the next singlequote is part of the string

	for {
		b, err := p.readByte()
		if err != nil {
			if err == io.EOF {
				return io.ErrUnexpectedEOF
			}
			return err
		}
		if b != '\'' {
			buf.WriteByte(b)
			continue
		}
		// end of string
		return nil
	}
}

func (p *parser) readDoubleQuote() (Token, error) {
	// https://www.gnu.org/software/bash/manual/html_node/Double-Quotes.html
	var t Token
	buf := p.buf()

	for {
		r, _, err := p.readRune()
		if err != nil {
			return nil, notEOF(err)
		}

		switch r {
		case '$':
			buf.app(&t)
			v, err := p.readVarCall()
			if err != nil {
				return nil, err
			}
			t = append(t, v)
		case '`':
			buf.app(&t)
			v, err := p.readBacktickCall()
			if err != nil {
				return nil, err
			}
			t = append(t, v)
		case '\\':
			// The backslash retains its special meaning only when followed by one of the following characters: ???$???, ???`???, ???"???, ???\???, or newline
			r, _, err = p.readRune()
			if err != nil {
				return nil, notEOF(err)
			}
			if r == '\r' {
				r, _, err = p.readRune()
				if err != nil {
					return nil, notEOF(err)
				}
			}

			switch r {
			case '$', '`', '"':
				buf.WriteRune(r)
			case '\n':
				// do nothing
			default:
				// no effect
				buf.WriteRune('\\')
				buf.WriteRune(r)
			}
		case '"':
			// end of string
			buf.app(&t)
			return t, nil
		default:
			buf.WriteRune(r)
		}
	}
}

func (p *parser) readVarCall() (Element, error) {
	// this can be a lot of things...
	// $VAR ??? *varElement
	// ${VAR...} ??? *varElement
	// $(cmd) ??? *shellCallElement
	// $'something' (ANSI-C quoting) ??? stringElement
	// $"something" (gettext) ??? ...

	buf := p.buf()

	r, _, err := p.readRune()
	if err != nil {
		return nil, notEOF(err)
	}

	switch r {
	case '\'':
		// https://www.gnu.org/software/bash/manual/html_node/ANSI_002dC-Quoting.html
		// ANSI-C Quoting
		err := p.readSingleQuote(buf)
		if err != nil {
			return nil, err
		}
		v := buf.value()
		v.value, _ = handleEscapes(v.value)
		return v, nil
	default:
		// TODO
	}
	return nil, errors.New("TODO")
}

func (p *parser) readBacktickCall() (*shellCallElement, error) {
	return nil, errors.New("TODO")
}

func (p *parser) readByte() (byte, error) {
	b, err := p.bio.ReadByte()
	if err != nil {
		return b, err
	}
	if b == '\n' {
		p.line += 1
		p.col = 1
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
		p.col = 1
	} else {
		p.col += 1
	}
	return r, l, nil
}

func (p *parser) skipLine() error {
	// read to EOL
	_, err := p.bio.ReadString('\n')
	if err != nil {
		return err
	}
	p.line += 1
	p.col = 1
	return nil
}
