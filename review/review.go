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
	"github.com/craftslab/lintflow/config"
	"github.com/craftslab/lintflow/proto"
)

type Review interface {
	Init() (string, error)
	Run([]proto.Format) error
	Deinit(string) error
}

type Config struct {
	Hash    string
	Name    string
	Reviews []config.Review
}

type review struct {
	cfg *Config
}

func New(cfg *Config) Review {
	return &review{
		cfg: cfg,
	}
}

func DefaultConfig() *Config {
	return &Config{}
}

func (r *review) Init() (string, error) {
	return "", nil
}

func (r *review) Run(data []proto.Format) error {
	return nil
}

func (r *review) Deinit(name string) error {
	return nil
}
