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
	"encoding/json"
	"io/ioutil"
	"time"

	"github.com/redis/go-redis/v9"

	clientv3 "go.etcd.io/etcd/client/v3"
	"golang.org/x/net/context"

	"github.com/BurntSushi/toml"
	"github.com/TimeWtr/gox/errorx"
	"gopkg.in/yaml.v3"
)

type ConfigSourceType string

const (
	ConfigSourceTypeFile  ConfigSourceType = "file"
	ConfigSourceTypeEtcd  ConfigSourceType = "etcd"
	ConfigSourceTypeRedis ConfigSourceType = "redis"
)

type DataType string

const (
	DataTypeJson DataType = "json"
	DataTypeYaml DataType = "yaml"
	DataTypeToml DataType = "toml"
)

// Parser the interface to parse rule config file.
type Parser interface {
	// Parse the method to parse config metadata.
	Parse() (Config, error)
}

// ConfigSource the interface to adapt multi config source, such as
// local file, etcd, nacos etc.
type ConfigSource interface {
	Read() ([]byte, error)
	SourceType() ConfigSourceType
	DataType() DataType
}

var _ ConfigSource = (*FileSource)(nil)

// FileSource the rule metadata source based on file system.
type FileSource struct {
	filepath string
	dataType DataType
}

func NewFileSource(filepath string, dataType DataType) ConfigSource {
	return &FileSource{
		filepath: filepath,
		dataType: dataType,
	}
}

func (f *FileSource) Read() ([]byte, error) {
	bs, err := ioutil.ReadFile(f.filepath)
	if err != nil {
		return nil, err
	}

	return bs, nil
}

func (f *FileSource) SourceType() ConfigSourceType {
	return ConfigSourceTypeFile
}

func (f *FileSource) DataType() DataType {
	return f.dataType
}

// EtcdSource the rule metadata source based on etcd.
type EtcdSource struct {
	client   *clientv3.Client
	key      string
	dataType DataType
}

func NewEtcdSource(client *clientv3.Client, key string, dataType DataType) ConfigSource {
	return &EtcdSource{
		client:   client,
		key:      key,
		dataType: dataType,
	}
}

func (e *EtcdSource) Read() ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	resp, err := e.client.Get(ctx, e.key)
	if err != nil {
		return nil, err
	}

	if len(resp.Kvs) == 0 {
		return nil, errorx.ErrConfigNotExists
	}

	return resp.Kvs[0].Value, nil
}

func (e *EtcdSource) SourceType() ConfigSourceType {
	return ConfigSourceTypeEtcd
}

func (e *EtcdSource) DataType() DataType {
	return e.dataType
}

type RedisSource struct {
	// redis client
	client redis.Cmdable
	// the key to store metadata
	key string
	// data type
	dataType DataType
}

func NewRedisSource(client redis.Cmdable, key string, dataType DataType) ConfigSource {
	return &RedisSource{
		client:   client,
		key:      key,
		dataType: dataType,
	}
}

func (r *RedisSource) Read() ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := r.client.Get(ctx, r.key).Result()
	if err != nil {
		return nil, err
	}

	return []byte(res), nil
}

func (r *RedisSource) SourceType() ConfigSourceType {
	return ConfigSourceTypeRedis
}

func (r *RedisSource) DataType() DataType {
	return r.dataType
}

// NewParser the parser initialize method.
func NewParser(cs ConfigSource) (Parser, error) {
	bs, err := cs.Read()
	if err != nil {
		return nil, err
	}

	switch cs.DataType() {
	case "json":
		return NewJsonParser(bs), nil
	case "yaml":
		return NewYamlParser(bs), nil
	case "toml":
		return NewTomlParser(bs), nil
	default:
		return nil, errorx.ErrFileType
	}
}

// YamlParser yaml parser to parse yaml type data.
type YamlParser struct {
	bs []byte
}

func NewYamlParser(bs []byte) Parser {
	return &YamlParser{
		bs: bs,
	}
}

func (y *YamlParser) Parse() (Config, error) {
	var cfg Config
	err := yaml.Unmarshal(y.bs, &cfg)
	if err != nil {
		return Config{}, err
	}

	return cfg, checker(cfg)
}

// JsonParser json parser to parse json type data.
type JsonParser struct {
	bs []byte
}

func NewJsonParser(bs []byte) Parser {
	return &JsonParser{
		bs: bs,
	}
}

func (j *JsonParser) Parse() (Config, error) {
	var cfg Config
	err := json.Unmarshal(j.bs, &cfg)
	if err != nil {
		return Config{}, err
	}

	return cfg, checker(cfg)
}

// TomlParser toml parser to parse toml type data.
type TomlParser struct {
	bs []byte
}

func NewTomlParser(bs []byte) Parser {
	return &TomlParser{
		bs: bs,
	}
}

func (t *TomlParser) Parse() (Config, error) {
	var cfg Config
	err := toml.Unmarshal(t.bs, &cfg)
	if err != nil {
		return Config{}, err
	}

	return cfg, checker(cfg)
}
