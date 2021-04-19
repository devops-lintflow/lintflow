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

	"github.com/craftslab/lintflow/config"
	"github.com/craftslab/lintflow/proto"
)

const (
	root = "../tests/gerrit-2021-03-06/21/c5d3440911e06ed4fc60252bd89e7756f9ae67ee"
)

// nolint: funlen
// nolint: goconst
func TestFilter(t *testing.T) {
	var f config.Filter
	var l lint

	f.Include.Extension = []string{".java", ".xml"}
	f.Include.File = []string{"message", "patch"}
	f.Include.Repo = []string{""}

	ret := l.filter(f, "", []string{".ext"})
	assert.Equal(t, 0, len(ret))

	ret = l.filter(f, "", []string{"foo.ext"})
	assert.Equal(t, 0, len(ret))

	ret = l.filter(f, "", []string{".java"})
	assert.Equal(t, 1, len(ret))

	ret = l.filter(f, "", []string{"foo.java"})
	assert.Equal(t, 1, len(ret))

	ret = l.filter(f, "", []string{".ext", ".java", ".xml", "foo.ext", "foo.java", "foo.xml"})
	assert.Equal(t, 4, len(ret))

	ret = l.filter(f, "", []string{"foo"})
	assert.Equal(t, 0, len(ret))

	ret = l.filter(f, "", []string{"message"})
	assert.Equal(t, 1, len(ret))

	ret = l.filter(f, "", []string{"foo", "message", "patch"})
	assert.Equal(t, 2, len(ret))

	ret = l.filter(f, "foo", []string{"foo", "message", "patch"})
	assert.Equal(t, 0, len(ret))

	f.Include.Repo = []string{"alpha", "beta"}

	ret = l.filter(f, "", []string{"foo", "message", "patch"})
	assert.Equal(t, 0, len(ret))

	ret = l.filter(f, "foo", []string{"foo", "message", "patch"})
	assert.Equal(t, 0, len(ret))

	ret = l.filter(f, "alpha", []string{"foo", "message", "patch"})
	assert.Equal(t, 2, len(ret))
}

func TestMarshal(t *testing.T) {
	var buf []string
	var l lint

	_, err := l.marshal(root, buf)
	assert.NotEqual(t, nil, err)

	buf = []string{"foo.base64"}
	_, err = l.marshal(root, buf)
	assert.NotEqual(t, nil, err)

	buf = []string{proto.Base64Patch}
	_, err = l.marshal(root, buf)
	assert.Equal(t, nil, err)
}
