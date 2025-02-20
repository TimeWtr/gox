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
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConsistentHash(t *testing.T) {
	ch := NewConsistentHash(nil, 3, 6)
	err := ch.AddNode([]Node{
		{
			Val:    "Node1",
			Weight: LowWeight,
		},
		{
			Val:    "Node2",
			Weight: MidWeight,
		},
		{
			Val:    "Node3",
			Weight: HighWeight,
		},
	}...)
	assert.Nil(t, err)

	var keys []string
	for i := 0; i < 1000000; i++ {
		keys = append(keys, "key"+strconv.Itoa(i))
	}

	for _, key := range keys {
		node, er := ch.GetNode([]byte(key))
		assert.Nil(t, er)
		t.Logf("key: %s -> node: %s\n", key, node)
	}

	err = ch.RemoveNode(Node{
		Val:    "Node1",
		Weight: LowWeight,
	})
	assert.Nil(t, err)
}

func TestNewConsistentHash_NoLog(t *testing.T) {
	ch := NewConsistentHash(nil, 3, 6)
	err := ch.AddNode([]Node{
		{
			Val:    "Node1",
			Weight: LowWeight,
		},
		{
			Val:    "Node2",
			Weight: MidWeight,
		},
		{
			Val:    "Node3",
			Weight: HighWeight,
		},
	}...)
	assert.Nil(t, err)

	var keys []string
	for i := 0; i < 1000000; i++ {
		keys = append(keys, "key"+strconv.Itoa(i))
	}

	for _, key := range keys {
		_, er := ch.GetNode([]byte(key))
		assert.Nil(t, er)
	}

	err = ch.RemoveNode(Node{
		Val:    "Node1",
		Weight: LowWeight,
	})
	assert.Nil(t, err)
}

func BenchmarkNewConsistentHash(b *testing.B) {
	ch := NewConsistentHash(nil, 3, 6)
	err := ch.AddNode([]Node{
		{
			Val:    "Node1",
			Weight: LowWeight,
		},
		{
			Val:    "Node2",
			Weight: MidWeight,
		},
		{
			Val:    "Node3",
			Weight: HighWeight,
		},
	}...)
	assert.Nil(b, err)

	for i := 0; i < b.N; i++ {
		key := "key" + strconv.Itoa(i)
		node, er := ch.GetNode([]byte(key))
		assert.Nil(b, er)
		b.Logf("key: %s -> node: %s\n", key, node)
	}

	err = ch.RemoveNode(Node{
		Val:    "Node1",
		Weight: LowWeight,
	})
	assert.Nil(b, err)
}

func BenchmarkNewConsistentHash_NoLog(b *testing.B) {
	ch := NewConsistentHash(nil, 3, 6)
	err := ch.AddNode([]Node{
		{
			Val:    "Node1",
			Weight: LowWeight,
		},
		{
			Val:    "Node2",
			Weight: MidWeight,
		},
		{
			Val:    "Node3",
			Weight: HighWeight,
		},
	}...)
	assert.Nil(b, err)

	for i := 0; i < b.N; i++ {
		key := "key" + strconv.Itoa(i)
		_, er := ch.GetNode([]byte(key))
		assert.Nil(b, er)
	}

	err = ch.RemoveNode(Node{
		Val:    "Node1",
		Weight: LowWeight,
	})
	assert.Nil(b, err)
}
