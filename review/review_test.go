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
	"fmt"
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
	changeGerrit   = 42
	commitGerrit   = "aed052e8b66810795d3a894a9095db41e5854b70"
	queryAfter     = "after:2024-08-01"
	queryBefore    = "before:2024-09-01"
	revisionGerrit = 72
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

	dir, repo, files, meta, patch, err := r.Fetch(root, commitGerrit)
	assert.Equal(t, nil, err)

	fmt.Printf("  dir: %s\n", dir)
	fmt.Printf(" repo: %s\n", repo)
	fmt.Printf("files: %v\n", files)
	fmt.Printf("meta: %s\n", meta)
	fmt.Printf("patch: %s\n", patch)

	buf := make([]format.Report, 1)
	buf[0] = format.Report{
		File:    "lintshell/test.sh",
		Line:    1,
		Type:    format.TypeError,
		Details: "Disapproved by review",
	}

	vote := config.Vote{
		Label:       "Lint-Verified",
		Approval:    "+1",
		Disapproval: "-1",
		Message:     "Voting Lint-Verified by review",
	}

	err = r.Vote(commitGerrit, buf, vote)
	assert.Equal(t, nil, err)

	err = r.Clean(root)
	assert.Equal(t, nil, err)
}
