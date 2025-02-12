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
	"io/ioutil"

	"github.com/BurntSushi/toml"
	"github.com/TimeWtr/gox/errorx"
	"gopkg.in/yaml.v3"
)

// Parser the interface to parse rule config file.
type Parser interface {
	Parse(filepath string) ([]Rule, error)
}

func NewParser(fileType string) (Parser, error) {
	switch fileType {
	case "json":
		return NewJsonParser(), nil
	case "yaml":
		return NewYamlParser(), nil
	case "toml":
		return NewTomlParser(), nil
	default:
		return nil, errorx.ErrFileType
	}
}

type YamlParser struct{}

func NewYamlParser() Parser {
	return &YamlParser{}
}

func (y *YamlParser) Parse(filepath string) ([]Rule, error) {
	bs, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	var rules []Rule
	return rules, yaml.Unmarshal(bs, &rules)
}

type JsonParser struct{}

func NewJsonParser() Parser {
	return &JsonParser{}
}

func (j *JsonParser) Parse(filepath string) ([]Rule, error) {
	bs, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	var rules []Rule
	return rules, json.Unmarshal(bs, &rules)
}

type TomlParser struct{}

func NewTomlParser() Parser {
	return &TomlParser{}
}

func (t *TomlParser) Parse(filepath string) ([]Rule, error) {
	bs, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	var rules []Rule
	return rules, toml.Unmarshal(bs, &rules)
}
