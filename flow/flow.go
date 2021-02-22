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

	"github.com/craftslab/lintflow/lint"
	"github.com/craftslab/lintflow/proto"
	"github.com/craftslab/lintflow/review"
)

type Flow interface {
	Run() ([]proto.Format, error)
}

type Config struct {
	Lint   lint.Lint
	Review review.Review
}

type flow struct {
	config *Config
}

func New(_ context.Context, config *Config) Flow {
	return &flow{
		config: config,
	}
}

func DefaultConfig() *Config {
	return &Config{}
}

func (f *flow) Run() ([]proto.Format, error) {
	return nil, nil
}
