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

	"github.com/TimeWtr/gox/log"
	"go.uber.org/zap"

	"github.com/TimeWtr/gox/errorx"
)

type EI interface {
	Register(ctx context.Context, latitude string, rate uint64, capacity int) error
	Unregister(ctx context.Context, latitude string) error
	Notify(ctx context.Context, latitude string) (chan<- Metrics, error)
	DynamicController() error
	Close() error
}

type Options func(*Executor)

func WithDecisionStrategy(stg DecisionStrategy) Options {
	return func(e *Executor) {
		e.stg = stg
	}
}

func WithLogger(lg log.Logger) Options {
	return func(e *Executor) {
		e.lg = lg
	}
}

type Executor struct {
	// the channel collection for reporting metrics data.
	ch map[string]chan Metrics
	// locker lock the cf if request rate need to modify.
	mu *sync.RWMutex
	// the interface to operate request config rate.
	cf Configuration
	// limiter interface
	limiter DisLimiter
	// the strategy for deciding whether to modify rate.
	stg DecisionStrategy
	// logger
	lg log.Logger
	// close channel
	closeCh chan struct{}
}

func NewExecutor(cf Configuration, limiter DisLimiter, opts ...Options) *Executor {
	logger, _ := zap.NewDevelopment()

	e := &Executor{
		ch:      map[string]chan Metrics{},
		mu:      new(sync.RWMutex),
		cf:      cf,
		limiter: limiter,
		stg:     NewBS(),
		lg:      log.NewZapLogger(logger),
		closeCh: make(chan struct{}),
	}

	for _, opt := range opts {
		opt(e)
	}

	return e
}

// Register register latitude and request rate
func (e *Executor) Register(ctx context.Context, latitude string, rate uint64, capacity int) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.ch[latitude] = make(chan Metrics, capacity)
	return e.cf.Set(ctx, latitude, rate)
}

// Unregister unregister latitude and request rate.
func (e *Executor) Unregister(ctx context.Context, latitude string) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	delete(e.ch, latitude)
	return e.cf.Del(ctx, latitude)
}

// Notify the function to get the specified channel reported metrics.
func (e *Executor) Notify(ctx context.Context, latitude string) (chan<- Metrics, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	ch, ok := e.ch[latitude]
	if !ok {
		return nil, errorx.ErrMetricsChannelNotExists
	}

	return ch, nil
}

func (e *Executor) DynamicController(interval time.Duration) error {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		select {
		case <-e.closeCh:
			e.lg.Infof("receive closed signal")
			break
		default:
			for latitude, ch := range e.ch {
				select {
				case metrics := <-ch:
					ctx1, cancel := context.WithTimeout(context.Background(), 2*time.Second)
					res := e.stg.AdjustRate(ctx1, metrics)
					cancel()
					if res.Err != nil {
						// log error message
						e.lg.Errorf("judge request rate error", log.Field{
							Key:   "latitude",
							Value: latitude,
						}, log.Field{
							Key:   "error",
							Value: res.Err.Error(),
						})
						return res.Err
					}

					if !res.Adjust {
						continue
					}
					// modify
					e.lg.Infof("judge request rate adjusted", log.Field{
						Key:   "latitude",
						Value: latitude,
					}, log.Field{
						Key:   "rate",
						Value: res.Rate,
					})
				default:
				}
			}
		}
	}

	return nil
}

func (e *Executor) Close() error {
	close(e.closeCh)
	return nil
}

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
