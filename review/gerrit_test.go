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

package review

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/craftslab/lintflow/proto"
)

func initHandle(t *testing.T) gerrit {
	c, err := initConfig("../tests/config.yml")
	assert.Equal(t, nil, err)

	var g gerrit

	for _, item := range c.Spec.Review {
		if item.Name == "gerrit" {
			g.r = item
			break
		}
	}

	return g
}

func TestClean(t *testing.T) {
	dir, _ := os.Getwd()
	name := filepath.Join(dir, "gerrit-"+commitGerrit)
	err := os.Mkdir(name, os.ModePerm)
	assert.Equal(t, nil, err)

	h := initHandle(t)
	err = h.Clean(name)
	assert.Equal(t, nil, err)
}

func TestFetch(t *testing.T) {
	h := initHandle(t)

	_, err := h.Fetch("")
	assert.NotEqual(t, nil, err)

	_, err = h.Fetch("invalid")
	assert.NotEqual(t, nil, err)

	name, err := h.Fetch(commitGerrit)
	assert.Equal(t, nil, err)

	if _, e := os.Stat(filepath.Join(name, proto.StorePatch)); os.IsNotExist(e) {
		assert.NotEqual(t, nil, e)
	}

	if _, e := os.Stat(filepath.Join(name, proto.StoreSource)); os.IsNotExist(e) {
		assert.NotEqual(t, nil, e)
	}

	err = h.Clean(name)
	assert.Equal(t, nil, err)
}

func TestVote(t *testing.T) {
	// TODO
	assert.Equal(t, nil, nil)
}

func TestGet(t *testing.T) {
	h := initHandle(t)

	_, err := h.get(-1)
	assert.NotEqual(t, nil, err)

	_, err = h.get(changeGerrit)
	assert.Equal(t, nil, err)
}

func TestQuery(t *testing.T) {
	h := initHandle(t)

	_, err := h.query("commit:-1", 0)
	assert.NotEqual(t, nil, err)

	_, err = h.query("commit:"+commitGerrit, 0)
	assert.Equal(t, nil, err)
}

func TestReviewGerrit(t *testing.T) {
	h := initHandle(t)

	err := h.review(-1, -1, []byte{})
	assert.NotEqual(t, nil, err)

	buf := []byte(`{
		"message": "Voting Code-Review by lintflow",
		"labels": {
			"Code-Review": -1
		},
		"comments": {
			"AndroidManifest.xml": [
				{
					"line": 1,
					"message": "Commented by lintflow"
				}
			]
		}
	}`)

	err = h.review(changeGerrit, revisionGerrit, buf)
	assert.Equal(t, nil, err)
}
