package json2

import (
	"reflect"
	"testing"
)

func TestUnmarshal(t *testing.T) {
	type TestCase struct {
		input string
		exp   interface{}
	}

	tests := []TestCase{
		{`null`, nil},
		{`false`, false},
		{`true`, true},
		{`-1.0`, -1.0},
		{`"foo \" bar"`, `foo " bar`},
		{`["foo", "bar"]`, []interface{}{"foo", "bar"}},
		{`{"foo": "bar", "bar": "baz"}`, map[string]interface{}{"foo": "bar", "bar": "baz"}},
	}

	for _, test := range tests {
		act, err := Unmarshal(test.input)
		if err != (*ParseError)(nil) {
			t.Errorf("Error: %v\nInput: %s", err, test.input)
			continue
		}

		if !reflect.DeepEqual(act, test.exp) {
			t.Errorf("Expected: %#v\nActual: %#v", test.exp, act)
		}
	}
}

func TestUnmarshalError(t *testing.T) {
	type TestCase struct {
		input string
		err   error
	}

	tests := []TestCase{
		{``, &ParseError{Found: nil}},
		{` `, &ParseError{Found: nil}},
		{`nul`, &TokenError{Found: "nul"}},
		{`true false`, &ParseError{Found: &Token{
			Type:  Boolean,
			Value: false,
		}}},
		{`[,]`, &ParseError{Found: &Token{
			Type: Comma,
		}}},
		{`[12,]`, &ParseError{Found: &Token{
			Type: RSquare,
		}}},
		{`[12,13`, &ParseError{Found: nil}},
		{`{`, &ParseError{Found: nil}},
		{`{true:3}`, &ParseError{Found: &Token{
			Type:  Boolean,
			Value: true,
		}}},
		{`{"key":"value", "foo" "value"}`, &ParseError{Found: &Token{
			Type:  String,
			Value: "value",
		}}},
		{`{"key":"value", "foo":}`, &ParseError{Found: &Token{
			Type: RCurly,
		}}},
	}

	for _, test := range tests {
		_, err := Unmarshal(test.input)
		if !reflect.DeepEqual(err, test.err) {
			t.Errorf("Expected: %#v\nActual:%#v", err, test.err)
		}
	}
}
