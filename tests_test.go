package gsh

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"testing"
)

func TestAll(t *testing.T) {
	// run all tests in tests/ except bash.sh
	lst, _ := os.ReadDir("tests")
	for _, f := range lst {
		if f.Name() == "bash.sh" {
			continue
		}

		fn := filepath.Join("tests", f.Name())
		sess := New()
		log.Printf("running %s", fn)

		err := sess.RunFile(context.Background(), fn)
		if err != nil {
			t.Errorf("%s failed: %s", fn, err)
		}
	}
}
