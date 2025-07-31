package main

import "testing"

func TestSomeFunction(t *testing.T) {
	err := SomeFunction(5)
	if err != nil {
		t.Errorf("SomeFunction(5) returned an error: %v", err)
	}
}