package scan

import (
	"unicode"
)

// Scanner is a func that returns a subset of the input and a success bool.
type Scanner func([]rune) ([]rune, bool)

// If returns a scanner that accepts the a rune if it satisfies the condition.
func If(condition func(rune) bool) Scanner {
	return func(input []rune) ([]rune, bool) {
		if len(input) > 0 && condition(input[0]) {
			return input[0:1], true
		}

		return nil, false
	}
}

// Rune returns a scanner that accepts r.
func Rune(r rune) Scanner {
	return If(func(b rune) bool {
		return r == b
	})
}

// Space returns a scanner that accepts whitespace as defined in the unicode
// package.
func Space() Scanner {
	return func(input []rune) ([]rune, bool) {
		if len(input) > 0 && unicode.IsSpace(input[0]) {
			return input[0:1], true
		}

		return nil, false
	}
}

// And returns a scanner that accepts all scanners in sequence.
func And(scanners ...Scanner) Scanner {
	return func(input []rune) ([]rune, bool) {
		remaining := input
		accumulated := []rune{}

		for _, s := range scanners {
			if read, ok := s(remaining); !ok {
				return nil, false
			} else {
				accumulated = append(accumulated, read...)
				remaining = remaining[len(read):]
			}
		}

		return accumulated, true
	}
}

// Or returns a scanner that accepts the first successful scan in scanners.
func Or(scanners ...Scanner) Scanner {
	return func(input []rune) ([]rune, bool) {
		for _, s := range scanners {
			if read, ok := s(input); ok {
				return read, true
			}
		}

		return nil, false
	}
}

// Any returns a scanner that accepts any number of occurrences of scanner,
// including zero.
func Any(scanner Scanner) Scanner {
	return func(input []rune) ([]rune, bool) {
		remaining := input
		accumulated := []rune{}

		for {
			if read, ok := scanner(remaining); !ok {
				return accumulated, true
			} else {
				accumulated = append(accumulated, read...)
				remaining = remaining[len(read):]
			}
		}
	}
}

// N returns a scanner that accepts scanner exactly n times.
func N(n int, scanner Scanner) Scanner {
	return func(input []rune) ([]rune, bool) {
		scanners := make([]Scanner, n)
		for i := 0; i < n; i++ {
			scanners[i] = scanner
		}

		return And(scanners...)(input)
	}
}

// Maybe runs scanner and returns true regardless of the output.
func Maybe(scanner Scanner) Scanner {
	return func(input []rune) ([]rune, bool) {
		read, _ := scanner(input)
		return read, true
	}
}
