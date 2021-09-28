/*
Copyright 2015-2021 Gravitational, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package desc

import (
	"fmt"
	"strings"
)

// Comment represents a wrapper over a string containing a comment
type Comment string

// SlashSlash appends "//" to comment
func (s Comment) SlashSlash(addSpace bool) string {
	var r []string

	l := strings.Split(string(s), "\n")

	for _, s := range l {
		if strings.TrimSpace(s) == "" {
			continue
		}

		v := "//%v"
		if addSpace {
			v = "// %v"
		}

		r = append(r, fmt.Sprintf(v, s))
	}

	return strings.Join(r, "\n")
}

// ToSingleLine returns multiline comment as a single string
func (s Comment) ToSingleLine() string {
	return strings.TrimSpace(strings.Join(strings.Split(string(s), "\n"), " "))
}
