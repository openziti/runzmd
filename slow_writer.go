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
	"io"
	"os"
	"time"
)

func NewSlowWriter(newlinePause time.Duration) io.Writer {
	return &slowWriter{
		delay: newlinePause,
	}
}

type slowWriter struct {
	delay time.Duration
}

func (self *slowWriter) Write(p []byte) (n int, err error) {
	var buf []byte
	written := 0
	for _, b := range p {
		buf = append(buf, b)
		if b == '\n' {
			time.Sleep(self.delay)
			n, err := os.Stdout.Write(buf)
			buf = buf[0:0]
			written += n
			if err != nil {
				return written, err
			}
		}
	}

	if len(buf) > 0 {
		n, err := os.Stdout.Write(buf)
		written += n
		if err != nil {
			return written, err
		}
	}
	return written, nil
}
