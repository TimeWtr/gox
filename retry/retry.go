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

package retry

import (
	"math/rand"
	"time"

	"github.com/TimeWtr/gox/errorx"
)

// Strategy 重试策略接口
type Strategy interface {
	// Next 是否允许下次操作，允许则返回时间间隔，反之则返回error
	Next() (time.Duration, error)
}

// AgainNow 立即重试，非退避重试策略
type AgainNow struct {
	// 最大重试次数
	maxRetries int
	// 当前重试次数
	counter int
}

func NewAgainNow(maxRetries int) Strategy {
	return &AgainNow{
		maxRetries: maxRetries,
	}
}

func (a *AgainNow) Next() (time.Duration, error) {
	if a.counter >= a.maxRetries {
		return 0, errorx.ErrOverMaxRetries
	}

	a.counter++
	return 0, nil
}

// FixedInterval 固定时间间隔的退避重试策略
type FixedInterval struct {
	// 时间间隔
	interval time.Duration
	// 最大重试次数
	maxRetries int
	// 当前重试次数
	counter int
}

func NewFixedInterval(interval time.Duration, maxRetries int) Strategy {
	return &FixedInterval{
		interval:   interval,
		maxRetries: maxRetries,
	}
}

func (f *FixedInterval) Next() (time.Duration, error) {
	if f.counter >= f.maxRetries {
		return 0, errorx.ErrOverMaxRetries
	}

	f.counter++

	return f.interval, nil
}

// ExponentialBackoff 指数退避的重试策略，每次重试时间间隔x2
type ExponentialBackoff struct {
	// 当前重试时间间隔
	interval time.Duration
	// 最大重试次数
	maxRetries int
	// 当前重试次数
	counter int
}

func NewExponentialBackoff(interval time.Duration, maxRetries int) Strategy {
	return &ExponentialBackoff{
		interval:   interval,
		maxRetries: maxRetries,
	}
}

func (e *ExponentialBackoff) Next() (time.Duration, error) {
	if e.counter >= e.maxRetries {
		return 0, errorx.ErrOverMaxRetries
	}

	e.counter++
	e.interval <<= 1

	return e.interval, nil
}

// RandomInterval 随机重试时间间隔，单位为毫秒
type RandomInterval struct {
	// 最大重试次数
	maxRetries int
	// 当前重试次数
	counter int
	// 随机数的最大值，用于确定范围
	maxN int
}

func NewRandomInterval(maxRetries, maxN int) Strategy {
	return &RandomInterval{
		maxRetries: maxRetries,
		maxN:       maxN,
	}
}

func (r *RandomInterval) Next() (time.Duration, error) {
	if r.counter >= r.maxRetries {
		return 0, errorx.ErrOverMaxRetries
	}

	r.counter++

	rand.Seed(time.Now().UnixNano())
	randomNumber := rand.Intn(r.maxN)

	return time.Duration(randomNumber) * time.Millisecond, nil
}

// ComplexInterval 综合重试策略，结合指数和随机策略，在指数的基础上加上随机的时间间隔
type ComplexInterval struct {
	// 当前时间间隔
	interval time.Duration
	// 最大重试次数
	maxRetries int
	// 当前重试次数
	counter int
	// 随机数的最大值，用于确定范围
	maxN int
}

func NewComplexInterval(interval time.Duration, maxRetries, maxN int) Strategy {
	return &ComplexInterval{
		interval:   interval,
		maxRetries: maxRetries,
		maxN:       maxN,
	}
}

func (c *ComplexInterval) Next() (time.Duration, error) {
	if c.counter >= c.maxRetries {
		return 0, errorx.ErrOverMaxRetries
	}

	c.counter++
	rand.Seed(time.Now().UnixNano())
	randomNumber := rand.Intn(c.maxN)
	c.interval <<= 1

	return time.Duration(randomNumber)*time.Millisecond + c.interval, nil
}
