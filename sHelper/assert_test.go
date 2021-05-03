package sHelper

import "testing"

type fakeT struct {
	iTesting
}

func (fakeT) Helper() {
}

func (fakeT) Fatal(args ...interface{}) {
}

func TestAssertEqual(t *testing.T) {
	AssertEqual(t, "HELLO", "HELLO")
	AssertEqual(&fakeT{}, "HELLO", "WORLD")
}
