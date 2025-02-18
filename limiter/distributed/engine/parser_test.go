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
	"testing"

	"github.com/stretchr/testify/assert"
)

//func TestNewYamlParser(t *testing.T) {
//	fs := NewFileSource("./examples/rule.yaml", DataTypeYaml)
//	parser, err := NewParser(fs)
//	assert.NoError(t, err)
//	cfg, err := parser.Parse()
//	assert.NoError(t, err)
//	t.Logf("config: %+v", cfg)
//}

func TestNewJsonParser(t *testing.T) {
	fs := NewFileSource("./examples/rule-new.json", DataTypeJson)
	parser, err := NewParser(fs)
	assert.NoError(t, err)
	cfg, err := parser.Parse()
	assert.NoError(t, err)
	t.Logf("config: %+v", cfg)
}

//func TestNewTomlParser(t *testing.T) {
//	fs := NewFileSource("./examples/rule.toml", DataTypeToml)
//	parser, err := NewParser(fs)
//	assert.NoError(t, err)
//	cfg, err := parser.Parse()
//	assert.NoError(t, err)
//	t.Logf("config: %+v", cfg)
//}
