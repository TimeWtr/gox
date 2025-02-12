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
	"strconv"

	"github.com/TimeWtr/gox/errorx"
	"github.com/redis/go-redis/v9"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type Configuration interface {
	Set(ctx context.Context, latitude string, rate uint64) error
	Del(ctx context.Context, latitude string) error
}

const DefaultHashtableName = "metadata"

const (
	DefaultHashName          = "meta"
	DefaultCPUKey            = "CPU"
	DefaultMemoryKey         = "Memory"
	DefaultMemoryUsageKey    = "Memory Usage"
	DefaultMemoryUsedKey     = "Memory Used"
	DefaultRequestLatencyKey = "Request Latency"
	DefaultErrRateKey        = "Error Rate"
	DefaultActiveConnsKey    = "Active Conns"
)

// RedisConf Use redis as the configuration center to request
// current limiting original data and the storage structure is a hash
// table. the default table name is metadata.
type RedisConf struct {
	// redis client
	client redis.Cmdable
	// the name of hash table to store request rate metadata.
	hashTableName string
}

func NewRedisConfiguration(client redis.Cmdable, hashTableName ...string) Configuration {
	c := &RedisConf{
		client: client,
	}

	if len(hashTableName) > 0 {
		c.hashTableName = hashTableName[0]
	}

	return c
}

func (rc *RedisConf) Set(ctx context.Context, latitude string, rate uint64) error {
	_, err := rc.client.HSet(ctx, rc.hashTableName, latitude, rate).Result()
	return err
}

func (rc *RedisConf) Del(ctx context.Context, latitude string) error {
	n, err := rc.client.HDel(ctx, rc.hashTableName, latitude).Result()
	if err != nil || n != 1 {
		return errorx.ErrDelConfig
	}

	return nil
}

// EtcdConf Use etcd as the configuration center to request
// current limiting original data and the storage structure is a hash
// table.
type EtcdConf struct {
	// etcd client
	client *clientv3.Client
}

func NewEtcdConfiguration(client *clientv3.Client) Configuration {
	return &EtcdConf{
		client: client,
	}
}

func (e *EtcdConf) Set(ctx context.Context, latitude string, rate uint64) error {
	_, err := e.client.Put(ctx, latitude, strconv.FormatUint(rate, 10))
	if err != nil {
		return err
	}

	return nil
}

func (e *EtcdConf) Del(ctx context.Context, latitude string) error {
	_, err := e.client.Delete(ctx, latitude)
	return err
}
