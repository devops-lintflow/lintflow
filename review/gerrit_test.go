//go:build review_test

// go test -cover -covermode=atomic -parallel 2 -tags=review_test -v github.com/devops-lintflow/lintflow/review

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

	dir, repo, files, meta, patch, err := h.Fetch(root, commitGerrit)
	assert.Equal(t, nil, err)

	fmt.Printf("  dir: %s\n", dir)
	fmt.Printf(" repo: %s\n", repo)
	fmt.Printf("files: %v\n", files)
	fmt.Printf("meta: %s\n", meta)
	fmt.Printf("patch: %s\n", patch)

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
	fmt.Printf("change: %s\n", string(ret))

	buf, err = r.Query(queryAfter+" "+queryBefore, 0)
	assert.Equal(t, nil, err)

	ret, _ = json.Marshal(buf)
	fmt.Printf("change: %s\n", string(ret))
}

// nolint: dogsled
func TestVote(t *testing.T) {
	h := initHandle(t)

	buf := make([]format.Report, 1)
	buf[0] = format.Report{
		File:    "lintshell/test.sh",
		Line:    1,
		Type:    format.TypeError,
		Details: "Disapproved by gerrit",
	}

	vote := config.Vote{
		Label:       "Lint-Verified",
		Approval:    "+1",
		Disapproval: "-1",
		Message:     "Voting Lint-Verified by gerrit",
	}

	err := h.Vote(commitGerrit, buf, vote)
	assert.Equal(t, nil, err)
}

func TestGetMeta(t *testing.T) {
	h := initHandle(t)

	commit := "39fe82c424a319e9613126d2ef1c837e114440c5"
	number := 1

	data := `{
		"branch": "main",
		"owner": {
			"name": "name"
		},
		"project": "name",
		"updated": "2024-09-20 07:15:44.639000000",
		"revisions": {
			"39fe82c424a319e9613126d2ef1c837e114440c5": {
				"_number": 1,
				"commit": {
					"committer": {
						"tz": -480
					}
				}
			}
		}
	}`

	_query := map[string]interface{}{}
	err := json.Unmarshal([]byte(data), &_query)

	_meta, err := h.meta(commit, _query)
	assert.Equal(t, nil, err)

	dst := make([]byte, base64.StdEncoding.DecodedLen(len(_meta)))
	n, _ := base64.StdEncoding.Decode(dst, _meta)

	buf := map[string]interface{}{}
	err = json.Unmarshal(dst[:n], &buf)

	assert.Equal(t, nil, err)
	assert.Equal(t, _query["branch"], buf[metaBranch])
	assert.Equal(t, _query["owner"].(map[string]interface{})["name"].(string), buf[metaOwner].(map[string]interface{})[metaName])
	assert.Equal(t, _query["project"], buf[metaProject])
	assert.NotEqual(t, "", buf[metaUpdated])

	_, ok := buf[metaRevisions].(map[string]interface{})[commit]
	assert.Equal(t, true, ok)

	assert.Equal(t, number, int(buf[metaRevisions].(map[string]interface{})[commit].(map[string]interface{})[metaNumber].(float64)))
}

func TestGetContent(t *testing.T) {
	h := initHandle(t)

	_, err := h.get(h.urlContent(-1, -1, ""))
	assert.NotEqual(t, nil, err)

	buf, err := h.get(h.urlContent(changeGerrit, revisionGerrit, "lintshell/test.sh"))
	assert.Equal(t, nil, err)

	dst := make([]byte, len(buf))
	n, _ := base64.StdEncoding.Decode(dst, buf)
	assert.NotEqual(t, 0, n)

	fmt.Printf("content: %s\n", string(dst))
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

	fmt.Printf("patch: %s\n", string(dst))
}

func TestGetQuery(t *testing.T) {
	h := initHandle(t)

	_, err := h.get(h.urlQuery("commit:-1", []string{"ALL_REVISIONS"}, 0))
	assert.NotEqual(t, nil, err)

	buf, err := h.get(h.urlQuery("commit:"+commitGerrit, []string{"ALL_REVISIONS"}, 0))
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
			"lintshell/test.sh": []map[string]interface{}{
				{
					"line":    1,
					"message": "Commented by gerrit",
				},
			},
		},
		"labels": map[string]interface{}{
			"Code-Review": -1,
		},
		"message": "Voting Code-Review by gerrit",
	}

	err = h.post(h.urlReview(changeGerrit, revisionGerrit), buf)
	assert.Equal(t, nil, err)
}
