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
	name, err := f.cfg.Review.Fetch(commit)
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch")
	}

	bufLint, err := f.cfg.Lint.Run(name)
	if err != nil {
		return nil, errors.Wrap(err, "failed to lint")
	}

	bufRuntime := make([]interface{}, len(bufLint))
	for index, val := range bufLint {
		bufRuntime[index] = val
	}

	retRuntime, err := runtime.Run(f.routine, bufRuntime)
	if err != nil {
		return nil, errors.Wrap(err, "failed to run")
	}

	bufReview := make([]proto.Format, len(retRuntime))
	for index, val := range retRuntime {
		bufReview[index] = val.(proto.Format)
	}

	if err := f.cfg.Review.Vote(commit, bufReview); err != nil {
		return nil, errors.Wrap(err, "failed to vote")
	}

	if err := f.cfg.Review.Clean(name); err != nil {
		return nil, errors.Wrap(err, "failed to clean")
	}

	return bufReview, nil
}

func (f *flow) routine(data interface{}) interface{} {
	// TODO
	return nil
}
