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
)

// Limiter signal machine Limiter Unified Interface
type Limiter interface {
	// Allow To determine whether to allow the request to be processed
	Allow(ctx context.Context) (bool, error)
	// Close send signal to close the limiter
	Close()
}

type DisLimiter interface {
	// Allow To determine whether to allow the request to be processed
	// ctx context.Context
	// key is required if latitude is IP or User.
	Allow(ctx context.Context, key ...string) (bool, error)
	// Close send signal to close the limiter
	Close()
}
