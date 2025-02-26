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
	"errors"
	"fmt"

	"golang.org/x/sync/errgroup"
)

const (
	// ScopeTypeService type is service
	ScopeTypeService = "service"
	// ScopeTypeAPI type is api
	ScopeTypeAPI = "api"
	// ScopeTypeUser type is user
	ScopeTypeUser = "user"
	// ScopeTypeIP type is ip
	ScopeTypeIP = "ip"
)

type PriorityType string

const (
	PriorityTypeLow    PriorityType = "low"
	PriorityTypeMedium PriorityType = "medium"
	PriorityTypeHigh   PriorityType = "high"
)

func (p *PriorityType) String() string {
	return string(*p)
}

func (p *PriorityType) valid() error {
	if p == nil {
		return errors.New("priority not valid")
	}

	switch *p {
	case PriorityTypeLow, PriorityTypeMedium, PriorityTypeHigh:
		return nil
	default:
		return fmt.Errorf("priority type %s not valid", *p)
	}
}

type StrategyType string

const (
	StrategyQPS         StrategyType = "qps"
	StrategyConcurrency StrategyType = "concurrency"
	StrategyTotal       StrategyType = "total"
)

func (s *StrategyType) String() string {
	return string(*s)
}

func (s *StrategyType) valid() error {
	switch *s {
	case StrategyQPS, StrategyConcurrency, StrategyTotal:
		return nil
	default:
		return fmt.Errorf("strategy type %s not valid", s.String())
	}
}

type PeriodType string

func (p *PeriodType) String() string {
	return string(*p)
}

func (p *PeriodType) valid() error {
	if len(*p) == 0 {
		return errors.New("period value must not be empty")
	}

	_, err := parseTime(string(*p))
	return err
}

type Conf struct {
	BaseThreshold uint64       `json:"base_threshold" yaml:"base_threshold" toml:"base_threshold"`
	MinThreshold  uint64       `json:"min_threshold" yaml:"min_threshold" toml:"min_threshold"`
	Strategy      StrategyType `json:"strategy" yaml:"strategy" toml:"strategy"`
	Period        PeriodType   `json:"period" yaml:"period" toml:"period"`
	Priority      PriorityType `json:"priority" yaml:"priority" toml:"priority"`
	Rules         []Rule2      `json:"rules" yaml:"rules" toml:"rules"`
}

func (c *Conf) Check() error {
	// global strategy value valid
	err := c.Strategy.valid()
	if err != nil {
		return err
	}

	// global period value valid
	err = c.Period.valid()
	if err != nil {
		return err
	}

	// global priority value valid
	err = c.Priority.valid()
	if err != nil {
		return err
	}

	if len(c.Rules) == 0 {
		return errors.New("rules must not be empty")
	}

	// check all rules
	var eg errgroup.Group
	for _, rule := range c.Rules {
		eg.Go(func() error {
			return rule.check()
		})
	}

	return eg.Wait()
}

type Rule2 struct {
	Scope         Scope         `json:"scope" yaml:"scope" toml:"scope"`
	BaseThreshold uint64        `json:"base_threshold" yaml:"base_threshold" toml:"base_threshold"`
	MinThreshold  uint64        `json:"min_threshold" yaml:"min_threshold" toml:"min_threshold"`
	Strategy      StrategyType  `json:"strategy" yaml:"strategy" toml:"strategy"`
	Period        PeriodType    `json:"period" yaml:"period" toml:"period"`
	Priority      PriorityType  `json:"priority" yaml:"priority" toml:"priority"`
	Trigger       TriggerType   `json:"trigger" yaml:"trigger" toml:"trigger"`
	Algorithm     AlgorithmType `json:"algorithm" yaml:"algorithm" toml:"algorithm"`
	Children      []Rule2       `json:"children" yaml:"children" toml:"children"`
}

func (r *Rule2) check() error {
	// check scope type
	err := r.Scope.valid()
	if err != nil {
		return err
	}

	// check strategy value
	err = r.Strategy.valid()
	if err != nil {
		return err
	}

	// check period value
	err = r.Period.valid()
	if err != nil {
		return err
	}

	// check rule's priority value valid
	err = r.Priority.valid()
	if err != nil {
		return err
	}

	// check algorithm value valid
	if r.Scope.Value == ScopeTypeUser || r.Scope.Value == ScopeTypeIP {
		if err = r.Algorithm.Valid(r.Scope.Value); err != nil {
			return err
		}
	} else {
		if err = r.Algorithm.Valid(); err != nil {
			return err
		}
	}

	// check limit trigger
	if len(r.Trigger) != 0 {
		//for _, t := range r.Trigger {
		//	err = t.valid()
		//	if err != nil {
		//		return err
		//	}
		//}
		if err = r.Trigger.valid(); err != nil {
			return err
		}
	}

	if r.Children == nil {
		return nil
	}

	for _, child := range r.Children {
		err = child.check()
		if err != nil {
			return err
		}
	}

	return nil
}

type Scope struct {
	Type  string `json:"type" yaml:"type" toml:"type"`
	Value string `json:"value" yaml:"value" toml:"value"`
}

func (s *Scope) valid() error {
	switch s.Type {
	case ScopeTypeService:
		if s.Value == "" {
			return fmt.Errorf("service scope must have a value")
		}
	case ScopeTypeAPI:
		if s.Value == "" {
			return fmt.Errorf("api scope must have a value")
		}
	case ScopeTypeUser:
		if s.Value == "" {
			return fmt.Errorf("user scope must have a value")
		}
	case ScopeTypeIP:
		if s.Value == "" {
			return fmt.Errorf("ip scope must have a value")
		}
	default:
		return errors.New("unknown scope type")
	}

	return nil
}

//type TriggerItem struct {
//	Metric    string      `json:"metric" yaml:"metric" toml:"metric"`
//	Threshold json.Number `json:"threshold" yaml:"threshold" toml:"threshold"`
//}

type TriggerType string

func (t *TriggerType) valid() error {
	//_, ok := metricsMap[t.Metric]
	//if !ok {
	//	return fmt.Errorf("metric `%s` is not supported", t.Metric)
	//}
	//
	return nil
}

type Metrics struct {
	// used cpu percent,
	CPUUsage float64 `json:"cpu_usage,omitempty"`
	// memory percent
	MemUsage float64 `json:"mem_usage,omitempty"`
	// used memory size, unit is bytes.
	MemUsed uint64 `json:"mem_used,omitempty"`
	// request latency, only used while type is API.
	RequestLatency float64 `json:"request_latency,omitempty"`
	// request error rate.
	ErrRate float64 `json:"err_rate,omitempty"`
	// current active connections.
	ActiveConns uint64 `json:"active_conns,omitempty"`
}
