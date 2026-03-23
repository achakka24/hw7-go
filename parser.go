package main

import (
	"fmt"
	"strconv"
	"unicode"
)

type tokenKind int

const (
	tokenEOF tokenKind = iota
	tokenInt
	tokenIdent
	tokenPlus
	tokenMinus
	tokenStar
	tokenSlash
	tokenCaret
	tokenLParen
	tokenRParen
)

type token struct {
	kind  tokenKind
	text  string
	value int
}

type parser struct {
	tokens []token
	pos    int
}

func parse(input string) (Expr, error) {
	tokens, err := tokenize(input)
	if err != nil {
		return nil, err
	}
	p := &parser{tokens: tokens}
	expr, err := p.parseExpression()
	if err != nil {
		return nil, err
	}
	if p.peek().kind != tokenEOF {
		return nil, fmt.Errorf("unexpected token %q", p.peek().text)
	}
	return expr, nil
}

func tokenize(input string) ([]token, error) {
	tokens := make([]token, 0, len(input)+1)
	runes := []rune(input)
	for i := 0; i < len(runes); {
		r := runes[i]
		if unicode.IsSpace(r) {
			i++
			continue
		}
		switch r {
		case '+':
			tokens = append(tokens, token{kind: tokenPlus, text: "+"})
			i++
		case '-':
			tokens = append(tokens, token{kind: tokenMinus, text: "-"})
			i++
		case '*':
			tokens = append(tokens, token{kind: tokenStar, text: "*"})
			i++
		case '/':
			tokens = append(tokens, token{kind: tokenSlash, text: "/"})
			i++
		case '^':
			tokens = append(tokens, token{kind: tokenCaret, text: "^"})
			i++
		case '(':
			tokens = append(tokens, token{kind: tokenLParen, text: "("})
			i++
		case ')':
			tokens = append(tokens, token{kind: tokenRParen, text: ")"})
			i++
		default:
			if unicode.IsDigit(r) {
				start := i
				for i < len(runes) && unicode.IsDigit(runes[i]) {
					i++
				}
				text := string(runes[start:i])
				n, err := strconv.Atoi(text)
				if err != nil {
					return nil, fmt.Errorf("invalid integer %q", text)
				}
				tokens = append(tokens, token{kind: tokenInt, text: text, value: n})
				continue
			}
			if unicode.IsLetter(r) {
				start := i
				for i < len(runes) && (unicode.IsLetter(runes[i]) || unicode.IsDigit(runes[i]) || runes[i] == '_') {
					i++
				}
				text := string(runes[start:i])
				tokens = append(tokens, token{kind: tokenIdent, text: text})
				continue
			}
			return nil, fmt.Errorf("invalid character %q", string(r))
		}
	}
	tokens = append(tokens, token{kind: tokenEOF, text: "<eof>"})
	return tokens, nil
}

func (p *parser) parseExpression() (Expr, error) {
	return p.parseAddSub()
}

func (p *parser) parseAddSub() (Expr, error) {
	left, err := p.parseMulDiv()
	if err != nil {
		return nil, err
	}
	for {
		switch p.peek().kind {
		case tokenPlus:
			p.next()
			right, err := p.parseMulDiv()
			if err != nil {
				return nil, err
			}
			left = Add{Left: left, Right: right}
		case tokenMinus:
			p.next()
			right, err := p.parseMulDiv()
			if err != nil {
				return nil, err
			}
			left = Sub{Left: left, Right: right}
		default:
			return left, nil
		}
	}
}

func (p *parser) parseMulDiv() (Expr, error) {
	left, err := p.parsePow()
	if err != nil {
		return nil, err
	}
	for {
		switch p.peek().kind {
		case tokenStar:
			p.next()
			right, err := p.parsePow()
			if err != nil {
				return nil, err
			}
			left = Mul{Left: left, Right: right}
		case tokenSlash:
			p.next()
			right, err := p.parsePow()
			if err != nil {
				return nil, err
			}
			left = Div{Left: left, Right: right}
		default:
			return left, nil
		}
	}
}

func (p *parser) parsePow() (Expr, error) {
	base, err := p.parseUnary()
	if err != nil {
		return nil, err
	}
	if p.peek().kind == tokenCaret {
		p.next()
		exp, err := p.parsePow()
		if err != nil {
			return nil, err
		}
		return Pow{Base: base, Exponent: exp}, nil
	}
	return base, nil
}

func (p *parser) parseUnary() (Expr, error) {
	if p.peek().kind == tokenMinus {
		p.next()
		x, err := p.parseUnary()
		if err != nil {
			return nil, err
		}
		return Neg{X: x}, nil
	}
	return p.parsePrimary()
}

func (p *parser) parsePrimary() (Expr, error) {
	t := p.peek()
	switch t.kind {
	case tokenInt:
		p.next()
		return Const{Value: t.value}, nil
	case tokenIdent:
		p.next()
		return Var{Name: t.text}, nil
	case tokenLParen:
		p.next()
		x, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		if p.peek().kind != tokenRParen {
			return nil, fmt.Errorf("expected ')'")
		}
		p.next()
		return x, nil
	default:
		return nil, fmt.Errorf("expected primary expression, got %q", t.text)
	}
}

func (p *parser) peek() token {
	return p.tokens[p.pos]
}

func (p *parser) next() token {
	t := p.tokens[p.pos]
	p.pos++
	return t
}
