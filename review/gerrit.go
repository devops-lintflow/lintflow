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
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

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

	name := filepath.Join(dir, "gerrit-"+commit)
	if err := os.Mkdir(name, os.ModePerm); err != nil {
		return "", errors.Wrap(err, "failed to make "+name)
	}

	patch := filepath.Join(name, proto.StorePatch)
	if err := os.Mkdir(patch, os.ModePerm); err != nil {
		return "", errors.Wrap(err, "failed to make "+patch)
	}

	source := filepath.Join(name, proto.StoreSource)
	if err := os.Mkdir(source, os.ModePerm); err != nil {
		return "", errors.Wrap(err, "failed to make "+source)
	}

	// TODO

	return name, nil
}

func (g *gerrit) Vote(commit string, data []proto.Format) error {
	// TODO
	return nil
}

func (g *gerrit) get(change int) (map[string]interface{}, error) {
	_url := g.r.Host + ":" + strconv.Itoa(g.r.Port) + "/changes/" + strconv.Itoa(change) + "/detail"
	if g.r.User != "" && g.r.Pass != "" {
		_url = g.r.Host + ":" + strconv.Itoa(g.r.Port) + "/a/changes/" + strconv.Itoa(change) + "/detail"
	}

	req, err := http.NewRequest(http.MethodGet, _url, nil)
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

	buf := map[string]interface{}{}

	if err := json.Unmarshal(data[4:], &buf); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal")
	}

	return buf, nil
}

func (g *gerrit) query(search string, start int) (map[string]interface{}, error) {
	_url := g.r.Host + ":" + strconv.Itoa(g.r.Port) + "/changes/"
	if g.r.User != "" && g.r.Pass != "" {
		_url = g.r.Host + ":" + strconv.Itoa(g.r.Port) + "/a/changes/"
	}

	req, err := http.NewRequest(http.MethodGet, _url, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to request")
	}

	if g.r.User != "" && g.r.Pass != "" {
		req.SetBasicAuth(g.r.User, g.r.Pass)
	}

	q := req.URL.Query()
	q.Add("o", "CURRENT_REVISION")
	q.Add("q", search)
	q.Add("start", strconv.Itoa(start))
	req.URL.RawQuery = q.Encode()

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

	buf := []map[string]interface{}{}

	if err := json.Unmarshal(data[4:], &buf); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal")
	}

	if len(buf) == 0 {
		return nil, errors.New("failed to match")
	}

	return buf[0], nil
}

func (g *gerrit) review(change, revision int, data []byte) error {
	_url := g.r.Host + ":" + strconv.Itoa(g.r.Port) + "/changes/" + strconv.Itoa(change) +
		"/revisions/" + strconv.Itoa(revision) + "/review"
	if g.r.User != "" && g.r.Pass != "" {
		_url = g.r.Host + ":" + strconv.Itoa(g.r.Port) + "/a/changes/" + strconv.Itoa(change) +
			"/revisions/" + strconv.Itoa(revision) + "/review"
	}

	req, err := http.NewRequest(http.MethodPost, _url, bytes.NewBuffer(data))
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
