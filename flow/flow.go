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
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/craftslab/lintflow/config"
	"github.com/craftslab/lintflow/lint"
	"github.com/craftslab/lintflow/proto"
	"github.com/craftslab/lintflow/review"
	"github.com/craftslab/lintflow/runtime"
)

type Flow interface {
	Run(string) ([]proto.Format, error)
}

type Config struct {
	Config config.Config
	Lint   lint.Lint
	Review review.Review
}

type flow struct {
	cfg    *Config
	filter *config.Filter
}

func New(_ context.Context, cfg *Config) Flow {
	helper := func(c config.Config) *config.Filter {
		var buf config.Filter
		for _, item := range c.Spec.Lint {
			buf.Include.Extension = append(buf.Include.Extension, item.Filter.Include.Extension...)
			buf.Include.File = append(buf.Include.File, item.Filter.Include.File...)
			buf.Include.Repo = append(buf.Include.Repo, item.Filter.Include.Repo...)
		}
		return &buf
	}

	return &flow{
		cfg:    cfg,
		filter: helper(cfg.Config),
	}
}

func DefaultConfig() *Config {
	return &Config{}
}

func (f *flow) Run(commit string) ([]proto.Format, error) {
	var err error
	var ret []proto.Format

	buf, err := runtime.Run(f.routine, []interface{}{commit})
	if err != nil {
		return nil, errors.Wrap(err, "failed to run")
	}

	err = nil

	for _, val := range buf {
		if val == nil {
			err = errors.New("invalid data")
			break
		}
		if len(val.([]proto.Format)) != 0 {
			ret = append(ret, val.([]proto.Format)...)
		}
	}

	return ret, err
}

func (f *flow) routine(data interface{}) interface{} {
	d, _ := os.Getwd()
	t := time.Now()
	root := filepath.Join(d, "gerrit-"+t.Format("2006-01-02"))

	commit := data.(string)

	dir, files, err := f.cfg.Review.Fetch(root, commit, f.match)
	defer func() { _ = f.cfg.Review.Clean(root) }()
	if err != nil {
		log.Println(err)
		return nil
	}

	buf, err := f.cfg.Lint.Run(dir, files, f.match)
	if err != nil {
		log.Println(err)
		return nil
	}

	if buf == nil {
		return []proto.Format{}
	}

	if err := f.cfg.Review.Vote(commit, buf); err != nil {
		log.Println(err)
		return nil
	}

	return buf
}

func (f *flow) match(filter *config.Filter, repo, file string) bool {
	matchExtension := func(filter *config.Filter, data string) bool {
		for _, val := range filter.Include.Extension {
			if val == filepath.Ext(strings.TrimSuffix(data, proto.Base64Content)) {
				return true
			}
		}
		return false
	}

	matchFile := func(filter *config.Filter, data string) bool {
		for _, val := range filter.Include.File {
			if val == filepath.Base(strings.TrimSuffix(data, proto.Base64Content)) {
				return true
			}
		}
		return false
	}

	matchRepo := func(filter *config.Filter, data string) bool {
		if len(filter.Include.Repo) == 0 {
			return true
		}
		for _, val := range filter.Include.Repo {
			if val == data {
				return true
			}
		}
		return false
	}

	if filter == nil {
		filter = f.filter
	}

	if repo != "" && !matchRepo(filter, repo) {
		return false
	}

	if !matchExtension(filter, file) && !matchFile(filter, file) {
		return false
	}

	return true
}
