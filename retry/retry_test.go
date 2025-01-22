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
	"testing"
	"time"

	"github.com/TimeWtr/gox/errorx"
	"github.com/stretchr/testify/assert"
)

func TestNewAgainNow(t *testing.T) {
	var timer *time.Timer
	strategy := NewAgainNow(5)
	for i := 0; i < 6; i++ {
		duration, err := strategy.Next()
		if err != nil {
			assert.Error(t, err, errorx.ErrOverMaxRetries)
			return
		}

		t.Log("允许重试，时间间隔：", duration)
		if timer == nil {
			timer = time.NewTimer(duration)
		} else {
			timer.Reset(duration)
		}
		<-timer.C
	}
}

func TestNewFixedInterval(t *testing.T) {
	var timer *time.Timer
	strategy := NewFixedInterval(time.Millisecond*50, 5)
	for i := 0; i < 6; i++ {
		duration, err := strategy.Next()
		if err != nil {
			assert.Error(t, err, errorx.ErrOverMaxRetries)
			return
		}

		t.Log("允许重试，时间间隔：", duration)
		assert.Equal(t, time.Millisecond*50, duration)
		if timer == nil {
			timer = time.NewTimer(duration)
		} else {
			timer.Reset(duration)
		}
		<-timer.C
	}
}

func TestNewExponentialBackoff(t *testing.T) {
	var timer *time.Timer
	strategy := NewExponentialBackoff(time.Millisecond*50, 5)
	targetDuration := time.Millisecond * 50
	for i := 0; i < 6; i++ {
		duration, err := strategy.Next()
		if err != nil {
			assert.Error(t, err, errorx.ErrOverMaxRetries)
			return
		}

		t.Log("允许重试，时间间隔：", duration)
		targetDuration <<= 1
		assert.Equal(t, targetDuration, duration)
		if timer == nil {
			timer = time.NewTimer(duration)
		} else {
			timer.Reset(duration)
		}
		<-timer.C
	}
}

func TestNewRandomInterval(t *testing.T) {
	var timer *time.Timer
	strategy := NewRandomInterval(5, 100)
	for i := 0; i < 6; i++ {
		duration, err := strategy.Next()
		if err != nil {
			assert.Error(t, err, errorx.ErrOverMaxRetries)
			return
		}

		t.Log("允许重试，时间间隔：", duration)
		if timer == nil {
			timer = time.NewTimer(duration)
		} else {
			timer.Reset(duration)
		}
		<-timer.C
	}
}

func TestNewComplexInterval(t *testing.T) {
	var timer *time.Timer
	strategy := NewComplexInterval(time.Millisecond*10, 5, 100)
	for i := 0; i < 6; i++ {
		duration, err := strategy.Next()
		if err != nil {
			assert.Error(t, err, errorx.ErrOverMaxRetries)
			return
		}

		t.Log("允许重试，时间间隔：", duration)
		if timer == nil {
			timer = time.NewTimer(duration)
		} else {
			timer.Reset(duration)
		}
		<-timer.C
	}
}
