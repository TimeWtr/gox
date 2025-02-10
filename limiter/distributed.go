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
	"time"

	"github.com/TimeWtr/gox/errorx"

	"github.com/redis/go-redis/v9"
)

const (
	DefaultHashName = "meta"
	DefaultInterval = time.Second * 5
)

const (
	DefaultCPUKey            = "CPU"
	DefaultMemoryKey         = "Memory"
	DefaultMemoryUsageKey    = "Memory Usage"
	DefaultMemoryUsedKey     = "Memory Used"
	DefaultRequestLatencyKey = "Request Latency"
	DefaultErrRateKey        = "Error Rate"
	DefaultActiveConnsKey    = "Active Conns"
)

type Options func(*DSlidingWindow)

// WithHashName set the name for metrics metadata hash table.
func WithHashName(hashTableName string) Options {
	return func(d *DSlidingWindow) {
		d.HashTableName = hashTableName
	}
}

// WithInterval set the window size
func WithInterval(interval time.Duration) Options {
	return func(d *DSlidingWindow) {
		d.interval = interval
	}
}

var _ DisLimiter = (*DSlidingWindow)(nil)

// DSlidingWindow distributed sliding window implement based on redis.
// this implement supports dynamic adjustment of the limit threshold
// and reception of the collected machine metrics. sliding window
// does not provide collection metrics functions.
type DSlidingWindow struct {
	// redis client
	client redis.Cmdable
	// window size
	interval time.Duration
	// the channel collection for reporting metrics data.
	ch map[string]chan Metrics
	// the name for metadata hash table in redis.
	HashTableName string
	// locker
	mu *sync.RWMutex
}

func NewDSlidingWindow(client redis.Cmdable, opts ...Options) DisLimiter {
	limiter := &DSlidingWindow{
		client: client,
		ch:     map[string]chan Metrics{},
		mu:     &sync.RWMutex{},
	}

	for _, opt := range opts {
		opt(limiter)
	}

	return limiter
}

// Register register latitude and request rate
func (d *DSlidingWindow) Register(ctx context.Context, latitude string, rate int, capacity int) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.ch[latitude] = make(chan Metrics, capacity)
	_, err := d.client.HSet(ctx, d.HashTableName, latitude, rate).Result()
	return err
}

// UnRegister unregister latitude and request rate
func (d *DSlidingWindow) UnRegister(ctx context.Context, latitude string) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	delete(d.ch, latitude)
	_, err := d.client.HDel(ctx, d.HashTableName, latitude).Result()
	return err
}

func (d *DSlidingWindow) Notify(ctx context.Context, latitude string) (chan<- Metrics, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	ch, ok := d.ch[latitude]
	if !ok {
		return nil, errorx.MetricsChannelNotExists
	}

	return ch, nil
}

func (d *DSlidingWindow) DynamicRate(ctx context.Context, latitude string, strategy DecisionStrategy) error {

	return nil
}

func (d *DSlidingWindow) Allow(ctx context.Context, latitude string) (bool, error) {

	return true, nil
}

func (d *DSlidingWindow) Close() {}

type Metrics struct {
	// used cpu percent,
	CPUUsage float64 `json:"cpu_usage,omitempty"`
	// memory percent
	MemUsage float64 `json:"mem_usage,omitempty"`
	// used memory size, unit is bytes.
	MemoryUsed uint64 `json:"memory_used,omitempty"`
	// request latency, only used while type is API.
	RequestLatency float64 `json:"request_latency,omitempty"`
	// request error rate.
	ErrRate float64 `json:"err_rate,omitempty"`
	// current active connections.
	ActiveConns uint64 `json:"active_conns,omitempty"`
}

type MetricsV1 struct {
	Timestamp time.Time      `json:"timestamp"`
	Data      map[string]any `json:"data"`
}
