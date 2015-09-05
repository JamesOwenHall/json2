package scan

import (
	"fmt"
	"strings"
	"unicode"
)

type TokenType int8

func (t TokenType) String() string {
	switch t {
	case LSquare:
		return "["
	case RSquare:
		return "]"
	case LCurly:
		return "{"
	case RCurly:
		return "}"
	case Colon:
		return ":"
	case Comma:
		return ","
	case Null:
		return "null"
	case Boolean:
		return "boolean"
	case Number:
		return "number"
	case String:
		return "string"
	default:
		return "unknown"
	}
}

const (
	LSquare TokenType = iota
	RSquare
	LCurly
	RCurly
	Colon
	Comma
	Null
	Boolean
	Number
	String
)

type Token struct {
	Type  TokenType
	Value interface{}
}

func (t *Token) Equals(other *Token) bool {
	if t == nil && other == nil {
		return true
	}

	if t == nil || other == nil {
		return false
	}

	if t.Type != other.Type {
		return false
	}

	switch tv := t.Value.(type) {
	case float64:
		if ov, ok := other.Value.(float64); !ok {
			return false
		} else {
			return ov < tv+0.000001 && ov > tv-0.000001
		}
	default:
		return t.Value == other.Value
	}
}

func (t *Token) String() string {
	if t == nil {
		return "<nil>"
	}

	return fmt.Sprintf("{%v: %v}", t.Type, t.Value)
}

type TokenError struct {
	Found string
}

func (t *TokenError) Error() string {
	return fmt.Sprintf("Unknown token %s", t.Found)
}

type Tokenizer struct {
	input []rune
}

func NewTokenizer(input string) *Tokenizer {
	return &Tokenizer{input: []rune(input)}
}

func (t *Tokenizer) All() ([]*Token, *TokenError) {
	result := make([]*Token, 0)
	for {
		tok, err := t.Next()
		if err != nil {
			return nil, err
		}

		if tok == nil {
			break
		}

		result = append(result, tok)
	}

	return result, nil
}

func (t *Tokenizer) Next() (*Token, *TokenError) {
	// Skip white space.
	spaces, _ := Any(Space())(t.input)
	t.input = t.input[len(spaces):]

	// Stop at EOF.
	if len(t.input) == 0 {
		return nil, nil
	}

	// Punctuation.
	if _, ok := Rune('[')(t.input); ok {
		t.input = t.input[1:]
		return &Token{Type: LSquare}, nil
	}
	if _, ok := Rune(']')(t.input); ok {
		t.input = t.input[1:]
		return &Token{Type: RSquare}, nil
	}
	if _, ok := Rune('{')(t.input); ok {
		t.input = t.input[1:]
		return &Token{Type: LCurly}, nil
	}
	if _, ok := Rune('}')(t.input); ok {
		t.input = t.input[1:]
		return &Token{Type: RCurly}, nil
	}
	if _, ok := Rune(':')(t.input); ok {
		t.input = t.input[1:]
		return &Token{Type: Colon}, nil
	}
	if _, ok := Rune(',')(t.input); ok {
		t.input = t.input[1:]
		return &Token{Type: Comma}, nil
	}

	// Keywords
	if ident, ok := Identifier()(t.input); ok {
		t.input = t.input[len(ident):]

		switch string(ident) {
		case "true":
			return &Token{Type: Boolean, Value: true}, nil
		case "false":
			return &Token{Type: Boolean, Value: false}, nil
		case "null":
			return &Token{Type: Null}, nil
		default:
			return nil, &TokenError{Found: string(ident)}
		}
	}

	// Number
	if num, ok := Num()(t.input); ok {
		t.input = t.input[len(num):]
		var val float64
		fmt.Sscanf(string(num), "%f", &val)
		return &Token{Type: Number, Value: val}, nil
	}

	// String
	if str, ok := Str()(t.input); ok {
		t.input = t.input[len(str):]
		return &Token{Type: String, Value: string(Unescape(str))}, nil
	}

	return nil, &TokenError{Found: string(t.input[:1])}
}

// Identifier returns a scanner that reads an identifier.  It must start with
// a letter or underscore and it must contain only letters, underscores and
// digits.  JSON doesn't support general identifiers, but this allows the
// tokenizer to return more useful error messages.
func Identifier() Scanner {
	return func(input []rune) ([]rune, bool) {
		return And(
			Or(If(unicode.IsLetter), Rune('_')),
			Any(Or(If(unicode.IsLetter), Rune('_'), If(unicode.IsDigit))),
		)(input)
	}
}

// Num returns a scanner that reads a number.
func Num() Scanner {
	return func(input []rune) ([]rune, bool) {
		return And(
			Maybe(Rune('-')),
			Or(
				Rune('0'),
				And(If(IsDigit19), Any(If(unicode.IsDigit))),
			),
			Maybe(
				And(
					Rune('.'),
					AtLeast(1, If(unicode.IsDigit)),
				),
			),
			Maybe(
				And(
					Or(Rune('e'), Rune('E')),
					Maybe(Or(Rune('+'), Rune('-'))),
					AtLeast(1, If(unicode.IsDigit)),
				),
			),
		)(input)
	}
}

// Str returns a scanner that reads a string.
func Str() Scanner {
	return And(
		Rune('"'),
		Any(
			Or(
				If(func(r rune) bool {
					return r != '\\' && r != '"'
				}),
				And(
					Rune('\\'),
					Or(
						Rune('"'),
						Rune('\\'),
						Rune('/'),
						Rune('b'),
						Rune('f'),
						Rune('n'),
						Rune('r'),
						Rune('t'),
						And(Rune('u'), N(4, If(IsHexDigit))),
					),
				),
			),
		),
		Rune('"'),
	)
}

func Unescape(in []rune) []rune {
	// Trim the double quotation marks.
	in = in[1 : len(in)-1]
	out := make([]rune, 0, len(in))

	for i := 0; i < len(in); i++ {
		if in[i] == '\\' {
			i++
			switch in[i] {
			case '"':
				out = append(out, '"')
			case '\\':
				out = append(out, '\\')
			case '/':
				out = append(out, '/')
			case 'b':
				out = append(out, '\b')
			case 'f':
				out = append(out, '\f')
			case 'n':
				out = append(out, '\n')
			case 'r':
				out = append(out, '\r')
			case 't':
				out = append(out, '\t')
			case 'u':
				hexStr := strings.ToLower(string(in[i : i+4]))
				i += 3
				var hex int
				fmt.Sscanf(hexStr, "%x", &hex)
				out = append(out, rune(hex))
			}
		} else {
			out = append(out, in[i])
		}
	}

	return out
}

func IsDigit19(r rune) bool {
	return '1' <= r && r <= '9'
}

func IsHexDigit(r rune) bool {
	return ('0' <= r && r <= '9') || ('a' <= r && r <= 'f') || ('A' <= r && r <= 'F')
}
