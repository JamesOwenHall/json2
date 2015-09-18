package json2

import (
	"fmt"
)

// ParseError represents the presence of an unexpected token.
type ParseError struct {
	Found *Token
}

func (p *ParseError) Error() string {
	if p.Found == nil {
		return "Unexpected end of input."
	}

	return fmt.Sprintf("Unexpected token %s", p.Found.String())
}

// Unmarshal parses the JSON string.
func Unmarshal(input string) (interface{}, error) {
	tokenizer := NewTokenizer(input)
	tokens, err := tokenizer.All()
	if err != nil {
		return nil, err
	}

	p := parser{tokens, 0}
	return p.parse()
}

type parser struct {
	tokens []*Token
	pos    int
}

func (p *parser) parse() (interface{}, *ParseError) {
	v, err := p.parseValue()
	if err != nil {
		return nil, err
	}

	if p.pos != len(p.tokens) {
		return nil, &ParseError{Found: p.tokens[p.pos]}
	}

	return v, nil
}

func (p *parser) parseValue() (interface{}, *ParseError) {
	if p.pos == len(p.tokens) {
		return nil, &ParseError{Found: nil}
	}

	switch tok := p.peek(); tok.Type {
	case Null:
		p.read()
		return nil, nil
	case Boolean, Number, String:
		p.read()
		return tok.Value, nil
	case LSquare:
		return p.parseArray()
	case LCurly:
		return p.parseObject()
	default:
		return nil, &ParseError{Found: tok}
	}
}

func (p *parser) parseArray() ([]interface{}, *ParseError) {
	var res []interface{}

	// Read '['.
	p.read()

	// Read first element.
	tok := p.peek()
	if tok != nil && tok.Type != RSquare {
		v, err := p.parseValue()
		if err != nil {
			return nil, err
		}
		res = append(res, v)
	}

	// Read the rest of the elements.
	tok = p.peek()
	for tok != nil && tok.Type == Comma {
		p.read()

		v, err := p.parseValue()
		if err != nil {
			return nil, err
		}

		res = append(res, v)
		tok = p.peek()
	}

	// Read ']'.
	tok = p.read()
	if tok == nil || tok.Type != RSquare {
		return nil, &ParseError{Found: tok}
	}

	return res, nil
}

func (p *parser) parseObject() (map[string]interface{}, *ParseError) {
	res := make(map[string]interface{})

	// Read '{'.
	p.read()

	// Read the first pair.
	tok := p.peek()
	if tok != nil && tok.Type != RCurly {
		k, v, err := p.parsePair()
		if err != nil {
			return nil, err
		}
		res[k] = v
	}

	// Read the rest of the pairs.
	tok = p.peek()
	for tok != nil && tok.Type == Comma {
		p.read()

		k, v, err := p.parsePair()
		if err != nil {
			return nil, err
		}
		res[k] = v

		tok = p.peek()
	}

	// Read '}'.
	tok = p.read()
	if tok == nil || tok.Type != RCurly {
		return nil, &ParseError{Found: tok}
	}

	return res, nil
}

func (p *parser) parsePair() (string, interface{}, *ParseError) {
	// Read key.
	tok := p.read()
	if tok == nil || tok.Type != String {
		return "", nil, &ParseError{Found: tok}
	}
	key := tok.Value.(string)

	// Read colon.
	tok = p.read()
	if tok == nil || tok.Type != Colon {
		return "", nil, &ParseError{Found: tok}
	}

	// Read value.
	value, err := p.parseValue()
	if err != nil {
		return "", nil, err
	}

	return key, value, nil
}

func (p *parser) peek() *Token {
	if p.pos == len(p.tokens) {
		return nil
	}

	return p.tokens[p.pos]
}

func (p *parser) read() *Token {
	tok := p.peek()
	if tok != nil {
		p.pos++
	}
	return tok
}
