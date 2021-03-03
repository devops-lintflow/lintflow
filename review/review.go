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
	"github.com/pkg/errors"

	"github.com/craftslab/lintflow/config"
	"github.com/craftslab/lintflow/proto"
)

type Review interface {
	Clean(string) error
	Fetch(string) (string, error)
	Vote(string, []proto.Format) error
}

type Config struct {
	Name    string
	Reviews []config.Review
}

type review struct {
	cfg *Config
	hdl Review
}

func New(cfg *Config) Review {
	reviews := map[string]Review{}
	for _, item := range cfg.Reviews {
		reviews[item.Name] = &gerrit{item}
	}

	h, p := reviews[cfg.Name]
	if !p {
		h = nil
	}

	return &review{
		cfg: cfg,
		hdl: h,
	}
}

func DefaultConfig() *Config {
	return &Config{}
}

func (r *review) Clean(name string) error {
	if r.hdl == nil {
		return errors.New("invalid handle")
	}

	if err := r.hdl.Clean(name); err != nil {
		return errors.Wrap(err, "failed to clean")
	}

	return nil
}

func (r *review) Fetch(commit string) (string, error) {
	if r.hdl == nil {
		return "", errors.New("invalid handle")
	}

	buf, err := r.hdl.Fetch(commit)
	if err != nil {
		return "", errors.Wrap(err, "failed to fetch")
	}

	return buf, nil
}

func (r *review) Vote(commit string, data []proto.Format) error {
	if r.hdl == nil {
		return errors.New("invalid handle")
	}

	if err := r.hdl.Vote(commit, data); err != nil {
		return errors.Wrap(err, "failed to vote")
	}

	return nil
}
