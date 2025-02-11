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
	"time"

	"github.com/TimeWtr/gox/limiter"

	"github.com/redis/go-redis/v9"
)

var _ limiter.DisLimiter = (*DSlidingWindow)(nil)

// DSlidingWindow distributed sliding window implement based on redis.
// this implement supports dynamic adjustment of the limit threshold
// and reception of the collected machine metrics. sliding window
// does not provide collection metrics functions.
type DSlidingWindow struct {
	// redis client
	client redis.Cmdable
	// window size
	interval time.Duration
}

func NewDSlidingWindow(client redis.Cmdable) limiter.DisLimiter {
	return &DSlidingWindow{
		client: client,
	}
}

func (d *DSlidingWindow) Allow(ctx context.Context, latitude string) (bool, error) {

	return true, nil
}

func (d *DSlidingWindow) Close() {}
