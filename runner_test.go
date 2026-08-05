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
	"testing"

	"github.com/fatih/color"
	markdown "github.com/openziti/go-term-markdown"
)

// TestEmitTemplatizedStripsAnsiFromTag verifies that a ${var} reference whose tag
// has been polluted with ANSI escape sequences (as syntax highlighting does) is
// still resolved to its variable value.
func TestEmitTemplatizedStripsAnsiFromTag(t *testing.T) {
	r := NewRunner()
	r.AddVariable("name", "MyValue")

	// mimics what a syntax highlighter produces for `${name}`: color codes wrap
	// and split the reference, leaving escape sequences inside the ${...} tag.
	input := "xxxx \x1b[1m\x1b[33m${\x1b[0m\x1b[34mname\x1b[0m\x1b[1m\x1b[33m}\x1b[0m"

	out := &strings.Builder{}
	if err := r.EmitTemplatized(input, out); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(out.String(), "MyValue") {
		t.Fatalf("expected output to contain expanded value, got %q", out.String())
	}
	if strings.Contains(out.String(), "name") {
		t.Fatalf("expected tag 'name' to be consumed, got %q", out.String())
	}
}

// TestEmitTemplatizedAfterHighlighting exercises the full path that caused the
// bug: a fenced code block with a language hint is syntax-highlighted during
// markdown rendering, and the ${var} inside it must still expand.
func TestEmitTemplatizedAfterHighlighting(t *testing.T) {
	// force color on so the renderer applies syntax highlighting regardless of
	// whether the test's stdout is a terminal.
	restore := color.NoColor
	color.NoColor = false
	defer func() { color.NoColor = restore }()

	r := NewRunner()
	r.AddVariable("name", "MyValue")

	rendered := markdown.Render("```bash\nxxxx ${name}\n```\n", r.LineWidth, r.LeftPad)

	out := &strings.Builder{}
	if err := r.EmitTemplatized(string(rendered), out); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(out.String(), "MyValue") {
		t.Fatalf("expected highlighted code block to expand ${name}, got %q", out.String())
	}
}

// TestEmitTemplatizedUnknownVariable ensures an unknown variable still errors,
// even after ANSI stripping.
func TestEmitTemplatizedUnknownVariable(t *testing.T) {
	r := NewRunner()

	out := &strings.Builder{}
	err := r.EmitTemplatized("hello ${missing}", out)
	if err == nil {
		t.Fatal("expected error for unknown template value, got nil")
	}
	if !strings.Contains(err.Error(), "missing") {
		t.Fatalf("expected error to mention the missing variable, got %v", err)
	}
}
