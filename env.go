package gsh

type Environ []Variable

type Variable interface {
	String() string
}

type StringVar struct {
	Name  string
	Value string
}

func (s StringVar) String() string {
	return s.Name + "=" + s.Value
}
