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

package liniter

import (
	"context"
	"sync"
	"time"

	"github.com/TimeWtr/gox/errorx"
)

// local node limiter

var (
	once sync.Once
	bk   *Buckets
)

// Buckets token bucket limiter
type Buckets struct {
	// token bucket
	ch chan struct{}
	// close signal channel
	closeCh chan struct{}
	// grant token interval
	interval time.Duration
}

func NewBuckets(interval time.Duration, capacity int64) Limiter {
	once.Do(func() {
		bk = &Buckets{
			ch:       make(chan struct{}, capacity),
			closeCh:  make(chan struct{}),
			interval: interval,
		}

		go func() {
			ticker := time.NewTicker(bk.interval)
			defer ticker.Stop()

			for {
				select {
				case <-bk.closeCh:
					close(bk.ch)
					return
				case <-ticker.C:
					select {
					case bk.ch <- struct{}{}:
					default:
					}
				}
			}
		}()
	})

	return bk
}

func (b *Buckets) Allow(ctx context.Context) (bool, error) {
	select {
	case <-ctx.Done():
		// context deadline
		return false, ctx.Err()
	case <-b.closeCh:
		// receive close signal
		return false, errorx.ErrClosed
	case <-b.ch:
		// get a token
		return true, nil
	default:
		return false, errorx.ErrOverMaxLimit
	}
}

func (b *Buckets) Close() {
	close(b.closeCh)
}
