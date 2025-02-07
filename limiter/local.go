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
	"sync"
	"sync/atomic"
	"time"

	"github.com/TimeWtr/gox/errorx"
)

// local node limiter

//var (
//	once sync.Once
//	bk   *Buckets
//)

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
	bk := &Buckets{
		ch:       make(chan struct{}, capacity),
		closeCh:  make(chan struct{}),
		interval: interval,
	}

	ticker := time.NewTicker(bk.interval)

	go func() {
		for {
			select {
			case <-bk.closeCh:
				close(bk.ch)
				ticker.Stop()
				return
			case <-ticker.C:
				select {
				case bk.ch <- struct{}{}:
				default:
				}
			}
		}
	}()

	return bk
}

func (b *Buckets) Allow(ctx context.Context) (bool, error) {
	select {
	case <-ctx.Done():
		// context timeout
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

// LeakyBucket The leaky bucket algorithm is implemented by ticker.
type LeakyBucket struct {
	// time duration
	ticker *time.Ticker
	// once do
	once sync.Once
}

func NewLeakyBucket(interval time.Duration) Limiter {
	return &LeakyBucket{
		ticker: time.NewTicker(interval),
		once:   sync.Once{},
	}
}

func (l *LeakyBucket) Allow(ctx context.Context) (bool, error) {
	select {
	case <-ctx.Done():
		return false, ctx.Err()
	case <-l.ticker.C:
		return true, nil
	default:
		return false, errorx.ErrOverMaxLimit
	}
}

func (l *LeakyBucket) Close() {
	l.once.Do(func() {
		l.ticker.Stop()
	})
}

// FixedWindow The fixed window algorithm is implemented by fix window(interval).
type FixedWindow struct {
	// window size
	interval time.Duration
	// start time
	startTime int64
	// request limit rate
	rate int64
	// current window request counter
	cnt int64
}

func NewFixedWindow(interval time.Duration, rate int64) Limiter {
	return &FixedWindow{
		interval:  interval,
		startTime: time.Now().UnixNano(),
		rate:      rate,
	}
}

func (f *FixedWindow) Allow(ctx context.Context) (bool, error) {
	select {
	case <-ctx.Done():
		return false, ctx.Err()
	default:
	}

	// Determine whether the current timestamp is within the window period.
	now := time.Now().UnixNano()
	cnt := atomic.LoadInt64(&f.cnt)
	if f.startTime+f.interval.Nanoseconds() <= now {
		// window expired
		if atomic.CompareAndSwapInt64(&f.startTime, f.startTime, now) {
			atomic.CompareAndSwapInt64(&f.cnt, cnt, 0)
		}
	}

	select {
	case <-ctx.Done():
		return false, ctx.Err()
	default:
	}

	cnt = atomic.AddInt64(&f.cnt, 1)
	if cnt >= f.rate {
		// over request limit
		return false, errorx.ErrOverMaxLimit
	}

	// counter ++
	atomic.AddInt64(&f.cnt, 1)

	return true, nil
}

func (f *FixedWindow) Close() {}

type LeakyWindow struct {
}

func (l *LeakyWindow) Allow(ctx context.Context) (bool, error) {
	//TODO implement me
	return true, nil
}

func (l *LeakyWindow) Close() {
	//TODO implement me
	return
}
