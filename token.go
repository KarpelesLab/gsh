package gsh

type Token []Element

type Element interface {
	Resolve(ctx Context) (string, error)
	ToSource() string
}

type stringElement string

func (s stringElement) Resolve(ctx Context) (string, error) {
	return string(s), nil
}

func (s stringElement) ToSource() string {
	return Quote(string(s))
}

type varElement string

func (v varElement) Resolve(ctx Context) (string, error) {
	// TODO resolve this
	return "TODO", nil
}

func (v varElement) ToSource() string {
	// something like that
	return "${" + string(v) + "}"
}
