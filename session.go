package gsh

import "os"

type Session struct {
	Location
}

type Job struct {
	os.Process
}
