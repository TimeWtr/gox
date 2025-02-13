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
	"sync"
	"time"

	limiter2 "github.com/TimeWtr/gox/limiter"

	"github.com/TimeWtr/gox/log"
	"go.uber.org/zap"

	"github.com/TimeWtr/gox/errorx"
)

type EI interface {
	// Register the method to register latitude request rate.
	Register(ctx context.Context, latitude string, rate uint64, capacity int) error
	// Unregister the method to unregister latitude request rate.
	Unregister(ctx context.Context, latitude string) error
	// Notify the method to get the specified channel of sending metrics.
	Notify(ctx context.Context, latitude string) (chan<- Metrics, error)
	// DynamicController the method to dynamic adjust request
	// rate according to received metrics.
	DynamicController(interval time.Duration) error
	// Close the method to close distributed Executor.
	Close() error
}

type Options func(*Executor)

func WithLimiter(l limiter2.DisLimiter) Options {
	return func(e *Executor) {
		e.limiter = l
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
	limiter limiter2.DisLimiter
	// the strategy for deciding whether to modify rate.
	stg DecisionStrategy
	// logger
	lg log.Logger
	// close channel
	closeCh chan struct{}
}

func NewExecutor(cf Configuration, stg DecisionStrategy, opts ...Options) EI {
	logger, _ := zap.NewDevelopment()

	e := &Executor{
		ch:      map[string]chan Metrics{},
		mu:      new(sync.RWMutex),
		cf:      cf,
		stg:     stg,
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

	errCh := make(chan error, len(e.ch))

	for range ticker.C {
		select {
		case <-e.closeCh:
			e.lg.Infof("receive closed signal")
			break
		default:
			for latitude, ch := range e.ch {
				select {
				case metrics := <-ch:
					go func() {
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

							errCh <- res.Err
							return
						}

						if !res.Adjust {
							return
						}
						// modify
						e.lg.Infof("judge request rate adjusted", log.Field{
							Key:   "latitude",
							Value: latitude,
						}, log.Field{
							Key:   "rate",
							Value: res.Rate,
						})
					}()
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
