package ignore

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

// ErrRegexCompile is returned when an error occurs while compiling regular
// expressions when parsing a .coverignore file.
var ErrRegexCompile = errors.New("failed to compile regex")

// File represents a .coverignore file and provides the functionality to match
// paths against its rules.
type File struct {
	patterns []*Pattern
}

// New creates a new File instance from a given .coverignore file.
func New(path string) (*File, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	defer file.Close()

	patterns, err := Parse(file)
	if err != nil {
		if errors.Is(err, ErrInvalidRegex) {
			return nil, fmt.Errorf("%w: %w", ErrRegexCompile, err)
		}

		return nil, fmt.Errorf("%w", err)
	}

	return &File{
		patterns: patterns,
	}, nil
}

// NewFromLines creates a new File instance from a list of strings. Useful when
// patterns are available in memory rather than in a file or for testing.
func NewFromLines(lines []string) (*File, error) {
	r := strings.NewReader(strings.Join(lines, "\n"))

	patterns, err := Parse(r)
	if err != nil {
		if errors.Is(err, ErrInvalidRegex) {
			return nil, fmt.Errorf("%w: %w", ErrRegexCompile, err)
		}

		return nil, fmt.Errorf("%w", err)
	}

	return &File{
		patterns: patterns,
	}, nil
}

func (f *File) Patterns() []*Pattern {
	return f.patterns
}

// Match checks if the given path matches any of the .coverignore rules, and
// return true if the path should be ignored according to the rules.
//
// The path is normalized to use forward slashes (/) regardless of the operating
// system.
func (f *File) Match(path string) bool {
	path = strings.ReplaceAll(path, string(os.PathSeparator), "/")

	var match bool

	for _, pat := range f.patterns {
		if pat.Regex.MatchString(path) {
			if pat.Negate {
				return false
			}

			match = true
		}
	}

	return match
}