package server

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
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
			t.Errorf("Before: %q", orig)
		}
	})
}

func TestGetHelloWorlod(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/helloworld", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	s := Server{}
	if err := s.GetHelloWorlod(c); err != nil {
		t.Fatalf("GetHelloWorlod() failed: %v", err)
	}

	expected := "HelloWorlod"
	if rec.Body.String() != expected {
		t.Errorf("expected %q but got %q", expected, rec.Body.String())
	}
}
