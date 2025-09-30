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
	"fmt"
	"github.com/fatih/color"
	"github.com/openziti/runzmd"
	"github.com/pkg/errors"
	"os"
	"strings"
)

type ZitiUpdateConfigAction struct{}

func (self *ZitiUpdateConfigAction) Execute(ctx *runzmd.ActionContext) error {
	if strings.EqualFold("true", ctx.Headers["templatize"]) {
		body, err := ctx.Runner.Template(ctx.Body)
		if err != nil {
			return err
		}
		ctx.Body = body
	}
	name := ctx.Headers["name"]

	buf := &strings.Builder{}
	buf.WriteString("About to execute:\n\n")

	line := fmt.Sprintf("ziti edge update config %v --data '%v'", name, ctx.Body)
	params := runzmd.ParseArgumentsWithStrings(line)
	if params[0] != "ziti" {
		return errors.Errorf("invalid parameter for ziti action, must start with 'ziti': %v", ctx.Body)
	}
	params[0] = line
	ctx.Runner.LeftPadBuilder(buf)
	buf.WriteString("  ")
	buf.WriteString(color.New(color.Bold).Sprint(line))
	buf.WriteRune('\n')
	buf.WriteRune('\n')
	ctx.Runner.LeftPadBuilder(buf)
	buf.WriteString("Continue [Y/N] (default Y): ")

	if !ctx.Runner.AssumeDefault {
		runzmd.Continue(buf.String(), true)
	}

	fmt.Println("")
	c := color.New(color.FgBlue, color.Bold)

	colorStdOut := !strings.EqualFold("false", ctx.Headers["colorStdOut"])

	allowRetry := strings.EqualFold("true", ctx.Headers["allowRetry"])
	failOk := strings.EqualFold("true", ctx.Headers["failOk"])
	_, _ = c.Printf("$ %v\n", line)
	done := false
	for !done {
		if err := runzmd.Exec(os.Args[0], colorStdOut, "edge", "update", "config", name, "--data", ctx.Body); err != nil {
			if failOk {
				return nil
			}
			if allowRetry {
				retry, err2 := runzmd.AskYesNoWithDefault(fmt.Sprintf("operation failed with err: %v. Retry [Y/N] (default Y):", err), true)
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
	return nil
}
