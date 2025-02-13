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
)

// DecisionStrategy The decision-making strategy interface decides whether to dynamically
// adjust the request rate limit based on the real-time incoming indicator data.
type DecisionStrategy interface {
	// AdjustRate Calculate and decide whether to adjust the request rate.
	AdjustRate(ctx context.Context, metrics Metrics) Value
}

type Value struct {
	// Whether to adjust,if so,it returns true, otherwise it returns false.
	Adjust bool
	// if Adjust is true, it returns rate number, normal is zero.
	Rate float64
	// if decision is fail, it returns error.
	Err error
}

type BS struct {
	conf Config
}

func NewBS(p Parser) (DecisionStrategy, error) {
	cf, err := p.Parse()
	if err != nil {
		return nil, err
	}
	return &BS{
		conf: cf,
	}, nil
}

func (b *BS) AdjustRate(ctx context.Context, metrics Metrics) Value {
	select {
	case <-ctx.Done():
		return Value{
			Err: ctx.Err(),
		}
	default:
	}

	return Value{}
}

func (b *BS) checker(metrics Metrics) (alarm bool, err error) {

	return false, nil
}
