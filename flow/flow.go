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

package flow

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/devops-lintflow/lintflow/config"
	"github.com/devops-lintflow/lintflow/format"
	"github.com/devops-lintflow/lintflow/lint"
	"github.com/devops-lintflow/lintflow/review"
)

type Flow interface {
	Run(context.Context, string) error
}

type Config struct {
	Config config.Config
	Lint   lint.Lint
	Review review.Review
}

type flow struct {
	cfg *Config
}

func New(_ context.Context, cfg *Config) Flow {
	return &flow{
		cfg: cfg,
	}
}

func DefaultConfig() *Config {
	return &Config{}
}

func (f *flow) Run(ctx context.Context, commit string) error {
	d, _ := os.Getwd()
	t := time.Now()

	root := filepath.Join(d, "gerrit-"+t.Format("2006-01-02"))

	dir, repo, files, err := f.cfg.Review.Fetch(root, commit)

	defer func() {
		_ = f.cfg.Review.Clean(root)
	}()

	if err != nil {
		return errors.Wrap(err, "failed to clean reivew")
	}

	buf, err := f.cfg.Lint.Run(ctx, dir, repo, files, f.matchFilter)
	if err != nil {
		return errors.Wrap(err, "failed to run lint")
	}

	if buf == nil || len(buf) == 0 {
		return nil
	}

	labels := f.buildLabel(buf)

	for label, data := range labels {
		if vote := f.buildVote(label); vote.Label != "" {
			if err := f.cfg.Review.Vote(commit, data, vote); err != nil {
				return errors.Wrap(err, "failed to vote reivew")
			}
		}
	}

	return nil
}

func (f *flow) matchFilter(filter *config.Filter, repo, file string) bool {
	matchExtension := func(filter *config.Filter, data string) bool {
		for _, val := range filter.Include.Extensions {
			if val == filepath.Ext(strings.TrimSuffix(data, format.Base64Content)) {
				return true
			}
		}
		return false
	}

	matchFile := func(filter *config.Filter, data string) bool {
		for _, val := range filter.Include.Files {
			if val == filepath.Base(strings.TrimSuffix(data, format.Base64Content)) {
				return true
			}
		}
		return false
	}

	matchRepo := func(filter *config.Filter, data string) bool {
		if len(filter.Include.Repos) == 0 {
			return true
		}
		for _, val := range filter.Include.Repos {
			if val == data {
				return true
			}
		}
		return false
	}

	if filter == nil {
		return false
	}

	if repo != "" && !matchRepo(filter, repo) {
		return false
	}

	if !matchExtension(filter, file) && !matchFile(filter, file) {
		return false
	}

	return true
}

func (f *flow) buildLabel(data map[string][]format.Format) map[string][]format.Format {
	helper := func(name string) string {
		var buf string
		lints := f.cfg.Config.Spec.Lints
		for i := range lints {
			if lints[i].Name == name {
				buf = lints[i].Vote
				break
			}
		}
		return buf
	}

	buf := map[string][]format.Format{}

	for key, val := range data {
		if vote := helper(key); vote != "" {
			if _, ok := buf[vote]; !ok {
				buf[vote] = val
			} else {
				buf[vote] = append(buf[vote], val...)
			}
		}
	}

	return buf
}

func (f *flow) buildVote(label string) config.Vote {
	var buf config.Vote

	votes := f.cfg.Config.Spec.Review.Votes

	for _, item := range votes {
		if item.Label == label {
			buf = item
			break
		}
	}

	return buf
}
