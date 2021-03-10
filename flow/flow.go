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

	"github.com/pkg/errors"

	"github.com/craftslab/lintflow/lint"
	"github.com/craftslab/lintflow/proto"
	"github.com/craftslab/lintflow/review"
	"github.com/craftslab/lintflow/runtime"
)

type Flow interface {
	Run(string) ([]proto.Format, error)
}

type Config struct {
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

func (f *flow) Run(commit string) ([]proto.Format, error) {
	buf, err := runtime.Run(f.routine, []interface{}{commit})
	if err != nil {
		return nil, errors.Wrap(err, "failed to run")
	}

	ret := make([]proto.Format, len(buf))
	for index, val := range buf {
		if val == nil {
			return nil, errors.New("invalid data")
		}
		ret[index] = val.(proto.Format)
	}

	return ret, nil
}

func (f *flow) routine(data interface{}) interface{} {
	commit := data.(string)

	root, files, err := f.cfg.Review.Fetch(commit)
	defer func() { _ = f.cfg.Review.Clean(root) }()
	if err != nil {
		return nil
	}

	buf, err := f.cfg.Lint.Run(root, files)
	if err != nil {
		return nil
	}

	if err := f.cfg.Review.Vote(commit, buf); err != nil {
		return nil
	}

	return buf
}
