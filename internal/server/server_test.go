package server

import (
	"fmt"
	"testing"
)

func FuzzCall(f *testing.F) {
	testcases := []string{"Hello, world", " ", "!12345"}
	for _, tc := range testcases {
		f.Add(tc) // Use f.Add to provide a seed corpus
	}
	f.Fuzz(func(t *testing.T, orig string) {
		// TODO: start api and call it
		fmt.Println(orig)

		if orig == "" {
			t.Errorf("Before: %q, after: %q", orig)
		}
	})
}
