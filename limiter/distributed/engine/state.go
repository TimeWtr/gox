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

// CircuitState the status for limiter.
// if not limit, the status is StatusClosed.
// if limited, the status is StatusOpen.
// if limiter is recovering, the status is StatusRecover.
type CircuitState int

const (
	StatusClosed  CircuitState = iota // normal status
	StatusOpen                        // limit status
	StatusRecover                     // status recover
)

func (s CircuitState) String() string {
	switch s {
	case StatusClosed:
		return "closed status"
	case StatusOpen:
		return "open status"
	case StatusRecover:
		return "recover status"
	default:
		return "unknown state"
	}
}
