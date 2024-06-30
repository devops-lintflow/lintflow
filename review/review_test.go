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
	"io"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"

	"github.com/devops-lintflow/lintflow/config"
	"github.com/devops-lintflow/lintflow/format"
)

const (
	changeGerrit   = 41
	commitGerrit   = "8f71e42dbcd8c68d849e483c04670f58621aab9c"
	queryAfter     = "after:2024-02-27"
	queryBefore    = "before:2024-02-28"
	revisionGerrit = 1
)

func initConfig(name string) (*config.Config, error) {
	c := config.New()
	if c == nil {
		return &config.Config{}, errors.New("failed to new")
	}

	fi, err := os.Open(name)
	if err != nil {
		return c, errors.Wrap(err, "failed to open")
	}

	defer func() {
		_ = fi.Close()
	}()

	buf, err := io.ReadAll(fi)
	if err != nil {
		return c, errors.Wrap(err, "failed to readall")
	}

	if err := yaml.Unmarshal(buf, c); err != nil {
		return c, errors.Wrap(err, "failed to unmarshal")
	}

	return c, nil
}

// nolint: dogsled
func TestReview(t *testing.T) {
	c, err := initConfig("../tests/config.yml")
	assert.Equal(t, nil, err)

	cfg := DefaultConfig()
	cfg.Review = c.Spec.Review

	r := New(cfg)

	d, _ := os.Getwd()
	ti := time.Now()
	root := filepath.Join(d, "gerrit-"+ti.Format("2006-01-02"))

	_, _, _, err = r.Fetch(root, commitGerrit)
	assert.Equal(t, nil, err)

	buf := make([]format.Report, 0)

	vote := config.Vote{
		Label:       "Lint-Verified",
		Approval:    "+1",
		Disapproval: "-1",
		Message:     "Voting Lint-Verified by lintflow",
	}

	err = r.Vote(commitGerrit, buf, vote)
	assert.Equal(t, nil, err)

	buf = make([]format.Report, 1)
	buf[0] = format.Report{
		File:    "AndroidManifest.xml",
		Line:    1,
		Type:    format.TypeError,
		Details: "Disapproved",
	}

	err = r.Vote(commitGerrit, buf, vote)
	assert.Equal(t, nil, err)

	err = r.Clean(root)
	assert.Equal(t, nil, err)
}
