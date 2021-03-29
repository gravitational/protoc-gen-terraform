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

package plugin

import (
	"strconv"
	"strings"

	"github.com/gogo/protobuf/protoc-gen-gogo/descriptor"
)

// appendSlashSlash appends "//" to comment
func appendSlashSlash(s string) string {
	var r []string

	l := strings.Split(s, "\n")

	for _, s := range l {
		if strings.Trim(s, " \n") != "" {
			r = append(r, "//"+s)
		}
	}

	return strings.Join(r, "\n")
}

// getLocationPath returns location path converted to string
func getLocationPath(l *descriptor.SourceCodeInfo_Location) string {
	s := make([]string, len(l.GetPath()))

	for i, v := range l.GetPath() {
		s[i] = strconv.Itoa(int(v))
	}

	return strings.Join(s, ",")
}
