package main

import (
	"fmt"
)

func main() {
	err := SomeFunction(5)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}

func SomeFunction(input int64) error {
	if input < 0 {
		return fmt.Errorf("input must be non-negative, got %d", input)
	}

	return nil
}