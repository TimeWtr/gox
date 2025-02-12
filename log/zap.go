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

package log

import "go.uber.org/zap"

type ZapLogger struct {
	z *zap.Logger
}

func NewZapLogger(z *zap.Logger) *ZapLogger {
	return &ZapLogger{
		z: z,
	}
}

func (z *ZapLogger) Debugf(msg string, args ...Field) {
	z.z.Debug(msg, z.transfer(args...)...)
}

func (z *ZapLogger) Infof(msg string, args ...Field) {
	z.z.Info(msg, z.transfer(args...)...)
}

func (z *ZapLogger) Warnf(msg string, args ...Field) {
	z.z.Warn(msg, z.transfer(args...)...)
}

func (z *ZapLogger) Errorf(msg string, args ...Field) {
	z.z.Error(msg, z.transfer(args...)...)
}

func (z *ZapLogger) transfer(args ...Field) []zap.Field {
	res := make([]zap.Field, len(args))
	for _, arg := range args {
		res = append(res, zap.Any(arg.Key, arg.Value))
	}

	return res
}
