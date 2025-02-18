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

package distributed

import (
	"context"
	"testing"

	"github.com/TimeWtr/gox/limiter/distributed/engine"

	"github.com/stretchr/testify/assert"
)

func TestNewBS(t *testing.T) {
	fs := engine.NewFileSource("./engine/examples/rule-new.json", engine.DataTypeYaml)
	p, err := engine.NewParser(fs)
	assert.Nil(t, err)
	bs, err := NewBS(p)
	assert.Nil(t, err)
	bs.AdjustRate(context.Background(), engine.Metrics{})
}
