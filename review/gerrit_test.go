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
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/devops-lintflow/lintflow/config"
	"github.com/devops-lintflow/lintflow/format"
)

func initHandle(t *testing.T) gerrit {
	c, err := initConfig("../tests/config.yml")
	assert.Equal(t, nil, err)

	return gerrit{
		r: c.Spec.Review,
	}
}

func TestClean(t *testing.T) {
	d, _ := os.Getwd()
	root := filepath.Join(d, "gerrit-test-clean")

	err := os.Mkdir(root, os.ModePerm)
	assert.Equal(t, nil, err)

	h := initHandle(t)
	err = h.Clean(root)
	assert.Equal(t, nil, err)
}

// nolint: dogsled
func TestFetch(t *testing.T) {
	h := initHandle(t)

	d, _ := os.Getwd()
	root := filepath.Join(d, "gerrit-test-fetch")

	_, _, _, err := h.Fetch(root, commitGerrit)
	assert.Equal(t, nil, err)

	err = h.Clean(root)
	assert.Equal(t, nil, err)
}

// nolint: dogsled
func TestQuery(t *testing.T) {
	var buf []interface{}
	var err error
	var ret []byte

	r := initHandle(t)

	buf, err = r.Query("change:"+strconv.Itoa(changeGerrit), 0)
	assert.Equal(t, nil, err)

	ret, _ = json.Marshal(buf)
	fmt.Println(string(ret))

	buf, err = r.Query(queryAfter+" "+queryBefore, 0)
	assert.Equal(t, nil, err)

	fmt.Println(len(buf))
}

// nolint: dogsled
func TestVote(t *testing.T) {
	h := initHandle(t)

	buf := make([]format.Format, 0)

	vote := config.Vote{
		Label:       "Lint-Verified",
		Approval:    "+1",
		Disapproval: "-1",
		Message:     "Voting Lint-Verified by lintflow",
	}

	err := h.Vote("", buf, vote)
	assert.NotEqual(t, nil, err)

	err = h.Vote(commitGerrit, buf, vote)
	assert.Equal(t, nil, err)

	buf = make([]format.Format, 1)
	buf[0] = format.Format{
		File:    "AndroidManifest.xml",
		Line:    1,
		Type:    format.TypeError,
		Details: "Disapproved",
	}

	err = h.Vote(commitGerrit, buf, vote)
	assert.Equal(t, nil, err)
}

func TestGetContent(t *testing.T) {
	h := initHandle(t)

	_, err := h.get(h.urlContent(-1, -1, ""))
	assert.NotEqual(t, nil, err)

	buf, err := h.get(h.urlContent(changeGerrit, revisionGerrit, url.PathEscape("AndroidManifest.xml")))
	assert.Equal(t, nil, err)

	dst := make([]byte, len(buf))
	n, _ := base64.StdEncoding.Decode(dst, buf)
	assert.NotEqual(t, 0, n)
}

func TestGetDetail(t *testing.T) {
	h := initHandle(t)

	_, err := h.get(h.urlDetail(-1))
	assert.NotEqual(t, nil, err)

	buf, err := h.get(h.urlDetail(changeGerrit))
	assert.Equal(t, nil, err)

	_, err = h.unmarshal(buf)
	assert.Equal(t, nil, err)
}

func TestGetFiles(t *testing.T) {
	h := initHandle(t)

	_, err := h.get(h.urlFiles(-1, -1))
	assert.NotEqual(t, nil, err)

	buf, err := h.get(h.urlFiles(changeGerrit, revisionGerrit))
	assert.Equal(t, nil, err)

	_, err = h.unmarshal(buf)
	assert.Equal(t, nil, err)
}

func TestGetPatch(t *testing.T) {
	h := initHandle(t)

	_, err := h.get(h.urlPatch(-1, -1))
	assert.NotEqual(t, nil, err)

	buf, err := h.get(h.urlPatch(changeGerrit, revisionGerrit))
	assert.Equal(t, nil, err)

	dst := make([]byte, len(buf))
	n, _ := base64.StdEncoding.Decode(dst, buf)
	assert.NotEqual(t, 0, n)
}

func TestGetQuery(t *testing.T) {
	h := initHandle(t)

	_, err := h.get(h.urlQuery("commit:-1", []string{"CURRENT_REVISION"}, 0))
	assert.NotEqual(t, nil, err)

	buf, err := h.get(h.urlQuery("commit:"+commitGerrit, []string{"CURRENT_REVISION"}, 0))
	assert.Equal(t, nil, err)

	_, err = h.unmarshalList(buf)
	assert.Equal(t, nil, err)
}

func TestPostReview(t *testing.T) {
	h := initHandle(t)

	err := h.post(h.urlReview(-1, -1), nil)
	assert.NotEqual(t, nil, err)

	buf := map[string]interface{}{
		"comments": map[string]interface{}{
			"AndroidManifest.xml": []map[string]interface{}{
				{
					"line":    1,
					"message": "Commented by lintflow",
				},
			},
		},
		"labels": map[string]interface{}{
			"Code-Review": -1,
		},
		"message": "Voting Code-Review by lintflow",
	}

	err = h.post(h.urlReview(changeGerrit, revisionGerrit), buf)
	assert.Equal(t, nil, err)
}
