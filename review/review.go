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

	"github.com/devops-lintflow/lintflow/config"
	"github.com/devops-lintflow/lintflow/proto"
)

type Review interface {
	Clean(string) error
	Fetch(string, string) (string, string, []string, error)
	Vote(string, []proto.Format, config.Vote) error
}

type Config struct {
	Review config.Review
}

type review struct {
	cfg *Config
	hdl Review
}

func New(cfg *Config) Review {
	return &review{
		cfg: cfg,
		hdl: &gerrit{cfg.Review},
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

func (r *review) Fetch(root, commit string) (dname, rname string, flist []string, emsg error) {
	if r.hdl == nil {
		return "", "", nil, errors.New("invalid handle")
	}

	dir, repo, files, err := r.hdl.Fetch(root, commit)
	if err != nil {
		return "", "", nil, errors.Wrap(err, "failed to fetch")
	}

	return dir, repo, files, nil
}

func (r *review) Vote(commit string, data []proto.Format, vote config.Vote) error {
	if r.hdl == nil {
		return errors.New("invalid handle")
	}

	if err := r.hdl.Vote(commit, data, vote); err != nil {
		return errors.Wrap(err, "failed to vote")
	}

	return nil
}
