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
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/reviewdog/reviewdog/diff"

	"github.com/devops-lintflow/lintflow/config"
	"github.com/devops-lintflow/lintflow/format"
)

const (
	commitMsg   = "/COMMIT_MSG"
	commitQuery = "commit"
)

const (
	diffBin    = "Binary files differ"
	diffSep    = "diff --git"
	diffPrefix = "b/"
)

const (
	metaBranch   = "branch"
	metaOwner    = "owner"
	metaProject  = "project"
	metaRevision = "revision"
	metaUpdated  = "updated"
	metaUrl      = "url"
)

const (
	suffixMeta  = "meta"
	suffixPatch = "patch"
)

const (
	queryLimit   = 1000
	urlChanges   = "/changes/"
	urlContent   = "/content"
	urlDetail    = "/detail"
	urlFiles     = "/files/"
	urlNumber    = "&n="
	urlOption    = "&o="
	urlPatch     = "/patch"
	urlPrefix    = "/a"
	urlQuery     = "?q="
	urlReview    = "/review"
	urlRevisions = "/revisions/"
	urlStart     = "&start="
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

// nolint:funlen,gocritic,gocyclo
func (g *gerrit) Fetch(root, commit string) (dname, rname string, flist []string, mname, pname string, emsg error) {
	filterFiles := func(data map[string]interface{}) map[string]interface{} {
		buf := make(map[string]interface{})
		for key, val := range data {
			if v, ok := val.(map[string]interface{})["status"]; ok {
				if v.(string) == "D" || v.(string) == "R" {
					continue
				}
			}
			buf[key] = val
		}
		return buf
	}

	// Query commit
	buf, err := g.get(g.urlQuery(commitQuery+":"+commit, []string{"ALL_COMMITS", "ALL_REVISIONS", "DETAILED_ACCOUNTS"}, 0))
	if err != nil {
		return "", "", nil, "", "", errors.Wrap(err, "failed to query")
	}

	queryRet, err := g.unmarshalList(buf)
	if err != nil {
		return "", "", nil, "", "", errors.Wrap(err, "failed to unmarshalList")
	}

	changeNum := int(queryRet[0].(map[string]interface{})["_number"].(float64))

	revisions := queryRet[0].(map[string]interface{})["revisions"].(map[string]interface{})
	current := revisions[commit].(map[string]interface{})
	revisionNum := int(current["_number"].(float64))

	path := filepath.Join(root, strconv.Itoa(changeNum), commit)

	// Get files
	buf, err = g.get(g.urlFiles(changeNum, revisionNum))
	if err != nil {
		return "", "", nil, "", "", errors.Wrap(err, "failed to get files")
	}

	fs, err := g.unmarshal(buf)
	if err != nil {
		return "", "", nil, "", "", errors.Wrap(err, "failed to unmarshal")
	}

	// Match files
	fs = filterFiles(fs)

	// Get content
	for key := range fs {
		buf, err = g.get(g.urlContent(changeNum, revisionNum, key))
		if err != nil {
			return "", "", nil, "", "", errors.Wrap(err, "failed to get content")
		}

		file := filepath.Base(key)
		if key == commitMsg {
			file = strings.TrimPrefix(commitMsg, "/")
		}

		err = g.write(filepath.Join(path, filepath.Dir(key)), file, string(buf))
		if err != nil {
			return "", "", nil, "", "", errors.Wrap(err, "failed to write content")
		}
	}

	var files []string

	for key := range fs {
		if key == commitMsg {
			files = append(files, strings.TrimPrefix(commitMsg, "/"))
		} else {
			files = append(files, filepath.Join(filepath.Dir(key), filepath.Base(key)))
		}
	}

	// Get meta
	buf, err = g.meta(commit, queryRet[0])
	if err != nil {
		return "", "", nil, "", "", errors.Wrap(err, "failed to get meta")
	}

	meta := fmt.Sprintf("%d-%s.%s", changeNum, commit[:7], suffixMeta)

	err = g.write(path, meta, string(buf))
	if err != nil {
		return "", "", nil, "", "", errors.Wrap(err, "failed to write meta")
	}

	// Get patch
	buf, err = g.get(g.urlPatch(changeNum, revisionNum))
	if err != nil {
		return "", "", nil, "", "", errors.Wrap(err, "failed to get patch")
	}

	patch := fmt.Sprintf("%d-%s.%s", changeNum, commit[:7], suffixPatch)

	err = g.write(path, patch, string(buf))
	if err != nil {
		return "", "", nil, "", "", errors.Wrap(err, "failed to write patch")
	}

	return path, queryRet[0].(map[string]interface{})["project"].(string), files, meta, patch, nil
}

func (g *gerrit) Query(search string, start int) ([]interface{}, error) {
	helper := func(search string, start int) []interface{} {
		buf, err := g.get(g.urlQuery(search, []string{"ALL_REVISIONS", "DETAILED_ACCOUNTS"}, start))
		if err != nil {
			return nil
		}
		ret, err := g.unmarshalList(buf)
		if err != nil {
			return nil
		}
		return ret
	}

	buf := helper(search, start)
	if len(buf) == 0 {
		return []interface{}{}, nil
	}

	more, ok := buf[len(buf)-1].(map[string]interface{})["_more_changes"].(bool)
	if !ok {
		more = false
	}

	if !more {
		return buf, nil
	}

	if b, err := g.Query(search, start+len(buf)); err == nil {
		buf = append(buf, b...)
	}

	return buf, nil
}

// nolint:funlen,gocyclo
func (g *gerrit) Vote(commit string, data []format.Report, vote config.Vote) error {
	match := func(data format.Report, diffs []*diff.FileDiff) bool {
		for _, d := range diffs {
			if strings.Replace(d.PathNew, diffPrefix, "", 1) != data.File {
				continue
			}
			if data.Line <= 0 {
				return true
			}
			for _, h := range d.Hunks {
				for _, l := range h.Lines {
					if l.Type == diff.LineAdded && l.LnumNew == data.Line {
						return true
					}
				}
			}
		}
		return false
	}

	build := func(data []format.Report, diffs []*diff.FileDiff) (map[string]interface{}, map[string]interface{}, string) {
		if len(data) == 0 {
			return nil, map[string]interface{}{vote.Label: vote.Approval}, vote.Message
		}
		c := map[string]interface{}{}
		for _, item := range data {
			if item.Details == "" || (item.File != commitMsg && !match(item, diffs)) {
				continue
			}
			l := item.Line
			if l <= 0 {
				l = 1
			}
			b := map[string]interface{}{"line": l, "message": item.Details, "unresolved": false}
			if _, ok := c[item.File]; !ok {
				c[item.File] = []map[string]interface{}{b}
			} else {
				c[item.File] = append(c[item.File].([]map[string]interface{}), b)
			}
		}
		if len(c) == 0 {
			return nil, map[string]interface{}{vote.Label: vote.Approval}, vote.Message
		} else {
			return c, map[string]interface{}{vote.Label: vote.Disapproval}, vote.Message
		}
	}

	// Query commit
	ret, err := g.get(g.urlQuery(commitQuery+":"+commit, []string{"ALL_REVISIONS", "DETAILED_ACCOUNTS"}, 0))
	if err != nil {
		return errors.Wrap(err, "failed to query")
	}

	c, err := g.unmarshalList(ret)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshalList")
	}

	revisions := c[0].(map[string]interface{})["revisions"].(map[string]interface{})
	current := revisions[commit].(map[string]interface{})
	revisionNum := int(current["_number"].(float64))

	// Get patch
	ret, err = g.get(g.urlPatch(int(c[0].(map[string]interface{})["_number"].(float64)), revisionNum))
	if err != nil {
		return errors.Wrap(err, "failed to patch")
	}

	// Parse diff
	dec := make([]byte, base64.StdEncoding.DecodedLen(len(ret)))
	if _, err = base64.StdEncoding.Decode(dec, ret); err != nil {
		return errors.Wrap(err, "failed to decode")
	}

	index := bytes.Index(dec, []byte(diffSep))
	if index < 0 {
		return errors.New("failed to index")
	}

	var b []byte

	for _, item := range bytes.SplitAfter(dec[index:], []byte(diffSep)) {
		if !bytes.Contains(item, []byte(diffBin)) {
			b = bytes.Join([][]byte{b, item}, []byte(""))
		}
	}

	diffs, err := diff.ParseMultiFile(bytes.NewReader(b))
	if err != nil {
		return errors.Wrap(err, "failed to parse")
	}

	// Review commit
	comments, labels, message := build(data, diffs)
	fmt.Printf("  labels: %v\n", labels)
	fmt.Printf(" message: %s\n", message)
	buf := map[string]interface{}{"comments": comments, "labels": labels, "message": message}
	if err := g.post(g.urlReview(int(c[0].(map[string]interface{})["_number"].(float64)), revisionNum), buf); err != nil {
		return errors.Wrap(err, "failed to review")
	}

	return nil
}

func (g *gerrit) write(dir, file, data string) error {
	_ = os.MkdirAll(dir, os.ModePerm)

	f, err := os.Create(filepath.Join(dir, file))
	if err != nil {
		return errors.Wrap(err, "failed to create")
	}
	defer func() { _ = f.Close() }()

	w := bufio.NewWriter(f)
	if _, err := w.WriteString(data); err != nil {
		return errors.Wrap(err, "failed to write")
	}
	defer func() { _ = w.Flush() }()

	return nil
}

func (g *gerrit) unmarshal(data []byte) (map[string]interface{}, error) {
	buf := map[string]interface{}{}

	if err := json.Unmarshal(data[4:], &buf); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal")
	}

	return buf, nil
}

func (g *gerrit) unmarshalList(data []byte) ([]interface{}, error) {
	var buf []interface{}

	if err := json.Unmarshal(data[4:], &buf); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal")
	}

	if len(buf) == 0 {
		return nil, errors.New("failed to match")
	}

	return buf, nil
}

func (g *gerrit) urlContent(change, revision int, name string) string {
	buf := g.r.Url + urlChanges + strconv.Itoa(change) +
		urlRevisions + strconv.Itoa(revision) + urlFiles + url.QueryEscape(name) + urlContent

	if g.r.User != "" && g.r.Pass != "" {
		buf = g.r.Url + urlPrefix + urlChanges + strconv.Itoa(change) +
			urlRevisions + strconv.Itoa(revision) + urlFiles + url.QueryEscape(name) + urlContent
	}

	return buf
}

func (g *gerrit) urlDetail(change int) string {
	buf := g.r.Url + urlChanges + strconv.Itoa(change) + urlDetail

	if g.r.User != "" && g.r.Pass != "" {
		buf = g.r.Url + urlPrefix + urlChanges + strconv.Itoa(change) + urlDetail
	}

	return buf
}

func (g *gerrit) urlFiles(change, revision int) string {
	buf := g.r.Url + urlChanges + strconv.Itoa(change) +
		urlRevisions + strconv.Itoa(revision) + urlFiles

	if g.r.User != "" && g.r.Pass != "" {
		buf = g.r.Url + urlPrefix + urlChanges + strconv.Itoa(change) +
			urlRevisions + strconv.Itoa(revision) + urlFiles
	}

	return buf
}

func (g *gerrit) urlPatch(change, revision int) string {
	buf := g.r.Url + urlChanges + strconv.Itoa(change) +
		urlRevisions + strconv.Itoa(revision) + urlPatch

	if g.r.User != "" && g.r.Pass != "" {
		buf = g.r.Url + urlPrefix + urlChanges + strconv.Itoa(change) +
			urlRevisions + strconv.Itoa(revision) + urlPatch
	}

	return buf
}

func (g *gerrit) urlQuery(search string, option []string, start int) string {
	query := urlQuery + url.PathEscape(search) +
		urlOption + strings.Join(option, urlOption) +
		urlStart + strconv.Itoa(start) +
		urlNumber + strconv.Itoa(queryLimit)

	buf := g.r.Url + urlChanges + query
	if g.r.User != "" && g.r.Pass != "" {
		buf = g.r.Url + urlPrefix + urlChanges + query
	}

	return buf
}

func (g *gerrit) urlReview(change, revision int) string {
	buf := g.r.Url + urlChanges + strconv.Itoa(change) +
		urlRevisions + strconv.Itoa(revision) + urlReview

	if g.r.User != "" && g.r.Pass != "" {
		buf = g.r.Url + urlPrefix + urlChanges + strconv.Itoa(change) +
			urlRevisions + strconv.Itoa(revision) + urlReview
	}

	return buf
}

func (g *gerrit) meta(rev string, _query interface{}) ([]byte, error) {
	helper := func(offset int) string {
		sign := "+"
		if offset < 0 {
			sign = "-"
			offset = -offset
		}
		hours := offset / 60
		minutes := offset % 60
		return fmt.Sprintf("%s%02d:%02d", sign, hours, minutes)
	}

	owner := _query.(map[string]interface{})["owner"]
	updated := _query.(map[string]interface{})["updated"].(string)

	revisions := _query.(map[string]interface{})["revisions"]
	revision := revisions.(map[string]interface{})[rev]
	commit := revision.(map[string]interface{})["commit"]
	committer := commit.(map[string]interface{})["committer"]
	tz := committer.(map[string]interface{})["tz"].(float64)

	date, _ := time.Parse(time.DateTime, updated)
	updated = fmt.Sprintf("%d-%d-%dT%d:%d:%d%s", date.Year(), date.Month(), date.Day(),
		date.Hour(), date.Minute(), date.Second(), helper(int(tz)))

	buf := map[string]string{
		metaBranch:   _query.(map[string]interface{})["branch"].(string),
		metaOwner:    owner.(map[string]interface{})["name"].(string),
		metaProject:  _query.(map[string]interface{})["project"].(string),
		metaRevision: rev,
		metaUpdated:  updated,
		metaUrl:      g.r.Url,
	}

	ret, err := json.Marshal(buf)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal")
	}

	dst := make([]byte, base64.StdEncoding.EncodedLen(len(ret)))
	base64.StdEncoding.Encode(dst, ret)

	return dst, nil
}

func (g *gerrit) get(_url string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, _url, http.NoBody)
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

	data, err := io.ReadAll(rsp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read")
	}

	return data, nil
}

func (g *gerrit) post(_url string, data map[string]interface{}) error {
	buf, err := json.Marshal(data)
	if err != nil {
		return errors.Wrap(err, "failed to marshal")
	}

	req, err := http.NewRequest(http.MethodPost, _url, bytes.NewBuffer(buf))
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

	_, err = io.ReadAll(rsp.Body)
	if err != nil {
		return errors.Wrap(err, "failed to read")
	}

	return nil
}
