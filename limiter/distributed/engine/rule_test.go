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

func TestBuildRuleTrees(t *testing.T) {
	testCases := []struct {
		name    string
		ctx     func() EvalContext
		wantRes bool
		wantErr error
	}{
		{
			name: "over threshold and or expression",
			ctx: func() EvalContext {
				return WithEvalContext(map[string]float64{
					"cpu_usage": 0.9,
					"mem_usage": 0.8,
				})
			},
			wantRes: true,
			wantErr: nil,
		},
		{
			name: "threshold and or expression",
			ctx: func() EvalContext {
				return WithEvalContext(map[string]float64{
					"cpu_usage": 0.8,
					"mem_usage": 0.8,
				})
			},
			wantRes: false,
			wantErr: nil,
		},
		{
			name: "not exist metrics",
			ctx: func() EvalContext {
				return WithEvalContext(map[string]float64{
					"xxx": 0.3,
				})
			},
			wantRes: false,
			wantErr: fmt.Errorf("trigger field cpu_usage not exist metrics"),
		},
	}

	fs := NewFileSource("./examples/rule.json", DataTypeJson)
	parser, err := NewParser(fs)
	assert.NoError(t, err)
	cfg, err := parser.Parse()
	assert.NoError(t, err)
	t.Logf("config: %+v", cfg)
	trees, err := BuildRuleTrees(cfg.Rules)
	assert.NoError(t, err)
	t.Logf("trees: %+v", trees)
	for _, rt := range trees[0].GetChildren() {
		t.Logf("rt:%+v\n", rt)
		if rt.triggerAST != nil {
			t.Logf("ast string: %s\n", rt.triggerAST.String())
			for _, tc := range testCases {
				t.Run(tc.name, func(t *testing.T) {
					ok, err := rt.triggerAST.Evaluate(tc.ctx())
					assert.Equal(t, tc.wantErr, err)
					assert.Equal(t, tc.wantRes, ok)
				})
			}
		}
	}
}
