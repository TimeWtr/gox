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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLex(t *testing.T) {
	testCases := []struct {
		name    string
		input   string
		wantRes []Token
		wantErr error
	}{
		{
			name:    "err number, last is dot",
			input:   "test > 0.00.",
			wantRes: []Token{},
			wantErr: fmt.Errorf("invalid character in number: 0.00, position: 11"),
		},
		{
			name:    "err number, two dot",
			input:   "test > 0.0.0",
			wantRes: []Token{},
			wantErr: fmt.Errorf("invalid character in number: 0.0, position: 10"),
		},
		{
			name:  "letter,operator and number",
			input: "test > 0.01",
			wantRes: []Token{
				{
					Tp:    TokenIdentifier,
					Value: "test",
				},
				{
					Tp:    TokenOperator,
					Value: ">",
				},
				{
					Tp:    TokenNumber,
					Value: "0.01",
				},
			},
			wantErr: nil,
		},
		{
			name:  "letter,operator,and,number",
			input: "cpu_usage > 0.9 and mem_usage >= 0.8",
			wantRes: []Token{
				{
					Tp:    TokenIdentifier,
					Value: "cpu_usage",
				},
				{
					Tp:    TokenOperator,
					Value: ">",
				},
				{
					Tp:    TokenNumber,
					Value: "0.9",
				},
				{
					Tp:    TokenLogicalOp,
					Value: "AND",
				},
				{
					Tp:    TokenIdentifier,
					Value: "mem_usage",
				},
				{
					Tp:    TokenOperator,
					Value: ">=",
				},
				{
					Tp:    TokenNumber,
					Value: "0.8",
				},
			},
			wantErr: nil,
		},
		{
			name:  "letter,operator,and,number,Paren",
			input: "cpu_usage > 0.9 and (mem_usage >= 0.8 OR err_rate > 0.2)",
			wantRes: []Token{
				{
					Tp:    TokenIdentifier,
					Value: "cpu_usage",
				},
				{
					Tp:    TokenOperator,
					Value: ">",
				},
				{
					Tp:    TokenNumber,
					Value: "0.9",
				},
				{
					Tp:    TokenLogicalOp,
					Value: "AND",
				},
				{
					Tp:    TokenLParen,
					Value: "(",
				},
				{
					Tp:    TokenIdentifier,
					Value: "mem_usage",
				},
				{
					Tp:    TokenOperator,
					Value: ">=",
				},
				{
					Tp:    TokenNumber,
					Value: "0.8",
				},
				{
					Tp:    TokenLogicalOp,
					Value: "OR",
				},
				{
					Tp:    TokenIdentifier,
					Value: "err_rate",
				},
				{
					Tp:    TokenOperator,
					Value: ">",
				},
				{
					Tp:    TokenNumber,
					Value: "0.2",
				},
				{
					Tp:    TokenRParen,
					Value: ")",
				},
			},
			wantErr: nil,
		},
		{
			name:    "error identifier",
			input:   "** _test ",
			wantRes: []Token{},
			wantErr: fmt.Errorf("invalid character in identifier: *, position: 0"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tokens, err := lex(tc.input)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}

			assert.Equal(t, len(tc.wantRes), len(tokens))
			for i, token := range tokens {
				assert.Equal(t, tc.wantRes[i], token)
			}
		})
	}
}

func TestParseTrigger(t *testing.T) {
	expr, err := parseTrigger("cpu_usage > 0.9 and (mem_usage >= 0.8 OR err_rate > 0.2)")
	assert.NoError(t, err)
	//expr.String()
	if expr.GetType() == NodeLogical {
		children := expr.GetChildren()
		children[0].String()
		children[1].String()
	}
}
