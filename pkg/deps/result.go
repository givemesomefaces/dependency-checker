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
	"encoding/json"
	"fmt"
	"math"
	"sort"
	"strings"
)

type SpdxID string

// Result is a single item that represents a resolved dependency license.
type Result struct {
	BlackDep  string
	ParentDep string
}

// Report is a collection of resolved Result.
type Report struct {
	Resolved []*Result
}

// Resolve marks the dependency's license is resolved.
func (report *Report) Resolve(result *Result) {
	if result.ParentDep == "" {
		result.ParentDep = "-"
	}
	report.Resolved = append(report.Resolved, result)
	report.Resolved = removeDuplicate(report.Resolved)
}
func removeDuplicate(resultList []*Result) []*Result {
	resultMap := map[string]bool{}
	for _, v := range resultList {
		data, _ := json.Marshal(v)
		resultMap[string(data)] = true
	}
	var result []*Result
	for k := range resultMap {
		var t *Result
		json.Unmarshal([]byte(k), &t)
		result = append(result, t)
	}
	return result
}
func (report *Report) String() string {
	sort.SliceStable(report.Resolved, func(i, j int) bool {
		return report.Resolved[i].BlackDep < report.Resolved[j].BlackDep
	})

	dWidth, lWidth := .0, .0
	for _, r := range report.Resolved {
		dWidth = math.Max(float64(len(r.BlackDep)), dWidth)
		lWidth = math.Max(float64(len(r.ParentDep)), lWidth)
	}

	rowTemplate := fmt.Sprintf("%%-%dv | %%%dv\n", int(dWidth), int(lWidth))
	s := fmt.Sprintf(rowTemplate, "Black-List", "Path")
	s += fmt.Sprintf(rowTemplate, strings.Repeat("-", int(dWidth)), strings.Repeat("-", int(lWidth)))
	for _, r := range report.Resolved {
		s += fmt.Sprintf(rowTemplate, r.BlackDep, r.ParentDep)
	}

	return s
}
