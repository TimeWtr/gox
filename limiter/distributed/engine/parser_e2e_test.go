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
  "rate": 1000,
  "restrictions": [
    {
      "ruleName": "cpu_usage",
      "threshold": 80,
      "action": "decrease",
      "amount": 20
    },
    {
      "ruleName": "mem_usage",
      "threshold": 80,
      "action": "decrease",
      "amount": 10
    },
    {
      "ruleName": "request_latency",
      "threshold": 200,
      "action": "decrease",
      "amount": 15
    },
    {
      "ruleName": "err_rate",
      "threshold": 0.1,
      "action": "decrease",
      "amount": 10
    }
  ],
  "grayRecover": {
    "grayScale": 0.1,
    "recoverTime": 60
  }
}
`

var yamlContent = `"rate": 1000
"restrictions":
  - "ruleName": "cpu_usage"
    "threshold": 80
    "action": "decrease"
    "amount": 20
  - "ruleName": "mem_usage"
    "threshold": 85
    "action": "decrease"
    "amount": 10
  - "ruleName": "request_latency"
    "threshold": 200
    "action": "decrease"
    "amount": 15
"grayRecover":
  "grayScale": 0.1
  "recoverTime": 60`

var tomlContent = `rate = 1000

[[restrictions]]
ruleName = "cpu_usage"
threshold = 80
action = "decrease"
amount = 20

[[restrictions]]
ruleName = "mem_usage"
threshold = 80
action = "decrease"
amount = 10

[[restrictions]]
ruleName = "request_latency"
threshold = 200
action = "decrease"
amount = 15

[grayRecover]
grayScale = 0.1
recoverTime = 60`

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
//	t.Log(string(bs))
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
//	t.Log(string(bs))
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
//	t.Log(string(bs))
//}
