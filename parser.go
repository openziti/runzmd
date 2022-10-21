/*
   Copyright NetFoundry Inc.

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

   https://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package runzmd

import (
	"strings"
	"unicode"
)

func ParseArgumentsWithStrings(val string) []string {
	var result []string

	current := &strings.Builder{}
	inString := false
	for _, r := range val {
		if r == '\'' {
			if inString {
				inString = false
			} else {
				inString = true
			}
		} else if !inString && unicode.IsSpace(r) {
			if current.Len() > 0 {
				result = append(result, current.String())
				current.Reset()
			}
		} else {
			current.WriteRune(r)
		}
	}

	if current.Len() > 0 {
		result = append(result, current.String())
	}

	return result
}
