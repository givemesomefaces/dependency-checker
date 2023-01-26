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

package commands

import (
	"fmt"
	"github.com/lvlifeng/eye/internal/logger"
	"github.com/lvlifeng/eye/pkg/deps"
	"github.com/spf13/cobra"
)

var DepsCheckCommand = &cobra.Command{
	Use:     "check",
	Aliases: []string{"r"},
	Long:    "check all dependencies of a module and their transitive dependencies",
	RunE: func(cmd *cobra.Command, args []string) error {
		report := deps.Report{}

		configDeps := Config.Dependencies()
		if len(configDeps.BlackList) != 0 {
			if err := deps.Resolve(configDeps, &report); err != nil {
				return err
			}
		}
		logger.Log.Infof("Checking dependencies completed!")
		if len(report.Hit) != 0 {
			fmt.Println(report.String())
			return fmt.Errorf("found %d dependencies hit the blacklist", len(report.Hit))
		}
		return nil
	},
}
