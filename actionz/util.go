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

package actionz

import (
	"bytes"
	"fmt"
	"github.com/Jeffail/gabs/v2"
	"github.com/pkg/errors"
	"os"
	"os/exec"
	"path"
	"strings"
)

func zitiList(params ...string) ([]*gabs.Container, error) {
	result, err := runZitiJson(params...)
	if err != nil {
		return nil, err
	}
	return result.S("data").Children(), nil
}

func getZitiPath() (string, error) {
	zitiPath := os.Args[0]
	if _, file := path.Split(zitiPath); file == "ziti" {
		return zitiPath, nil
	}

	zitiPath, err := exec.LookPath("ziti")
	if err != nil {
		return "", errors.Wrap(err, "ziti executable not found in path")
	}
	return zitiPath, nil
}

func runZitiJson(params ...string) (*gabs.Container, error) {
	path, err := getZitiPath()
	if err != nil {
		return nil, err
	}

	cmd := exec.Command(path, params...)

	outCollector := &bytes.Buffer{}
	cmd.Stdout = outCollector
	if err = cmd.Run(); err != nil {
		return nil, errors.Wrapf(err, "error running ziti command 'ziti %v'", strings.Join(params, " "))
	}

	result, err := gabs.ParseJSON(outCollector.Bytes())
	if err != nil {
		return nil, errors.Wrapf(err, "error parsing JSON output from ziti command 'ziti %v'", strings.Join(params, " "))
	}
	return result, nil
}

func wrapGabs(c *gabs.Container) *gabsWrapper {
	return &gabsWrapper{Container: c}
}

type gabsWrapper struct {
	*gabs.Container
}

func (self *gabsWrapper) String(path string) string {
	child := self.Path(path)
	if child == nil || child.Data() == nil {
		return ""
	}
	return self.toString(child.Data())
}

func (self *gabsWrapper) Bool(path string) bool {
	child := self.Path(path)
	if child == nil || child.Data() == nil {
		return false
	}
	if val, ok := child.Data().(bool); ok {
		return val
	}
	return strings.EqualFold("true", fmt.Sprintf("%v", child.Data()))
}

func (self *gabsWrapper) toString(val interface{}) string {
	if val, ok := val.(string); ok {
		return val
	}
	return fmt.Sprintf("%v", val)
}

func (self *gabsWrapper) Float64(path string) float64 {
	child := self.Path(path)
	if child == nil || child.Data() == nil {
		return 0
	}
	if val, ok := child.Data().(float64); ok {
		return val
	}
	return 0
}

func (self *gabsWrapper) StringSlice(path string) []string {
	child := self.Path(path)
	if child == nil || child.Data() == nil {
		return nil
	}
	if val, ok := child.Data().([]string); ok {
		return val
	}
	if vals, ok := child.Data().([]interface{}); ok {
		var result []string
		for _, val := range vals {
			result = append(result, self.toString(val))
		}
		return result
	}
	return nil
}
