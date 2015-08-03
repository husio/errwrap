package errwrap

import (
	"errors"
	"io"
	"testing"
)

var (
	err1 = errors.New("first error")
	err2 = errors.New("second error")
	err3 = errors.New("third error")
)

func TestErrorMessage(t *testing.T) {
	var testCases = []struct {
		err      error
		expected string
	}{
		{
			Err(err1, "first"),
			"first: first error",
		},
		{
			Err(Err(err1, "first"), "second"),
			"second: first: first error",
		},
		{
			Err(Err(Err(err1, "first"), "second"), "third"),
			"third: second: first: first error",
		},
	}

	for i, tc := range testCases {
		if got := tc.err.Error(); got != tc.expected {
			t.Errorf("%d: expected %q, got %q", i, tc.expected, got)
		}
	}
}

func TestIs(t *testing.T) {
	var testCases = []struct {
		err   error
		anyof []error
		is    bool
	}{
		{
			io.EOF,
			[]error{err1, err2, err3},
			false,
		},
		{
			io.EOF,
			[]error{err1, err2, io.EOF, err3},
			true,
		},
		{
			io.EOF,
			[]error{io.EOF, err1, err2, err3},
			true,
		},
		{
			Err(err1, "first"),
			[]error{io.EOF, err2, err1},
			true,
		},
		{
			Err(Err(err1, "first"), "second"),
			[]error{io.EOF, err2, err1},
			true,
		},
		{
			Err(Err(err1, "first"), "second"),
			[]error{io.EOF, err2},
			false,
		},
	}

	for i, tc := range testCases {
		if got := Is(tc.err, tc.anyof...); got != tc.is {
			t.Errorf("%d: expected %v", i, tc.is)
		}
	}
}

func TestErrDoesNotMutate(t *testing.T) {
	origin := Err(err1, "origin")

	e1 := Err(origin, "1")
	e2 := Err(origin, "2")

	if !Is(e1, e2) {
		t.Error("e1 is not e2")
	}

	if e1.Error() != "1: origin: first error" {
		t.Fatalf("invalid e1 error: %q", e1.Error())
	}

	if e2.Error() != "2: origin: first error" {
		t.Fatalf("invalid e2 error: %q", e2.Error())
	}
}
