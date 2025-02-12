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

package local

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/TimeWtr/gox/errorx"

	"github.com/stretchr/testify/assert"
)

func TestNewBuckets(t *testing.T) {
	buckets := NewBuckets(time.Millisecond*5, 10)

	var wg sync.WaitGroup
	wg.Add(2)
	closeCh := make(chan struct{})
	go func() {
		defer wg.Done()

		for {
			time.Sleep(5 * time.Millisecond)
			select {
			case <-closeCh:
				return
			default:
			}

			ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*500)
			ok, err := buckets.Allow(ctx)
			cancel()
			if err != nil {
				t.Log("error:", err)
				continue
			}

			t.Log("ok:", ok)
			if !ok {
				continue
			}
		}
	}()

	go func() {
		defer wg.Done()
		time.Sleep(time.Second * 10)
		buckets.Close()
		close(closeCh)
	}()

	wg.Wait()
}

func TestNewBuckets_limit(t *testing.T) {
	buckets := NewBuckets(time.Millisecond*100, 2)

	var wg sync.WaitGroup
	wg.Add(2)
	closeCh := make(chan struct{})
	go func() {
		defer wg.Done()

		for {
			time.Sleep(5 * time.Millisecond)
			select {
			case <-closeCh:
				return
			default:
			}

			ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*500)
			ok, err := buckets.Allow(ctx)
			cancel()
			if err != nil {
				t.Log("error:", err)
				continue
			}

			t.Log("ok:", ok)
			if !ok {
				continue
			}
		}
	}()

	go func() {
		defer wg.Done()
		time.Sleep(time.Second * 10)
		buckets.Close()
		close(closeCh)
	}()

	wg.Wait()
}

func TestLeakyBucket_Allow(t *testing.T) {
	lb := NewLeakyBucket(50 * time.Nanosecond)
	var wg sync.WaitGroup
	wg.Add(2)
	closeCh := make(chan struct{})
	go func() {
		defer wg.Done()
		for {
			select {
			case <-closeCh:
				t.Log("receive close signal")
				return
			default:
			}

			ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*500)
			ok, err := lb.Allow(ctx)
			cancel()
			if err != nil {
				t.Log("error:", err)
				continue
			}

			t.Log("ok:", ok)
			if !ok {
				continue
			}
		}
	}()

	go func() {
		defer wg.Done()
		time.Sleep(time.Second * 5)
		lb.Close()
		close(closeCh)
	}()

	wg.Wait()
}

func TestLeakyBucket_Context_Deadline_Exceeded(t *testing.T) {
	lb := NewLeakyBucket(time.Millisecond * 5)
	var wg sync.WaitGroup
	wg.Add(2)
	closeCh := make(chan struct{})
	go func() {
		defer wg.Done()
		for {
			select {
			case <-closeCh:
				return
			default:
			}

			ctx, cancel := context.WithTimeout(context.Background(), time.Nanosecond)
			_, err := lb.Allow(ctx)
			cancel()
			assert.Error(t, context.DeadlineExceeded, err)
		}
	}()

	go func() {
		defer wg.Done()
		time.Sleep(time.Second * 2)
		lb.Close()
		close(closeCh)

	}()

	wg.Wait()
}

func TestLeakyBucket_Err_Limit(t *testing.T) {
	lb := NewLeakyBucket(time.Second * 5)
	var wg sync.WaitGroup
	wg.Add(2)

	closeCh := make(chan struct{})
	go func() {
		defer wg.Done()
		for {
			select {
			case <-closeCh:
				return
			default:
			}

			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			_, err := lb.Allow(ctx)
			cancel()
			assert.Error(t, errorx.ErrOverMaxLimit, err)
		}
	}()

	go func() {
		defer wg.Done()
		time.Sleep(time.Second * 2)
		lb.Close()
		close(closeCh)

	}()

	wg.Wait()
}

func TestNewFixedWindow(t *testing.T) {
	fw := NewFixedWindow(time.Second*5, 1)
	var wg sync.WaitGroup
	wg.Add(2)
	closeCh := make(chan struct{})
	go func() {
		defer wg.Done()

		ticker := time.NewTicker(2 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-closeCh:
				return
			case <-ticker.C:
			}

			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			ok, err := fw.Allow(ctx)
			cancel()
			if err != nil {
				t.Log("error:", err)
			}

			t.Log("ok:", ok)
			if !ok {
				continue
			}
		}
	}()

	go func() {
		defer wg.Done()
		defer close(closeCh)
		defer fw.Close()

		time.Sleep(time.Second * 5)
	}()

	wg.Wait()
}

func TestSlidingWindow_Allow(t *testing.T) {
	lb := NewSlidingWindow(time.Second*5, 5)
	var wg sync.WaitGroup

	wg.Add(2)
	closeCh := make(chan struct{})
	go func() {
		defer wg.Done()

		ticker := time.NewTicker(500 * time.Millisecond)
		defer ticker.Stop()
		for {
			select {
			case <-closeCh:
				t.Log("receive close signal")
				return
			case <-ticker.C:
				ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
				ok, err := lb.Allow(ctx)
				cancel()
				if err != nil {
					t.Log("error:", err)
				} else {
					t.Log("ok:", ok)
				}
			default:
			}
		}
	}()

	go func() {
		defer wg.Done()

		time.Sleep(time.Second * 11)
		lb.Close()
		close(closeCh)
	}()

	wg.Wait()
}

func TestSlidingWindow_Context_Timeout(t *testing.T) {
	lb := NewSlidingWindow(time.Second*5, 1)
	ctx, cancel := context.WithTimeout(context.Background(), time.Nanosecond)
	defer cancel()
	ok, err := lb.Allow(ctx)
	assert.False(t, ok)
	assert.Error(t, context.DeadlineExceeded, err)
}

func BenchmarkSlidingWindow_Allow(b *testing.B) {
	lb := NewSlidingWindow(time.Second, 100000)
	for i := 0; i < b.N; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		ok, err := lb.Allow(ctx)
		cancel()
		if err != nil {
			b.Log(err)
			continue
		}
		b.Log("ok:", ok)
	}
}
