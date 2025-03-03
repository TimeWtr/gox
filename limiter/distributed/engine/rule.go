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

	_ "golang.org/x/sync/errgroup"
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
	RedisCluster RedisCluster `json:"redis_cluster" toml:"redis_cluster" yaml:"redis_cluster"`
	Rules        Rule         `json:"rules" yaml:"rules" toml:"rules"`
}

func (c *Conf) Check() error {
	if err := c.RedisCluster.Check(); err != nil {
		return err
	}

	if c.Rules.Children == nil {
		return nil
	}

	// fast path
	if len(c.Rules.Children) == 1 {
		return c.Rules.Children[0].check()
	}

	// low path
	var eg errgroup.Group
	for _, rule := range c.Rules.Children {
		eg.Go(func() error {
			return rule.check()
		})
	}

	return eg.Wait()
}

type RedisCluster struct {
	Addr []string `json:"addr" toml:"addr" yaml:"addr"`
}

func (c *RedisCluster) Check() error {
	if len(c.Addr) == 0 {
		return errors.New("redis cluster address must not be empty")
	}

	return nil
}

type Rule struct {
	Scope         Scope         `json:"scope,omitempty" yaml:"scope,omitempty" toml:"scope,omitempty"`
	BaseThreshold uint64        `json:"base_threshold" yaml:"base_threshold" toml:"base_threshold"`
	MinThreshold  uint64        `json:"min_threshold" yaml:"min_threshold" toml:"min_threshold"`
	Strategy      StrategyType  `json:"strategy" yaml:"strategy" toml:"strategy"`
	Period        PeriodType    `json:"period" yaml:"period" toml:"period"`
	Priority      PriorityType  `json:"priority" yaml:"priority" toml:"priority"`
	Trigger       TriggerType   `json:"trigger,omitempty" yaml:"trigger,omitempty" toml:"trigger,omitempty"`
	TriggerAST    Expr          `json:"-"` // parse and generate ast
	Algorithm     AlgorithmType `json:"algorithm,omitempty" yaml:"algorithm,omitempty" toml:"algorithm,omitempty"`
	Children      []Rule        `json:"children" yaml:"children" toml:"children"`
}

func (r *Rule) check() error {
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

type RuleTreeInter interface {
	GetScope() Scope
	GetBaseThreshold() uint64
	GetMinThreshold() uint64
	GetPeriod() PeriodType
	GetPriority() PriorityType
	GetTriggerAST() Expr
	GetAlgorithm() AlgorithmType
	GetChildren() []RuleTree
}

var _ RuleTreeInter = (*RuleTree)(nil)

type RuleTree struct {
	scope         Scope
	baseThreshold uint64
	minThreshold  uint64
	strategy      StrategyType
	period        PeriodType
	priority      PriorityType
	triggerAST    Expr
	algorithm     AlgorithmType
	children      []RuleTree
}

func (r *RuleTree) GetScope() Scope {
	return r.scope
}

func (r *RuleTree) GetBaseThreshold() uint64 {
	return r.baseThreshold
}

func (r *RuleTree) GetMinThreshold() uint64 {
	return r.minThreshold
}

func (r *RuleTree) GetPeriod() PeriodType {
	return r.period
}

func (r *RuleTree) GetPriority() PriorityType {
	return r.priority
}

func (r *RuleTree) GetTriggerAST() Expr {
	return r.triggerAST
}

func (r *RuleTree) GetAlgorithm() AlgorithmType {
	return r.algorithm
}

func (r *RuleTree) GetChildren() []RuleTree {
	return r.children
}

// BuildRuleTrees the method to build the rule trees.
func BuildRuleTrees(r Rule) ([]RuleTree, error) {
	return builder(r)
}

func builder(rs Rule) ([]RuleTree, error) {
	if rs.BaseThreshold == 0 {
		return nil, errors.New("rule must not be nil")
	}

	var trees []RuleTree

	rt := &RuleTree{
		scope:         rs.Scope,
		baseThreshold: rs.BaseThreshold,
		minThreshold:  rs.MinThreshold,
		strategy:      rs.Strategy,
		period:        rs.Period,
		priority:      rs.Priority,
	}

	if rs.Trigger != "" {
		expr, err := parseTrigger(string(rs.Trigger))
		if err != nil {
			return nil, err
		}
		rt.triggerAST = expr
	}

	if rs.Children != nil {
		for _, child := range rs.Children {
			tree, er := builder(child)
			if er != nil {
				return nil, er
			}
			rt.children = tree
		}
	}
	trees = append(trees, *rt)

	return trees, nil
}
