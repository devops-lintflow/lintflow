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

package printer

import (
	"github.com/craftslab/lintflow/proto"
)

type Printer interface {
	Run([]proto.Format) error
}

type Config struct {
	Name string
}

type printer struct {
	config *Config
}

func New(config *Config) Printer {
	return &printer{
		config: config,
	}
}

func DefaultConfig() *Config {
	return &Config{}
}

func (p *printer) Run([]proto.Format) error {
	return nil
}
