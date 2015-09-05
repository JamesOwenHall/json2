package scan

import (
	"testing"
)

func TestTokenizer(t *testing.T) {
	type TestCase struct {
		input    string
		expected []*Token
	}

	tests := []TestCase{
		{"", nil},
		{"[{:,}]", []*Token{
			&Token{Type: LSquare},
			&Token{Type: LCurly},
			&Token{Type: Colon},
			&Token{Type: Comma},
			&Token{Type: RCurly},
			&Token{Type: RSquare},
		}},
		{"\n[\t{ :\n, }    ] \t", []*Token{
			&Token{Type: LSquare},
			&Token{Type: LCurly},
			&Token{Type: Colon},
			&Token{Type: Comma},
			&Token{Type: RCurly},
			&Token{Type: RSquare},
		}},
		{"null false true", []*Token{
			&Token{Type: Null},
			&Token{Type: Boolean, Value: false},
			&Token{Type: Boolean, Value: true},
		}},
		{"0 1 -1 2.4 -7.6 45.67 -98.71 0.2e1 0.2e+2 0.2e-3", []*Token{
			&Token{Type: Number, Value: 0.0},
			&Token{Type: Number, Value: 1.0},
			&Token{Type: Number, Value: -1.0},
			&Token{Type: Number, Value: 2.4},
			&Token{Type: Number, Value: -7.6},
			&Token{Type: Number, Value: 45.67},
			&Token{Type: Number, Value: -98.71},
			&Token{Type: Number, Value: 2.0},
			&Token{Type: Number, Value: 20.0},
			&Token{Type: Number, Value: 0.0002},
		}},
		{`"" "foo" "\"\"" "\n\t\r\b\f"`, []*Token{
			&Token{Type: String, Value: ""},
			&Token{Type: String, Value: "foo"},
			&Token{Type: String, Value: `""`},
			&Token{Type: String, Value: "\n\t\r\b\f"},
		}},
	}

	for _, test := range tests {
		tokens, err := NewTokenizer(test.input).All()
		if err != nil {
			t.Errorf("Unexpected error %s from %s", err.Error(), test.input)
			continue
		}

		if len(tokens) != len(test.expected) {
			t.Errorf("Expected %v, got %v from %s", test.expected, tokens, test.input)
			continue
		}

		for i, exp := range test.expected {
			if !exp.Equals(tokens[i]) {
				t.Errorf("Expected %v, got %v from %s", exp, tokens[i], test.input)
			}
		}
	}
}
