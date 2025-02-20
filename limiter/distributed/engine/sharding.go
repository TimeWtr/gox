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
	"hash/crc32"
	"sort"
	"strconv"
	"sync"

	"github.com/TimeWtr/gox/errorx"
)

const (
	LowWeight WeightType = iota
	MidWeight
	HighWeight
)

type WeightType int

func (wt *WeightType) String() string {
	switch *wt {
	case LowWeight:
		return "low weight"
	case MidWeight:
		return "middle weight"
	case HighWeight:
		return "high weight"
	default:
		return "unknown wight"
	}
}

func (wt *WeightType) valid() error {
	switch *wt {
	case LowWeight, MidWeight, HighWeight:
		return nil
	default:
		return fmt.Errorf("unknown WeightType: %d", *wt)
	}
}

type HashFunc func(data []byte) uint32

type ConsistentHash struct {
	// the hash handle function
	fn HashFunc
	// node replicas num
	replicas int
	// max node replicas num, some node has stronger ability.
	maxReplicas int
	// all node list include real node and virtual node.
	keys []int
	// the relationship between real node and virtual node.
	mp map[int]string
	// locker
	mu *sync.RWMutex
}

func NewConsistentHash(fn HashFunc, replicas, maxReplicas int) *ConsistentHash {
	// default hash function is crc32.ChecksumIEEE
	if fn == nil {
		fn = crc32.ChecksumIEEE
	}

	return &ConsistentHash{
		fn:          fn,
		replicas:    replicas,
		maxReplicas: maxReplicas,
		mp:          make(map[int]string),
		mu:          new(sync.RWMutex),
	}
}

// AddNode the method for adding nodes to hash ring.
func (c *ConsistentHash) AddNode(nodes ...Node) error {
	if len(nodes) == 0 {
		return errorx.ErrEmptyNode
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	for _, node := range nodes {
		err := node.Weight.valid()
		if err != nil {
			return err
		}

		// dynamic calculate replicas num according to the node weight.
		replicas := c.calculateReplicas(node.Weight)
		for i := 0; i < replicas; i++ {
			hash := int(c.fn([]byte(node.Val + strconv.Itoa(i))))
			c.mp[hash] = node.Val
			c.keys = append(c.keys, hash)
		}
	}

	// sort the slice.
	sort.Ints(c.keys)

	return nil
}

// RemoveNode the method for removing node from hash ring.
func (c *ConsistentHash) RemoveNode(nodes ...Node) error {
	if len(nodes) == 0 {
		return errorx.ErrEmptyNode
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	for _, node := range nodes {
		err := node.Weight.valid()
		if err != nil {
			return err
		}

		// dynamic calculate replicas num according to the node weight.
		replicas := c.calculateReplicas(node.Weight)
		for i := 0; i < replicas; i++ {
			hash := int(c.fn([]byte(node.Val + strconv.Itoa(i))))
			delete(c.mp, hash)
			//index := sort.SearchInts(c.keys, hash)
			//c.keys = append(c.keys[:index], c.keys[index+1:]...)
		}
	}

	c.keys = c.keys[:0]
	for hash := range c.mp {
		c.keys = append(c.keys, hash)
	}
	sort.Ints(c.keys)

	return nil
}

// GetNode the method for calculating hash value and return the real node.
func (c *ConsistentHash) GetNode(key []byte) (string, error) {
	if len(c.keys) == 0 {
		return "", errorx.ErrEmptyNode
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	// find the first node lager then the hash value.
	hash := int(c.fn(key))
	index := sort.Search(len(c.keys), func(i int) bool {
		return c.keys[i] >= hash
	})

	if index == len(c.keys) {
		index = 0
	}

	return c.mp[c.keys[index]], nil
}

func (c *ConsistentHash) calculateReplicas(weight WeightType) int {
	replicas := 0
	switch weight {
	case LowWeight:
		replicas = c.replicas
	case MidWeight:
		replicas = c.maxReplicas - c.replicas/2
	case HighWeight:
		replicas = c.maxReplicas
	}

	return replicas
}

type Node struct {
	// node value
	Val string
	// node weight
	Weight WeightType
}
