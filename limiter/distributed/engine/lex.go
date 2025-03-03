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
	"strconv"
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

func (t TokenType) String() string {
	switch t {
	case TokenIdentifier:
		return "identifier"
	case TokenNumber:
		return "number"
	case TokenOperator:
		return "operator"
	case TokenLogicalOp:
		return "logical op"
	case TokenLParen:
		return "lparen"
	case TokenRParen:
		return "rparen"
	default:
		return "unknown"
	}
}

var (
	metricsMap = map[string]struct{}{
		"cpu_usage":       {},
		"mem_usage":       {},
		"err_rate":        {},
		"mem_used":        {},
		"request_latency": {},
		"active_conns":    {},
	}

	operatorsMap = map[string]struct{}{
		">":  {},
		"<":  {},
		">=": {},
		"<=": {},
		"=":  {},
	}
)

type Token struct {
	Tp    TokenType
	Value string
}

func lex(content string) ([]Token, error) {
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

type NodeType int

const (
	NodeLogical NodeType = iota
	NodeCondition
)

func (t *NodeType) String() string {
	switch *t {
	case NodeLogical:
		return "logical Node"
	case NodeCondition:
		return "condition Node"
	default:
		return "unknown node"
	}
}

type EvalContext struct {
	metrics map[string]float64
}

func WithEvalContext(metrics map[string]float64) EvalContext {
	return EvalContext{metrics: metrics}
}

type Expr interface {
	// GetType get the type of node
	GetType() NodeType
	// GetOperator get the operator only if node type is logical.
	GetOperator() string
	// GetChildren get the children expression.
	GetChildren() []Expr
	// GetCondition get the condition only if node type is condition.
	GetCondition() *Condition
	// Evaluate the core logic to evaluate metric value.
	Evaluate(EvalContext) (bool, error)
	// String return the value type string.
	String() string
}

var _ Expr = (*LogicalExpr)(nil)

// LogicalExpr the struct of logical expr, such as cpu_usage > 70 OR mem_usage >= 80
type LogicalExpr struct {
	Operator string // OR 、AND
	Left     Expr
	Right    Expr
}

func (e *LogicalExpr) GetType() NodeType {
	return NodeLogical
}

func (e *LogicalExpr) GetOperator() string {
	return e.Operator
}

func (e *LogicalExpr) GetChildren() []Expr {
	return []Expr{e.Left, e.Right}
}

func (e *LogicalExpr) GetCondition() *Condition {
	return nil
}

func (e *LogicalExpr) String() string {
	return fmt.Sprintf("%s %s %s", e.Left, e.Operator, e.Right)
}

func (e *LogicalExpr) Evaluate(ctx EvalContext) (bool, error) {
	l, err := e.Left.Evaluate(ctx)
	if err != nil {
		return false, err
	}

	r, err := e.Right.Evaluate(ctx)
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
	Operator string // >, >=, <, <=, =
	Value    float64
}

var _ Expr = (*Condition)(nil)

func (c *Condition) GetType() NodeType {
	return NodeCondition
}

func (c *Condition) GetOperator() string {
	return ""
}

func (c *Condition) GetChildren() []Expr {
	return nil
}

func (c *Condition) GetCondition() *Condition {
	return c
}

func (c *Condition) String() string {
	return fmt.Sprintf("%s %s %f", c.Field, c.Operator, c.Value)
}

func (c *Condition) Evaluate(ctx EvalContext) (bool, error) {
	actualValue, ok := ctx.metrics[c.Field]
	if !ok {
		return false, fmt.Errorf("trigger field %s not exist metrics", c.Field)
	}

	switch c.Operator {
	case ">":
		return actualValue > c.Value, nil
	case "<":
		return actualValue < c.Value, nil
	case ">=":
		return actualValue >= c.Value, nil
	case "<=":
		return actualValue <= c.Value, nil
	case "=":
		return actualValue == c.Value, nil
	default:
		return false, fmt.Errorf("invalid condition operator: %s", c.Operator)
	}
}

type ParenExpr struct {
	Ep Expr
}

type TriggerParser struct {
	// all lex tokens
	tokens []Token
	// current token index
	pos int
}

func newTriggerParser(tokens []Token) *TriggerParser {
	return &TriggerParser{tokens: tokens}
}

func (t *TriggerParser) parse() (Expr, error) {
	return t.parseExpression()
}

func (t *TriggerParser) Evaluate(metrics map[string]float64) (bool, error) {
	var stack []string
	for _, token := range t.tokens {
		if token.Tp == TokenLParen {
			stack = append(stack, "(")
		} else if token.Tp == TokenRParen {
			if len(stack) != 0 && stack[len(stack)-1] == "(" {
				stack = stack[:len(stack)-1]
			} else {
				return false, nil
			}
		}
	}

	return true, nil
}

func (t *TriggerParser) parseExpression() (Expr, error) {
	left, err := t.parseTerm()
	if err != nil {
		return nil, err
	}

	for {
		token := t.peek()
		switch token.Tp {
		case TokenLogicalOp:
			t.consume()
			operator := token.Value
			right, er := t.parseTerm()
			if er != nil {
				return nil, er
			}
			return &LogicalExpr{
				Operator: operator,
				Left:     left,
				Right:    right,
			}, nil

		default:
			return left, nil
		}
	}
}

func (t *TriggerParser) parseTerm() (Expr, error) {
	token := t.peek()
	if token.Tp == TokenLParen {
		t.consume()
		expr, err := t.parseExpression()
		if err != nil {
			return nil, err
		}

		if t.peek().Tp != TokenRParen {
			return nil, fmt.Errorf("expected ')' but got '%s'", t.peek().Tp.String())
		}
		t.consume()

		return expr, nil
	}

	return t.parseCondition()
}

// parseCondition parse the condition unit.
func (t *TriggerParser) parseCondition() (Expr, error) {
	// get and validate field.
	filedToken := t.peek()
	if filedToken.Tp != TokenIdentifier {
		return nil, fmt.Errorf("expected identifier, got %v", filedToken.Tp)
	}
	field := filedToken.Value
	_, ok := metricsMap[field]
	if !ok {
		return nil, fmt.Errorf("expected metrics field, got %v", field)
	}

	// get and validate operator.
	t.consume()
	operatorToken := t.peek()
	if operatorToken.Tp != TokenOperator {
		return nil, fmt.Errorf("expected operator, got %v", operatorToken.Tp)
	}
	operator := operatorToken.Value
	_, ok = operatorsMap[operator]
	if !ok {
		return nil, fmt.Errorf("expected operator, got %v", operator)
	}

	// get and validate value.
	t.consume()
	valueToken := t.peek()
	if valueToken.Tp != TokenNumber {
		return nil, fmt.Errorf("expected number, got %v", valueToken.Tp)
	}
	value, err := strconv.ParseFloat(valueToken.Value, 64)
	if err != nil {
		return nil, err
	}

	t.consume()

	return &Condition{
		Field:    field,
		Operator: operator,
		Value:    value,
	}, nil
}

// peek return the next token without consuming it.
func (t *TriggerParser) peek() Token {
	if t.pos >= len(t.tokens) {
		return Token{Tp: -1, Value: ""}
	}

	return t.tokens[t.pos]
}

// consume moves to next token
func (t *TriggerParser) consume() {
	t.pos++
}

// parseTrigger the main method for parsing trigger and generate Expr.
func parseTrigger(trigger string) (Expr, error) {
	tokens, err := lex(trigger)
	if err != nil {
		return nil, err
	}

	return newTriggerParser(tokens).parse()
}
