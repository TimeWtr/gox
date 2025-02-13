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
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/TimeWtr/gox/errorx"
)

const (
	// ReduceAction decrease request rate.
	ReduceAction = "decrease"
)

type Rule struct {
	RuleName  string      `yaml:"ruleName" json:"ruleName" toml:"ruleName"`
	Threshold json.Number `yaml:"threshold" json:"threshold" toml:"threshold"`
	Action    string      `yaml:"action" json:"action" toml:"action"`
	Amount    int         `yaml:"amount" json:"amount" toml:"amount"`
}

type GrayRecover struct {
	GrayScale   float32 `yaml:"grayScale" json:"grayScale" toml:"grayScale"`
	RecoverTime int64   `yaml:"recoverTime" json:"recoverTime" toml:"recoverTime"`
}

type Config struct {
	// Rate the rate for request rate.
	Rate int `yaml:"rate" json:"rate" toml:"rate"`
	// Restrictions the conditions for executing current limiting.
	Restrictions []Rule `yaml:"restrictions" json:"restrictions" toml:"restrictions"`
	// GrayRecover Conditions for grayscale recovery request rate
	GrayRecover `yaml:"grayRecover" json:"grayRecover" toml:"grayRecover"`
}

// Metrics handle
// metrics tag and struct field cache.
var metricsFiledMap = make(map[string]reflect.StructField)

func init() {
	tp := reflect.TypeOf(Metrics{})
	for i := 0; i < tp.NumField(); i++ {
		f := tp.Field(i)
		jsonTag := f.Tag.Get("json")
		if jsonTag == "" || jsonTag == "-" {
			continue
		}

		jsonName := strings.Split(jsonTag, ",")[0]
		metricsFiledMap[jsonName] = f
	}
}

// _check the method to check ruleName.
func _check(rules []Rule) error {
	for _, rule := range rules {
		_, ok := metricsFiledMap[rule.RuleName]
		if !ok {
			return fmt.Errorf("rule `%s` not exists", rule.RuleName)
		}
	}

	return nil
}

func _parseThreshold(threshold json.Number) (any, error) {
	if i, err := strconv.Atoi(threshold.String()); err == nil {
		return i, nil
	}

	if f, err := strconv.ParseFloat(threshold.String(), 64); err == nil {
		return f, nil
	}

	return nil, errorx.ErrThresholdType
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
