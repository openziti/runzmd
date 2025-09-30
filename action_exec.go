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
	"github.com/fatih/color"
	"strings"
)

type ExecActionHandler struct{}

func (self *ExecActionHandler) Execute(ctx *ActionContext) error {
	if strings.EqualFold("true", ctx.Headers["templatize"]) {
		body, err := ctx.Runner.Template(ctx.Body)
		if err != nil {
			return err
		}
		ctx.Body = body
	}
	lines := strings.Split(ctx.Body, "\n")
	var cmds [][]string
	buf := &strings.Builder{}
	buf.WriteString("About to execute:\n\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if len(line) > 0 {
			params := ParseArgumentsWithStrings(line)
			params[0] = line
			cmds = append(cmds, params)
			ctx.Runner.LeftPadBuilder(buf)
			buf.WriteString("  ")
			buf.WriteString(color.New(color.Bold).Sprint(line))
			buf.WriteRune('\n')
		}
	}
	buf.WriteRune('\n')
	ctx.Runner.LeftPadBuilder(buf)
	buf.WriteString("Continue [Y/N] (default Y): ")

	if !ctx.Runner.AssumeDefault {
		Continue(buf.String(), true)
	}

	fmt.Println("")
	c := color.New(color.FgBlue, color.Bold)

	colorStdOut := !strings.EqualFold("false", ctx.Headers["colorStdOut"])

	allowRetry := strings.EqualFold("true", ctx.Headers["allowRetry"])
	failOk := strings.EqualFold("true", ctx.Headers["failOk"])
	for _, cmd := range cmds {
		_, _ = c.Printf("$ %v\n", cmd[0])
		done := false
		for !done {
			if err := Exec(cmd[0], colorStdOut, cmd[1:]...); err != nil {
				if failOk {
					return nil
				}
				if allowRetry {
					retry, err2 := AskYesNoWithDefault(fmt.Sprintf("operation failed with err: %v. Retry [Y/N] (default Y):", err), true)
					if err2 != nil {
						fmt.Printf("error while asking about retry: %v\n", err2)
						return err
					}
					if !retry {
						return err
					}
				} else {
					return err
				}
			} else {
				done = true
			}
		}
	}
	return nil
}
