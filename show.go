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
	"fmt"
	"github.com/alecthomas/chroma/quick"
	"github.com/pkg/errors"
)

func NewShowActionHandler() *ShowActionHandler {
	return &ShowActionHandler{
		data: map[string]string{},
	}
}

type ShowActionHandler struct {
	data map[string]string
}

func (self *ShowActionHandler) Execute(ctx *ActionContext) error {
	name, ok := ctx.Headers["src"]
	if !ok {
		return errors.New("no name specified in show")
	}

	content, ok := self.data[name]
	if !ok {
		return errors.Errorf("no contents found to show for name: '%v'", name)
	}

	if !ctx.Runner.AssumeDefault {
		view, err := AskYesNoWithDefault(fmt.Sprintf("View source for %v? [Y/N] (default N): ", name), false)
		if err != nil {
			return err
		}
		if !view {
			return nil
		}
	}

	codeWriter := NewSlowWriter(ctx.Runner.NewLinePause)

	fmt.Println("")
	fmt.Printf("%v:\n", name)
	for i := 0; i < 80; i++ {
		fmt.Print("-")
	}
	fmt.Println("")
	if lang, found := ctx.Headers["highlight"]; found {
		if err := quick.Highlight(codeWriter, content, lang, "terminal", "friendly"); err != nil {
			return err
		}
	} else {
		_, _ = fmt.Fprintln(codeWriter, content)
	}
	for i := 0; i < 80; i++ {
		fmt.Print("-")
	}
	fmt.Println("")
	return ctx.Runner.HandlePause(ctx)
}

func (self *ShowActionHandler) Add(name, contents string) {
	self.data[name] = contents
}
