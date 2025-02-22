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
	"sync"
	"time"
)

// CircuitState the status for limiter.
// if not limit, the status is StatusClosed.
// if limited, the status is StatusOpen.
// if limiter is recovering, the status is StatusRecover.
type CircuitState int

const (
	StatusNormal     CircuitState = iota // normal status, default status
	StatusThrottling                     // limit status
	StatusRecovering                     // status recover
)

func (s CircuitState) String() string {
	switch s {
	case StatusNormal:
		return "normal status"
	case StatusThrottling:
		return "throttling status"
	case StatusRecovering:
		return "recovering status"
	default:
		return "unknown state"
	}
}

var once *sync.Once

type LimitStatus struct {
	// limit current status
	state CircuitState
	// limit time
	throttleSince time.Time
	// recover steps, item is time duration
	recoverSteps []int
	// current step
	currentStep int
	// When the grayscale is restored, if the threshold exceeds the threshold,
	// is it allowed to go back to the previous step
	rollback bool
	// locker
	mu *sync.RWMutex
}

func NewLimitStatus(rollback bool, recoverSteps []int) *LimitStatus {
	ls := &LimitStatus{}
	once.Do(func() {
		ls = &LimitStatus{
			state:        StatusNormal,
			recoverSteps: recoverSteps,
			rollback:     rollback,
			mu:           new(sync.RWMutex),
		}
	})

	return ls
}
