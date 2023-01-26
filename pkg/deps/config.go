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
	"os"
	"path/filepath"
	"strings"
)

var (
	mavenPom  = "pom.xml"
	golangMod = "go.mod"
)

type ConfigDeps struct {
	BlackList []ConfigDependency `yaml:"black-list"`
	Files     []string           `yaml:"files"`
}

type ConfigDependency struct {
	// mavenPom project
	GroupId    string `yaml:"groupId"`
	ArtifactId string `yaml:"artifactId"`
	Version    string `yaml:"version"`

	// go project
	Path string `yaml:"path"`
}

func (dep *ConfigDependency) Name(project string) string {
	var names []string
	if project == mavenPom {
		names = dep.MavenProjectName()
	}
	if project == golangMod {
		names = dep.GolangProjectName()
	}
	return strings.Join(names, ":")
}

func (dep *ConfigDependency) MavenProjectName() []string {
	var names []string
	if dep.GroupId != "" {
		names = append(names, dep.GroupId)
	}
	if dep.ArtifactId != "" {
		names = append(names, dep.ArtifactId)
	}
	if dep.Version != "" {
		names = append(names, dep.Version)
	}
	return names
}

func (dep *ConfigDependency) GolangProjectName() []string {
	var names []string
	if dep.Path != "" {
		names = append(names, dep.GroupId)
	}
	if dep.Version != "" {
		names = append(names, dep.Version)
	}
	return names
}
func (config *ConfigDeps) Finalize(configFile string) error {
	_, err := filepath.Abs(configFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	return nil
}
