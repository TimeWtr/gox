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

package limiter

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

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

func NewDSlidingWindow() Limiter {
	return &DSlidingWindow{}
}

// Register register latitude and request rate
func (d *DSlidingWindow) Register(ctx context.Context, latitude string, rate int) error {
	_, err := d.client.HSet(ctx, latitude, rate).Result()
	return err
}

// UnRegister unregister latitude and request rate
func (d *DSlidingWindow) UnRegister(ctx context.Context, latitude string) error {
	_, err := d.client.HDel(ctx, latitude).Result()
	return err
}

func (d *DSlidingWindow) Allow(ctx context.Context) (bool, error) {

	return true, nil
}

func (d *DSlidingWindow) Close() {}
