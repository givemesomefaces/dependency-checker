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

// HitResult is a single item that represents a dependency checked.
type HitResult struct {
	BlackDep  string
	ParentDep string
}

// Report is a collection of resolved HitResult.
type Report struct {
	Hit []*HitResult
}

// Resolve marks the dependency checked.
func (report *Report) Resolve(result *HitResult) {
	if result.ParentDep == "" {
		result.ParentDep = "-"
	}
	report.Hit = append(report.Hit, result)
	report.Hit = removeDuplicate(report.Hit)
}
func removeDuplicate(resultList []*HitResult) []*HitResult {
	resultMap := map[string]bool{}
	for _, v := range resultList {
		data, _ := json.Marshal(v)
		resultMap[string(data)] = true
	}
	var result []*HitResult
	for k := range resultMap {
		var t *HitResult
		json.Unmarshal([]byte(k), &t)
		result = append(result, t)
	}
	return result
}
func (report *Report) String() string {
	sort.SliceStable(report.Hit, func(i, j int) bool {
		return report.Hit[i].BlackDep < report.Hit[j].BlackDep
	})

	dWidth, lWidth := .0, .0
	for _, r := range report.Hit {
		dWidth = math.Max(float64(len(r.BlackDep)), dWidth)
		lWidth = math.Max(float64(len(r.ParentDep)), lWidth)
	}

	rowTemplate := fmt.Sprintf("%%-%dv | %%%dv\n", int(dWidth), int(lWidth))
	s := fmt.Sprintf(rowTemplate, "Black-List", "Path")
	s += fmt.Sprintf(rowTemplate, strings.Repeat("-", int(dWidth)), strings.Repeat("-", int(lWidth)))
	for _, r := range report.Hit {
		s += fmt.Sprintf(rowTemplate, r.BlackDep, r.ParentDep)
	}

	return s
}
