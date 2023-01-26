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

package deps

import (
	"fmt"
	"github.com/lvlifeng/eye/internal/logger"
)

type Checker interface {
	CanCheck(string) bool
	Check(string, *ConfigDeps, *Report) error
}

var Checkers = []Checker{
	new(MavenPomChecker),
}

func Check(config *ConfigDeps, report *Report) error {
resolveFile:
	for _, file := range config.Files {
		for _, checker := range Checkers {
			if !checker.CanCheck(file) {
				continue
			}
			logger.Log.Infof("Start checking dependencies, please wait!")
			if err := checker.Check(file, config, report); err != nil {
				return err
			}
			continue resolveFile
		}
		return fmt.Errorf("unable to find a checker to check dependency declaration file: %v", file)
	}
	return nil
}
