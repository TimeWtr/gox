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

package engine

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParseTime(t *testing.T) {
	testCases := []struct {
		name    string
		t       string
		wantErr error
		wantRes time.Duration
	}{
		{
			name:    "10s",
			t:       "10s",
			wantErr: nil,
			wantRes: 10 * time.Second,
		},
		{
			name:    "1m",
			t:       "1m",
			wantErr: nil,
			wantRes: 60 * time.Second,
		},
		{
			name:    "1h",
			t:       "1h",
			wantErr: nil,
			wantRes: 60 * time.Minute,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			d, err := parseTime(tc.t)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantRes, d)
			t.Log("duration:", d)
		})
	}
}
