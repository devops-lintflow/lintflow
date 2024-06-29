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
	"io"
	"os"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"

	"github.com/devops-lintflow/lintflow/config"
	"github.com/devops-lintflow/lintflow/proto"
)

func initConfig(name string) (*config.Config, error) {
	c := config.New()
	if c == nil {
		return &config.Config{}, errors.New("failed to new")
	}

	fi, err := os.Open(name)
	if err != nil {
		return c, errors.Wrap(err, "failed to open")
	}

	defer func() {
		_ = fi.Close()
	}()

	buf, err := io.ReadAll(fi)
	if err != nil {
		return c, errors.Wrap(err, "failed to readall")
	}

	if err := yaml.Unmarshal(buf, c); err != nil {
		return c, errors.Wrap(err, "failed to unmarshal")
	}

	return c, nil
}

// nolint: funlen
// nolint: goconst
func TestMatchFilter(t *testing.T) {
	filter := config.Filter{
		Include: config.Include{
			Extensions: []string{".java", ".xml"},
			Files:      []string{"message"},
			Repos:      []string{""},
		},
	}

	f := flow{}

	ret := f.matchFilter(&filter, "", ".ext")
	assert.Equal(t, false, ret)

	ret = f.matchFilter(&filter, "", "foo.ext")
	assert.Equal(t, false, ret)

	ret = f.matchFilter(&filter, "", ".java")
	assert.Equal(t, true, ret)

	ret = f.matchFilter(nil, "", ".java")
	assert.Equal(t, false, ret)

	ret = f.matchFilter(&filter, "", "foo.java")
	assert.Equal(t, true, ret)

	ret = f.matchFilter(nil, "", "foo.java")
	assert.Equal(t, false, ret)

	ret = f.matchFilter(&filter, "", "message")
	assert.Equal(t, true, ret)

	ret = f.matchFilter(nil, "", "message")
	assert.Equal(t, false, ret)

	filter.Include.Repos = []string{"alpha", "beta"}

	ret = f.matchFilter(&filter, "", "message")
	assert.Equal(t, true, ret)

	ret = f.matchFilter(nil, "", "message")
	assert.Equal(t, false, ret)

	ret = f.matchFilter(&filter, "foo", "message")
	assert.Equal(t, false, ret)

	ret = f.matchFilter(&filter, "alpha", "message")
	assert.Equal(t, true, ret)

	ret = f.matchFilter(nil, "alpha", "message")
	assert.Equal(t, false, ret)
}

func TestBuildLabel(t *testing.T) {
	c, err := initConfig("../tests/config.yml")
	assert.Equal(t, nil, err)

	cfg := DefaultConfig()
	cfg.Config = *c

	f := flow{
		cfg: cfg,
	}

	buf := map[string][]proto.Format{
		"lintai": {
			{
				File:    "/path/to/file1",
				Line:    1,
				Type:    proto.TypeError,
				Details: "Disapproved",
			},
		},
		"lintcpp": {
			{
				File:    "/path/to/file2",
				Line:    1,
				Type:    proto.TypeWarn,
				Details: "Disapproved",
			},
		},
	}

	ret := f.buildLabel(buf)
	assert.NotEqual(t, nil, ret)
	assert.NotEqual(t, 0, len(ret))
}

func TestBuildVote(t *testing.T) {
	c, err := initConfig("../tests/config.yml")
	assert.Equal(t, nil, err)

	cfg := DefaultConfig()
	cfg.Config = *c

	f := flow{
		cfg: cfg,
	}

	ret := f.buildVote("Lint-Verified")
	assert.Equal(t, "Lint-Verified", ret.Label)
}
