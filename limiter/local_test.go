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
	"testing"
	"time"

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
