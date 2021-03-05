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
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/pkg/errors"

	"github.com/craftslab/lintflow/config"
	"github.com/craftslab/lintflow/proto"
)

type gerrit struct {
	r config.Review
}

func (g *gerrit) Clean(name string) error {
	if err := os.RemoveAll(name); err != nil {
		return errors.Wrap(err, "failed to clean")
	}

	return nil
}

func (g *gerrit) Fetch(commit string) (string, error) {
	dir, _ := os.Getwd()

	n := filepath.Join(dir, "gerrit-"+commit)
	if err := os.Mkdir(n, os.ModePerm); err != nil {
		return "", errors.Wrap(err, "failed to make "+n)
	}

	p := filepath.Join(n, proto.StorePatch)
	if err := os.Mkdir(p, os.ModePerm); err != nil {
		return "", errors.Wrap(err, "failed to make "+p)
	}

	s := filepath.Join(n, proto.StoreSource)
	if err := os.Mkdir(s, os.ModePerm); err != nil {
		return "", errors.Wrap(err, "failed to make "+s)
	}

	// TODO

	return n, nil
}

func (g *gerrit) Vote(commit string, data []proto.Format) error {
	helper := func() (map[string]interface{}, map[string]interface{}, string) {
		if len(data) == 0 {
			return nil, map[string]interface{}{g.r.Vote.Label: g.r.Vote.Approval}, g.r.Vote.Message
		}
		c := map[string]interface{}{}
		for _, item := range data {
			b := map[string]interface{}{"line": item.Line, "message": item.Details}
			if _, p := c[item.File]; !p {
				c[item.File] = []map[string]interface{}{b}
			} else {
				c[item.File] = append(c[item.File].([]map[string]interface{}), b)
			}
		}
		return c, map[string]interface{}{g.r.Vote.Label: g.r.Vote.Disapproval}, g.r.Vote.Message
	}

	r, err := g.get(g.urlQuery("commit:"+commit, "CURRENT_REVISION", 0))
	if err != nil {
		return errors.Wrap(err, "failed to query")
	}

	ret, err := g.unmarshalList(r)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal")
	}

	revisions := ret["revisions"].(map[string]interface{})
	current := revisions[ret["current_revision"].(string)].(map[string]interface{})

	comments, labels, message := helper()
	buf := map[string]interface{}{"comments": comments, "labels": labels, "message": message}

	if err := g.post(g.urlReview(int(ret["_number"].(float64)), int(current["_number"].(float64))), buf); err != nil {
		return errors.Wrap(err, "failed to review")
	}

	return nil
}

func (g *gerrit) unmarshal(data []byte) (map[string]interface{}, error) {
	buf := map[string]interface{}{}

	if err := json.Unmarshal(data[4:], &buf); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal")
	}

	return buf, nil
}

func (g *gerrit) unmarshalList(data []byte) (map[string]interface{}, error) {
	var buf []map[string]interface{}

	if err := json.Unmarshal(data[4:], &buf); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal")
	}

	if len(buf) == 0 {
		return nil, errors.New("failed to match")
	}

	return buf[0], nil
}

func (g *gerrit) urlContent(change, revision int, name string) string {
	buf := g.r.Host + ":" + strconv.Itoa(g.r.Port) + "/changes/" + strconv.Itoa(change) +
		"/revisions/" + strconv.Itoa(revision) + "/files/" + name + "/content"

	if g.r.User != "" && g.r.Pass != "" {
		buf = g.r.Host + ":" + strconv.Itoa(g.r.Port) + "/a/changes/" + strconv.Itoa(change) +
			"/revisions/" + strconv.Itoa(revision) + "/files/" + name + "/content"
	}

	return buf
}

func (g *gerrit) urlDetail(change int) string {
	buf := g.r.Host + ":" + strconv.Itoa(g.r.Port) + "/changes/" + strconv.Itoa(change) + "/detail"

	if g.r.User != "" && g.r.Pass != "" {
		buf = g.r.Host + ":" + strconv.Itoa(g.r.Port) + "/a/changes/" + strconv.Itoa(change) + "/detail"
	}

	return buf
}

func (g *gerrit) urlPatch(change, revision int) string {
	buf := g.r.Host + ":" + strconv.Itoa(g.r.Port) + "/changes/" + strconv.Itoa(change) +
		"/revisions/" + strconv.Itoa(revision) + "/patch"

	if g.r.User != "" && g.r.Pass != "" {
		buf = g.r.Host + ":" + strconv.Itoa(g.r.Port) + "/a/changes/" + strconv.Itoa(change) +
			"/revisions/" + strconv.Itoa(revision) + "/patch"
	}

	return buf
}

func (g *gerrit) urlQuery(search, option string, start int) string {
	query := "?q=" + search + "&o=" + option + "&n=" + strconv.Itoa(start)

	buf := g.r.Host + ":" + strconv.Itoa(g.r.Port) + "/changes/" + query
	if g.r.User != "" && g.r.Pass != "" {
		buf = g.r.Host + ":" + strconv.Itoa(g.r.Port) + "/a/changes/" + query
	}

	return buf
}

func (g *gerrit) urlReview(change, revision int) string {
	buf := g.r.Host + ":" + strconv.Itoa(g.r.Port) + "/changes/" + strconv.Itoa(change) +
		"/revisions/" + strconv.Itoa(revision) + "/review"

	if g.r.User != "" && g.r.Pass != "" {
		buf = g.r.Host + ":" + strconv.Itoa(g.r.Port) + "/a/changes/" + strconv.Itoa(change) +
			"/revisions/" + strconv.Itoa(revision) + "/review"
	}

	return buf
}

func (g *gerrit) get(url string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to request")
	}

	if g.r.User != "" && g.r.Pass != "" {
		req.SetBasicAuth(g.r.User, g.r.Pass)
	}

	rsp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to do")
	}

	defer func() {
		_ = rsp.Body.Close()
	}()

	if rsp.StatusCode != http.StatusOK {
		return nil, errors.New("invalid status")
	}

	data, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read")
	}

	return data, nil
}

func (g *gerrit) post(url string, data map[string]interface{}) error {
	buf, err := json.Marshal(data)
	if err != nil {
		return errors.Wrap(err, "failed to marshal")
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(buf))
	if err != nil {
		return errors.Wrap(err, "failed to request")
	}

	req.Header.Set("Content-Type", "application/json;charset=utf-8")

	if g.r.User != "" && g.r.Pass != "" {
		req.SetBasicAuth(g.r.User, g.r.Pass)
	}

	rsp, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "failed to do")
	}

	defer func() {
		_ = rsp.Body.Close()
	}()

	if rsp.StatusCode != http.StatusOK {
		return errors.New("invalid status")
	}

	_, err = ioutil.ReadAll(rsp.Body)
	if err != nil {
		return errors.Wrap(err, "failed to read")
	}

	return nil
}
