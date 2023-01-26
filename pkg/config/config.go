// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package config

import (
	"errors"
	"github.com/lvlifeng/eye/internal/logger"
	"github.com/lvlifeng/eye/pkg/deps"
	"gopkg.in/yaml.v3"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
)

type DependencyYaml struct {
	Deps deps.ConfigDeps `yaml:"dependency"`
}

func Parse(filename string, bytes []byte) (*DependencyYaml, error) {
	var config DependencyYaml
	if err := yaml.Unmarshal(bytes, &config); err != nil {
		return nil, err
	}

	if err := config.Deps.Finalize(filename); err != nil {
		return nil, err
	}

	return &config, nil
}

func (config *DependencyYaml) Dependencies() *deps.ConfigDeps {
	return &config.Deps
}

type Config interface {
	Dependencies() *deps.ConfigDeps
}

func NewConfigFromFile(filename string) (Config, error) {
	var err error
	var bytes []byte

	// attempt to read configuration from specified file
	logger.Log.Infoln("Loading configuration from file:", filename)

	if bytes, err = os.ReadFile(filename); err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	if os.IsNotExist(err) {
		logger.Log.Infof("Config file %s does not exist, using the default config: dependency-default.yaml", filename)
		var eyeAbsPath string
		if eyeAbsPath, err = EyeAbsPath(); err != nil {
			return nil, err
		}
		if bytes, err = os.ReadFile(path.Join(eyeAbsPath, "dependency-default.yaml")); err != nil &&
			!os.IsNotExist(err) {
			return nil, err
		}
	}

	var config Config
	if config, err = Parse(filename, bytes); err == nil {
		return config, nil
	}
	return config, nil
}

// EyeAbsPath find the root directory of eye when user execute dep-eye command.
func EyeAbsPath() (string, error) {
	compile := regexp.MustCompile("/bin/(darwin|linux|windows)/dep-eye")

	var depEyeAbsPath string
	var err error
	// find the root directory of dep-eye command from "../../dep-eye d check"
	if depEyeAbsPath, err = filepath.Abs(os.Args[0]); err != nil {
		return "", err
	}
	if !compile.MatchString(depEyeAbsPath) {
		// find the root directory of dep-eye command from environmental variable. e.g. "dep-eye d check"
		if depEyeAbsPath, err = exec.LookPath(os.Args[0]); err != nil {
			return "", err
		}
	}
	eyePathIndexes := compile.FindAllStringIndex(depEyeAbsPath, -1)
	if eyePathIndexes == nil {
		return "", errors.New("dep-eye command execute failed, please check environmental variable")
	}
	lastEyePathIndex := eyePathIndexes[len(eyePathIndexes)-1]
	return depEyeAbsPath[0:lastEyePathIndex[0]], nil
}
