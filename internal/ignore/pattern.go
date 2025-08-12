package ignore

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strings"
)

// defaultPatternCapacity is the initial capacity allocated for coverignore
// patterns. This value trues to be a reasonable default that covers most common
// .coverignore files without over-allocating.
const defaultPatternCapacity int = 20

var (
	// ErrInvalidRegex is returned when a regular expression fails to compile.
	ErrInvalidRegex = errors.New("invalid regex")

	// ErrScanningFile is returned when scanning a file fails for any reason.
	ErrScanningFile = errors.New("failed to scan file")
)

// Pattern represents a parsed coverignore pattern.
type Pattern struct {
	// Regex is the compiled regular expression for this pattern.
	Regex *regexp.Regexp

	// Negate indicates whether the pattern should be negated.
	Negate bool
}

// Parse parses a .coverignore file into a list of patterns.
func Parse(r io.Reader) ([]*Pattern, error) {
	var (
		lineNumber int
		builder    strings.Builder
		patterns   = make([]*Pattern, 0, defaultPatternCapacity)
		scanner    = bufio.NewScanner(r)
	)

	for scanner.Scan() {
		lineNumber++

		line := scanner.Text()

		// Trim OS-specific carriage returns.
		line = strings.TrimRight(line, "\r")

		// Strip comments [Rule 2].
		if strings.HasPrefix(line, `#`) {
			continue
		}

		// Trim string [Rule 3].
		line = strings.Trim(line, " ")

		// Exit for no-ops and return nil which will prevent us from
		// appending a pattern against this line.
		if line == "" {
			continue
		}

		// Handle [Rule 4] which negates the match for patterns leading with "!".
		negatePattern := false
		if strings.HasPrefix(line, "!") {
			negatePattern = true

			line = line[1:]
		}

		// Handle [Rule 2, 4], when # or ! is escaped with a \.
		if regexp.MustCompile(`^([#!])`).MatchString(line) {
			line = line[1:]
		}

		// If we encounter a foo/*.blah in a folder, prepend the / char.
		if regexp.MustCompile(`([^/+])/.*\*\.`).MatchString(line) && !strings.HasPrefix(line, "/") {
			line = "/" + line
		}

		// Handle escaping the "." char.
		line = regexp.MustCompile(`\.`).ReplaceAllString(line, `\.`)

		const magicStar = "#$~"

		// Handle "/**/" usage.
		if strings.HasPrefix(line, "/**/") {
			line = line[1:]
		}

		line = regexp.MustCompile(`/\*\*/`).ReplaceAllString(line, `(/|/.+/)`)
		line = regexp.MustCompile(`\*\*/`).ReplaceAllString(line, `(|.`+magicStar+`/)`)
		line = regexp.MustCompile(`/\*\*`).ReplaceAllString(line, `(|/.`+magicStar+`)`)

		// Handle escaping the "*" char.
		line = regexp.MustCompile(`\\\*`).ReplaceAllString(line, `\`+magicStar)
		line = regexp.MustCompile(`\*`).ReplaceAllString(line, `([^/]*)`)

		// Handle escaping the "?" char.
		line = strings.ReplaceAll(line, "?", `\?`)

		line = strings.ReplaceAll(line, magicStar, "*")

		builder.Reset()

		if strings.HasSuffix(line, "/") {
			builder.WriteString(line)
			builder.WriteString("(|.*)$")
		} else {
			builder.WriteString(line)
			builder.WriteString("(|/.*)$")
		}

		expr := builder.String()

		if strings.HasPrefix(expr, "/") {
			expr = "^(|/)" + expr[1:]
		} else {
			expr = "^(|.*/)" + expr
		}

		regex, err := regexp.Compile(expr)
		if err != nil {
			return nil, fmt.Errorf("%w: %q on line %d: %w", ErrInvalidRegex, expr, lineNumber, err)
		}

		patterns = append(patterns, &Pattern{
			Regex:  regex,
			Negate: negatePattern,
		})
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrScanningFile, err)
	}

	return patterns, nil
}