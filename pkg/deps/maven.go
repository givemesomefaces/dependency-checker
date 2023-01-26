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
	"bufio"
	"bytes"
	"encoding/xml"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/lvlifeng/eye/internal/logger"
)

type MavenPomChecker struct {
	maven string
	repo  string
}

// CanCheck determine whether the file can be checked by name of the file
func (resolver *MavenPomChecker) CanCheck(mavenPomFile string) bool {
	base := filepath.Base(mavenPomFile)
	logger.Log.Debugln("Base name:", base)
	return base == mavenPom
}

// Check all dependencies declared in the pom.xml file.
func (resolver *MavenPomChecker) Check(mavenPomFile string, config *ConfigDeps, report *Report) error {
	if err := os.Chdir(filepath.Dir(mavenPomFile)); err != nil {
		return err
	}

	if err := resolver.CheckMVN(); err != nil {
		return err
	}

	deps, err := resolver.LoadDependencies()
	if err != nil {
		// attempt to download dependencies
		if err = resolver.DownloadDeps(); err != nil {
			return fmt.Errorf("dependencies download error")
		}
		// load again
		deps, err = resolver.LoadDependencies()
		if err != nil {
			return err
		}
	}

	return resolver.CheckDependencies(deps, config, report)
}

// CheckMVN check available mavenPom tools, find local repositories and download all dependencies
func (resolver *MavenPomChecker) CheckMVN() error {
	if err := resolver.FindMaven("./mvnw"); err == nil {
		logger.Log.Debugln("mvnw is found, will use mvnw by default")
	} else if err := resolver.FindMaven("mvn"); err != nil {
		return fmt.Errorf("neither found mvnw nor mvn")
	}

	if err := resolver.FindLocalRepository(); err != nil {
		return fmt.Errorf("can not find the local repository: %v", err)
	}

	return nil
}

func (resolver *MavenPomChecker) FindMaven(execName string) error {
	if _, err := exec.Command(execName, "--version").Output(); err != nil {
		return err
	}

	resolver.maven = execName
	return nil
}

func (resolver *MavenPomChecker) FindLocalRepository() error {
	output, err := exec.Command(resolver.maven, "help:evaluate", "-Dexpression=settings.localRepository", "-q", "-DforceStdout").Output() // #nosec G204
	if err != nil {
		return err
	}

	resolver.repo = string(output)
	return nil
}

func (resolver *MavenPomChecker) DownloadDeps() error {
	cmd := exec.Command(resolver.maven, "dependency:resolve") // #nosec G204
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err == nil {
		return nil
	}
	// the failure may be caused by the lack of submodules, try to install it
	install := exec.Command(resolver.maven, "clean", "install", "-Dcheckstyle.skip=true", "-Drat.skip=true", "-Dmaven.test.skip=true") // #nosec G204
	install.Stdout = os.Stdout
	install.Stderr = os.Stderr

	return install.Run()
}

func (resolver *MavenPomChecker) LoadDependencies() ([]*Dependency, error) {
	buf := bytes.NewBuffer(nil)

	cmd := exec.Command(resolver.maven, "dependency:tree") // #nosec G204
	cmd.Stdout = bufio.NewWriter(buf)
	cmd.Stderr = os.Stderr

	logger.Log.Debugf("Running command: [%v], please wait", cmd.String())
	err := cmd.Run()
	if err != nil {
		return nil, err
	}

	deps := LoadDependencies(buf.Bytes())
	return deps, nil
}

// CheckDependencies checks the dependencies of the given dependencies
func (resolver *MavenPomChecker) CheckDependencies(deps []*Dependency, config *ConfigDeps, report *Report) error {
	for _, dep := range deps {
		func() {
			state := NotFound
			err := resolver.CheckBlackList(config, dep, report)
			if err != nil {
				logger.Log.Warnf("Failed to check the dependency of <%s>: %v\n", dep.Name(), state.String())
			}
		}()
	}
	return nil
}

func (resolver *MavenPomChecker) CheckBlackList(config *ConfigDeps, dep *Dependency, report *Report) error {
	hit := false
	var hitBlackDep ConfigDependency
	for _, blackDep := range config.BlackList {
		hitBlackDep = blackDep
		blackDepGroupIdRegexp := regexp.MustCompile(blackDep.GroupId)
		if blackDep.GroupId != "" &&
			blackDepGroupIdRegexp != nil &&
			blackDepGroupIdRegexp.MatchString(dep.GroupId) {
			if blackDep.ArtifactId == "" {
				hit = true
				break
			} else {
				blackDepArtifactIdRegexp := regexp.MustCompile(blackDep.ArtifactId)
				if blackDep.Version == "" &&
					blackDepArtifactIdRegexp != nil &&
					blackDep.ArtifactId != "" &&
					blackDepArtifactIdRegexp.MatchString(dep.ArtifactId) {
					hit = true
					break
				}
				blackDepVersionRegexp := regexp.MustCompile(blackDep.Version)
				if blackDep.Version != "" &&
					blackDep.ArtifactId != "" &&
					blackDepArtifactIdRegexp != nil &&
					blackDepArtifactIdRegexp.MatchString(dep.ArtifactId) &&
					dep.Version != "" &&
					blackDepVersionRegexp != nil &&
					blackDepVersionRegexp.MatchString(dep.Version) {
					hit = true
					break
				}
			}

		}
	}
	if hit {
		report.Resolve(&HitResult{
			BlackDep:  hitBlackDep.Name(mavenPom),
			ParentDep: dep.Parent,
		})
	}
	return nil
}

func LoadDependencies(data []byte) []*Dependency {
	depsTree := LoadDependenciesTree(data)

	cnt := 0
	for _, dep := range depsTree {
		cnt += dep.Count()
	}

	deps := make([]*Dependency, 0, cnt)

	var queue []*Dependency
	for _, depTree := range depsTree {
		queue = append(queue, depTree)

		for len(queue) > 0 {
			dep := queue[0]
			queue = queue[1:]

			deps = append(deps, dep.Clone())
			queue = append(queue, dep.TransitiveDeps...)
		}
	}
	return deps
}

func LoadDependenciesTree(data []byte) []*Dependency {
	type Elem struct {
		*Dependency
		level int
	}

	var stack []Elem
	unique := make(map[string]struct{})

	reFind := regexp.MustCompile(`(?im)^.*? ([| ]*)(\+-|\\-) (?P<gid>\b.+?):(?P<aid>\b.+?):(?P<packaging>\b.+)(:\b.+)?:(?P<Version>\b.+):(?P<scope>\b.+?)(?P<optional>\b.+?)?$`) //nolint:lll // can't break down regex
	rawDeps := reFind.FindAllSubmatch(data, -1)

	deps := make([]*Dependency, 0, len(rawDeps))
	for _, rawDep := range rawDeps {
		dep := &Dependency{
			GroupId:    string(rawDep[reFind.SubexpIndex("gid")]),
			ArtifactId: string(rawDep[reFind.SubexpIndex("aid")]),
			Packaging:  string(rawDep[reFind.SubexpIndex("packaging")]),
			Version:    string(rawDep[reFind.SubexpIndex("Version")]),
			Scope:      string(rawDep[reFind.SubexpIndex("scope")]),
		}

		if _, have := unique[dep.Path()]; have {
			continue
		}

		if dep.Scope == "test" || dep.Scope == "provided" || dep.Scope == "system" {
			continue
		}

		unique[dep.Path()] = struct{}{}

		level := len(rawDep[1]) / 3
		dependence := string(rawDep[2])

		if level == 0 {
			deps = append(deps, dep)

			if len(stack) != 0 {
				stack = stack[:0]
			}

			stack = append(stack, Elem{dep, level})
			continue
		}

		tail := stack[len(stack)-1]

		if level == tail.level {
			stack[len(stack)-1] = Elem{dep, level}
			dep.AppendParent(stack[len(stack)-2].Dependency)
			stack[len(stack)-2].TransitiveDeps = append(stack[len(stack)-2].TransitiveDeps, dep)
		} else {
			stack = append(stack, Elem{dep, level})
			dep.AppendParent(tail.Dependency)
			tail.TransitiveDeps = append(tail.TransitiveDeps, dep)
		}

		if dependence == `\-` {
			stack = stack[:len(stack)-1]
		}
	}
	return deps
}

const (
	FoundLicenseInPomHeader State = 1 << iota
	FoundLicenseInJarLicenseFile
	FoundLicenseInJarManifestFile
	NotFound State = 0
)

type State int

func (s *State) String() string {
	if *s == 0 {
		return "no possible license found"
	}

	var m []string

	if *s&FoundLicenseInPomHeader != 0 {
		m = append(m, "failed to resolve license found in pom header")
	}
	if *s&FoundLicenseInJarLicenseFile != 0 {
		m = append(m, "failed to resolve license file found in jar")
	}
	if *s&FoundLicenseInJarManifestFile != 0 {
		m = append(m, "failed to resolve license content from manifest file found in jar")
	}

	return strings.Join(m, " | ")
}

type Dependency struct {
	GroupId, ArtifactId, Version, Packaging, Scope, Parent string
	TransitiveDeps                                         []*Dependency
}

func (dep *Dependency) Clone() *Dependency {
	return &Dependency{
		GroupId:    dep.GroupId,
		ArtifactId: dep.ArtifactId,
		Version:    dep.Version,
		Packaging:  dep.Packaging,
		Scope:      dep.Scope,
		Parent:     dep.Parent,
	}
}

func (dep *Dependency) Count() int {
	cnt := 1
	for _, tDep := range dep.TransitiveDeps {
		cnt += tDep.Count()
	}
	return cnt
}

func (dep *Dependency) Path() string {
	return fmt.Sprintf("%v/%v/%v", strings.ReplaceAll(dep.GroupId, ".", "/"), dep.ArtifactId, dep.Version)
}

func (dep *Dependency) Pom() string {
	return fmt.Sprintf("%v-%v.pom", dep.ArtifactId, dep.Version)
}

func (dep *Dependency) Jar() string {
	return fmt.Sprintf("%v-%v.jar", dep.ArtifactId, dep.Version)
}

func (dep *Dependency) Name() string {
	return fmt.Sprintf("%v:%v", dep.GroupId, dep.ArtifactId)
}

func (dep *Dependency) AllName() string {
	return fmt.Sprintf("%v:%v:%v", dep.GroupId, dep.ArtifactId, dep.Version)
}

func (dep *Dependency) AppendParent(parentDep *Dependency) {
	if parentDep != nil {
		if parentDep.Parent != "" {
			dep.Parent = parentDep.Parent + " -> " + parentDep.AllName()
		} else {
			dep.Parent = parentDep.AllName()
		}
	}
}

// PomFile is used to extract license from the pom.xml file
type PomFile struct {
	XMLName  xml.Name      `xml:"project"`
	Licenses []*XMLLicense `xml:"licenses>license,omitempty"`
}

// Raw return raw data
func (pom *PomFile) Raw() string {
	var contents []string
	for _, l := range pom.Licenses {
		contents = append(contents, l.Raw())
	}
	return strings.Join(contents, "\n")
}

type XMLLicense struct {
	Name         string `xml:"name,omitempty"`
	URL          string `xml:"url,omitempty"`
	Distribution string `xml:"distribution,omitempty"`
	Comments     string `xml:"comments,omitempty"`
}

func (l *XMLLicense) Raw() string {
	return fmt.Sprintf(`License: {Name: %s, URL: %s, Distribution: %s, Comments: %s, }`, l.Name, l.URL, l.Distribution, l.Comments)
}
