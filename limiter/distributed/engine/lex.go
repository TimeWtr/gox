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
			// token is space
			pos++
		case ch == '(':
			// token is left bracket
			tokens = append(tokens, Token{Tp: TokenLParen, Value: "("})
			pos++
		case ch == ')':
			// token is right bracket
			tokens = append(tokens, Token{Tp: TokenRParen, Value: ")"})
			pos++
		case unicode.IsDigit(ch) || ch == '.':
			// token is number or dot
			start := pos
			hasDot := ch == '.'
			pos++
			for pos < len(content) && (unicode.IsDigit(rune(content[pos])) || content[pos] == '.') {
				if content[pos] == '.' {
					if hasDot {
						return tokens, fmt.Errorf("unexpected digit at position %d", pos)
					}
					hasDot = true
				}
				pos++
			}

			// Abnormal situation, the number ends with a dot.
			if content[pos] == '.' {
				return tokens, fmt.Errorf("unexpected digit at position %d", pos)
			}

			tokens = append(tokens, Token{Tp: TokenNumber, Value: content[start:pos]})
		case unicode.IsLetter(ch) || ch == '_':
			// token is letter or _.
			start := pos
			pos++
			for pos < len(content) && (unicode.IsLetter(rune(content[pos])) || content[pos] == '_') {
				pos++
			}
			value := content[start:pos]
			upperValue := strings.ToUpper(value)
			if upperValue == "OR" || upperValue == "AND" {
				tokens = append(tokens, Token{Tp: TokenOperator, Value: upperValue})
			} else {
				tokens = append(tokens, Token{Tp: TokenIdentifier, Value: upperValue})
			}
		default:
			return tokens, fmt.Errorf("token (%v) is not valid", ch)
		}
	}

	return tokens, nil
}
