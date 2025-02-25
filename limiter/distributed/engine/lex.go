// Copyright 2025 TimeWtr
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package engine

import (
	"fmt"
	"strings"
	"unicode"
)

const (
	LogicAnd      = "and"
	LogicUpperAnd = "AND"
	LogicOr       = "or"
	LogicUpperOr  = "OR"
)

type TokenType int

const (
	TokenIdentifier TokenType = iota
	TokenNumber
	TokenOperator
	TokenLogicalOp
	TokenLParen
	TokenRParen
)

type Token struct {
	Tp    TokenType
	Value string
}

func Lex(content string) ([]Token, error) {
	var tokens []Token
	content = strings.TrimSpace(content)
	pos := 0
	for pos < len(content) {
		ch := rune(content[pos])
		switch {
		case unicode.IsSpace(ch):
			// space content, include
			pos++
			continue
		case ch == '(':
			tokens = append(tokens, Token{TokenLParen, string(ch)})
			pos++
		case ch == ')':
			tokens = append(tokens, Token{TokenRParen, string(ch)})
			pos++
		case unicode.IsDigit(ch), ch == '.':
			start := pos
			pos++
			hasDot := ch == '.'
			for pos < len(content) && (unicode.IsDigit(rune(content[pos])) || content[pos] == '.') {
				if content[pos] == '.' {
					if hasDot {
						return tokens, fmt.Errorf("invalid character in number: %s, position: %d",
							content[start:pos], pos)
					}
					hasDot = true
				}
				pos++
			}

			tokens = append(tokens, Token{TokenNumber, content[start:pos]})
		case ch == '>', ch == '<', ch == '=':
			start := pos
			pos++
			if pos < len(content) && content[pos] == '=' {
				tokens = append(tokens, Token{TokenOperator, content[start : pos+1]})
				pos++
			} else {
				tokens = append(tokens, Token{TokenOperator, string(ch)})
			}
		case unicode.IsLetter(ch), ch == '_':
			start := pos
			pos++
			for pos < len(content) && (unicode.IsLetter(rune(content[pos])) || content[pos] == '_' || unicode.IsDigit(rune(content[pos]))) {
				pos++
			}

			value := content[start:pos]
			upperValue := strings.ToUpper(value)
			switch upperValue {
			case LogicUpperOr, LogicUpperAnd:
				tokens = append(tokens, Token{TokenLogicalOp, upperValue})
			default:
				tokens = append(tokens, Token{TokenIdentifier, value})
			}
		default:
			return tokens, fmt.Errorf("invalid character in identifier: %s, position: %d", string(ch), pos)
		}
	}

	return tokens, nil
}

type Expr interface {
	Evaluate() (bool, error)
}

// LogicalExpr the struct of logical expr, such as cpu_usage > 70 OR mem_usage >= 80
type LogicalExpr struct {
	Operator string // OR 、AND
	Left     Expr
	Right    Expr
}

func (e *LogicalExpr) Evaluate() (bool, error) {
	l, err := e.Left.Evaluate()
	if err != nil {
		return false, err
	}

	r, err := e.Right.Evaluate()
	if err != nil {
		return false, err
	}

	switch e.Operator {
	case LogicAnd, LogicUpperAnd:
		return l && r, nil
	case LogicOr, LogicUpperOr:
		return l || r, nil
	default:
		return false, fmt.Errorf("unsupported logical operator: %s", e.Operator)
	}
}

// Condition the struct of trigger condition, such as cpu_usage > 80, mem_usage >= 80
type Condition struct {
	Field    string
	Operator string // >, >=, <, <=, ==
	Value    float64
}

func (c *Condition) Evaluate() (bool, error) {
	_, ok := metricsMap[c.Field]
	if !ok {
		return false, nil
	}

	return true, nil
}

type ParenExpr struct {
	Ep Expr
}
type TriggerParser struct {
	tokens []Token
	pos    int
}
