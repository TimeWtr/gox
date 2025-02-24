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

package errorx

import "errors"

// ErrOverMaxRetries Retry strategy error
var (
	ErrOverMaxRetries = errors.New("over max retry limit")
)

// ErrOverMaxLimit Over limit
var (
	ErrOverMaxLimit = errors.New("over max limit")
	ErrClosed       = errors.New("limiter closed")
)

var (
	ErrMetricsChannelNotExists = errors.New("metrics channel not exists")
	ErrDelConfig               = errors.New("delete rate config error")
	ErrFileType                = errors.New("unsupported file type")
)

var (
	ErrEmptyNode = errors.New("empty node")
)
