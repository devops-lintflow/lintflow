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
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/craftslab/lintflow/config"
	"github.com/craftslab/lintflow/proto"
)

const (
	changeGerrit   = 21
	commitGerrit   = "6c7447ffad57d13f7f2657a302dd176cff453c80"
	revisionGerrit = 1
)

const (
	NAME_GERRIT  = "gerrit"
	NAME_INVALID = "invalid"
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

	buf, err := ioutil.ReadAll(fi)
	if err != nil {
		return c, errors.Wrap(err, "failed to readall")
	}

	if err := yaml.Unmarshal(buf, c); err != nil {
		return c, errors.Wrap(err, "failed to unmarshal")
	}

	return c, nil
}

func TestReview(t *testing.T) {
	c, err := initConfig("../tests/config.yml")
	assert.Equal(t, nil, err)

	cfg := DefaultConfig()
	cfg.Name = NAME_GERRIT
	cfg.Reviews = c.Spec.Review

	r := New(cfg)

	_, err = r.Fetch("")
	assert.NotEqual(t, nil, err)

	name, err := r.Fetch(commitGerrit)
	assert.Equal(t, nil, err)

	buf := make([]proto.Format, 0)

	err = r.Vote(commitGerrit, buf)
	assert.Equal(t, nil, err)

	buf = make([]proto.Format, 1)
	buf[0] = proto.Format{
		Details: "Disapproved",
		File:    "AndroidManifest.xml",
		Line:    1,
		Type:    proto.TypeError,
	}

	err = r.Vote(commitGerrit, buf)
	assert.Equal(t, nil, err)

	err = r.Clean(name)
	assert.Equal(t, nil, err)
}
