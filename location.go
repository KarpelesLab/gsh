package gsh

import "os"

type Location interface {
	String() string
	Chdir(dest string) error
}

func ProcLocation() Location {
	return procLocation{}
}

type procLocation struct{}

func (p procLocation) String() string {
	res, err := os.Getwd()
	if err != nil {
		return "." // shouldn't happen
	}
	return res
}

func (p procLocation) Chdir(dest string) error {
	return os.Chdir(dest)
}
