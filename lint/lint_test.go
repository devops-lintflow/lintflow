// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package lint

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/devops-lintflow/lintflow/format"
)

const (
	commitMeta  = "42-a4bc7bd.meta"
	commitPatch = "42-a4bc7bd.patch"
)

func TestEncode(t *testing.T) {
	l := lint{}

	name := "lintshell"
	root := "../tests/project"
	files := []string{"COMMIT_MSG", "lintshell/test.sh"}
	meta := commitMeta
	patch := commitPatch

	req, err := l.encode(name, root, files, meta, patch)
	assert.Equal(t, nil, err)
	assert.Equal(t, name, req.Name)
	assert.Equal(t, len(files), len(req.LintFiles))
	assert.Equal(t, meta, req.LintMeta.Path)
	assert.Equal(t, patch, req.LintPatch.Path)
}

func TestDecode(t *testing.T) {
	l := lint{}

	reply := &LintReply{
		Name: "lintshell",
	}

	reply.LintReports = []*LintReport{
		{
			File:    "lintshell/test.sh",
			Line:    1,
			Type:    format.TypeError,
			Details: "Disapproved by review",
		},
	}

	buf, err := l.decode(reply)
	assert.Equal(t, nil, err)

	_, ok := buf[reply.Name]
	assert.Equal(t, true, ok)
	assert.Equal(t, len(reply.LintReports), len(buf[reply.Name]))
	assert.Equal(t, reply.LintReports[0].File, buf[reply.Name][0].File)
}
