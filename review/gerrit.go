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
	review config.Review
}

func (g *gerrit) Clean(name string) error {
	if err := os.RemoveAll(name); err != nil {
		return errors.Wrap(err, "failed to clean")
	}

	return nil
}

func (g *gerrit) Fetch(commit string) (string, error) {
	name, err := ioutil.TempDir(os.TempDir(), "*"+"-gerrit-"+commit)
	if err != nil {
		return "", errors.Wrap(err, "failed to tempdir")
	}

	patch := filepath.Join(name, proto.STORE_PATCH)
	if err := os.Mkdir(patch, os.ModePerm); err != nil {
		return "", errors.Wrap(err, "failed to make patch")
	}

	source := filepath.Join(name, proto.STORE_SOURCE)
	if err := os.Mkdir(source, os.ModePerm); err != nil {
		return "", errors.Wrap(err, "failed to make source")
	}

	// TODO

	return name, nil
}

func (g *gerrit) Vote(commit string, data []proto.Format) error {
	// TODO
	return nil
}

func (g *gerrit) get(id int) (map[string]interface{}, error) {
	_url := g.review.Host + ":" + strconv.Itoa(g.review.Port) + "/changes/" + strconv.Itoa(id) + "/detail"
	if g.review.User != "" && g.review.Pass != "" {
		_url = g.review.Host + ":" + strconv.Itoa(g.review.Port) + "/a/changes/" + strconv.Itoa(id) + "/detail"
	}

	req, err := http.NewRequest(http.MethodGet, _url, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to request")
	}

	if g.review.User != "" && g.review.Pass != "" {
		req.SetBasicAuth(g.review.User, g.review.Pass)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to do")
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("invalid status")
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read")
	}

	buf := map[string]interface{}{}

	if err := json.Unmarshal(data[4:], &buf); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal")
	}

	return buf, nil
}

func (g *gerrit) put(id int) error {
	// TODO
	return nil
}

func (g *gerrit) query(search string, start int) (map[string]interface{}, error) {
	_url := g.review.Host + ":" + strconv.Itoa(g.review.Port) + "/changes/"
	if g.review.User != "" && g.review.Pass != "" {
		_url = g.review.Host + ":" + strconv.Itoa(g.review.Port) + "/a/changes/"
	}

	req, err := http.NewRequest(http.MethodGet, _url, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to request")
	}

	if g.review.User != "" && g.review.Pass != "" {
		req.SetBasicAuth(g.review.User, g.review.Pass)
	}

	q := req.URL.Query()
	q.Add("o", "CURRENT_REVISION")
	q.Add("q", search)
	q.Add("start", strconv.Itoa(start))
	req.URL.RawQuery = q.Encode()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to do")
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("invalid status")
	}

	data, err := ioutil.ReadAll(resp.Body)
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
