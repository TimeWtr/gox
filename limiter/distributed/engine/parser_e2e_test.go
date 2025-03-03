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

var jsonContent = `{
  "base_threshold":1000,
  "min_threshold": 300,
  "strategy": "qps",
  "period": "1s",
  "priority": "high",
  "rules": [
    {
      "scope":{
        "type": "service",
        "value": "order_service"
      },
      "base_threshold":1000,
      "min_threshold": 300,
      "strategy": "qps",
      "priority": "medium",
      "period": "1s",
      "trigger": [
        {
          "metric": "cpu_usage",
          "threshold": 0.8
        },
        {
          "metric": "mem_usage",
          "threshold": 0.8
        },
        {
          "metric": "err_rate",
          "threshold": 0.2
        }
      ],
      "children": [
        {
          "scope": {
            "type": "api",
            "value": "/api/v1/order"
          },
          "base_threshold": 500,
          "min_threshold": 100,
          "strategy": "concurrency",
          "priority": "low",
          "period": "1s"
        },
        {
          "scope": {
            "type": "api",
            "value": "/api/v1/user"
          },
          "base_threshold": 300,
          "min_threshold": 100,
          "strategy": "qps",
          "priority": "low",
          "period": "1s",
          "children": [
            {
              "scope": {
                "type": "user",
                "value": "*"
              },
              "base_threshold": 5,
              "strategy": "total",
              "priority": "low",
              "period": "1m"
            },
            {
              "scope": {
                "type": "ip",
                "value": "*"
              },
              "base_threshold": 5,
              "priority": "low",
              "strategy": "total",
              "period": "1m"
            }
          ]
        }
      ]
    }
  ]
}
`

var yamlContent = `base_threshold: 1000
min_threshold: 300
strategy: qps
period: 1s
priority: high
rules:
  - scope:
      type: service
      value: order_service
    base_threshold: 1000
    min_threshold: 300
    strategy: qps
    priority: medium
    period: 1s
    trigger:
      - metric: cpu_usage
        threshold: 0.8
      - metric: mem_usage
        threshold: 0.8
      - metric: err_rate
        threshold: 0.2
    children:
      - scope:
          type: api
          value: /api/v1/order
        base_threshold: 500
        min_threshold: 100
        strategy: concurrency
        priority: low
        period: 1s
      - scope:
          type: api
          value: /api/v1/user
        base_threshold: 300
        min_threshold: 100
        strategy: qps
        priority: low
        period: 1s
        children:
          - scope:
              type: user
              value: "*"
            base_threshold: 5
            strategy: total
            priority: low
            period: 1m
          - scope:
              type: ip
              value: "*"
            base_threshold: 5
            priority: low
            strategy: total
            period: 1m`

var tomlContent = `base_threshold = 1000
min_threshold = 300
strategy = "qps"
period = "1s"
priority = "high"
[[rules]]
scope = { type = "service", value = "order_service" }
base_threshold = 1000
min_threshold = 300
strategy = "qps"
priority = "medium"
period = "1s"

[[rules.trigger]]
metric = "cpu_usage"
threshold = 0.8

[[rules.trigger]]
metric = "mem_usage"
threshold = 0.8

[[rules.trigger]]
metric = "err_rate"
threshold = 0.2

[[rules.children]]
scope = { type = "api", value = "/api/v1/order" }
base_threshold = 500
min_threshold = 100
strategy = "concurrency"
priority = "low"
period = "1s"

[[rules.children]]
scope = { type = "api", value = "/api/v1/user" }
base_threshold = 300
min_threshold = 100
strategy = "qps"
priority = "low"
period = "1s"

[[rules.children.children]]
scope = { type = "user", value = "*" }
base_threshold = 5
strategy = "total"
priority = "low"
period = "1m"

[[rules.children.children]]
scope = { type = "ip", value = "*" }
base_threshold = 5
priority = "low"
strategy = "total"
period = "1m"`

//func TestNewRedisSource_JSON(t *testing.T) {
//	client := redis.NewClient(&redis.Options{
//		Addr:     "127.0.0.1:6379",
//		Password: "root",
//	})
//
//	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
//	defer cancel()
//	_, err := client.Set(ctx, "test", []byte(jsonContent), 120).Result()
//	assert.Nil(t, err)
//
//	bs, err := NewRedisSource(client, "test", DataTypeJson).Read()
//	assert.Nil(t, err)
//
//	cf, err := NewJsonParser(bs).parse()
//	assert.Nil(t, err)
//	t.Logf("conf: %+v\n", cf)
//}
//
//func TestNewRedisSource_YAML(t *testing.T) {
//	client := redis.NewClient(&redis.Options{
//		Addr:     "127.0.0.1:6379",
//		Password: "root",
//	})
//
//	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
//	defer cancel()
//	_, err := client.Set(ctx, "test", []byte(yamlContent), 120).Result()
//	assert.Nil(t, err)
//
//	bs, err := NewRedisSource(client, "test", DataTypeYaml).Read()
//	assert.Nil(t, err)
//
//	cf, err := NewYamlParser(bs).parse()
//	assert.Nil(t, err)
//	t.Logf("conf: %+v\n", cf)
//}
//
//func TestNewRedisSource_TOML(t *testing.T) {
//	client := redis.NewClient(&redis.Options{
//		Addr:     "127.0.0.1:6379",
//		Password: "root",
//	})
//
//	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
//	defer cancel()
//	_, err := client.Set(ctx, "test", []byte(tomlContent), 120).Result()
//	assert.Nil(t, err)
//
//	bs, err := NewRedisSource(client, "test", DataTypeToml).Read()
//	assert.Nil(t, err)
//
//	p := NewTomlParser(bs)
//	cf, err := p.parse()
//	assert.Nil(t, err)
//	t.Logf("conf: %+v\n", cf)
//}
