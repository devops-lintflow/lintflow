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

package lint

import (
	"github.com/pkg/errors"

	"github.com/craftslab/lintflow/config"
	"github.com/craftslab/lintflow/proto"
)

type Lint interface {
	Run(root string, files []string) ([]proto.Format, error)
}

type Config struct {
	Lints []config.Lint
}

type lint struct {
	cfg *Config
}

func New(cfg *Config) Lint {
	return &lint{
		cfg: cfg,
	}
}

func DefaultConfig() *Config {
	return &Config{}
}

func (l *lint) Run(root string, files []string) ([]proto.Format, error) {
	var ret []proto.Format

	// TODO: goroutine
	for _, val := range l.cfg.Lints {
		buf, err := l.filter(val.Name, files)
		if err != nil {
			return nil, errors.Wrap(err, "failed to filter")
		}
		b, err := l.marshal(root, buf)
		if err != nil {
			return nil, errors.Wrap(err, "failed to marshal")
		}
		ret, err = l.routine(val.Host, val.Port, b)
		if err != nil {
			return nil, errors.Wrap(err, "failed to routine")
		}
	}

	return ret, nil
}

func (l *lint) filter(name string, data []string) ([]string, error) {
	// TODO
	return nil, nil
}

func (l *lint) marshal(root string, data []string) ([]byte, error) {
	// TODO
	return nil, nil
}

func (l *lint) routine(host string, port int, data []byte) ([]proto.Format, error) {
	// TODO
	return nil, nil
}
