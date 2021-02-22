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
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/craftslab/lintflow/config"
	"github.com/craftslab/lintflow/proto"
)

const (
	HASH_GERRIT = "0d43a7ec2377edeacc016878595229ca377c80f1"
)

const (
	NAME_GERRIT  = "gerrit"
	NAME_INVALID = "invalid"
)

func TestReview(t *testing.T) {
	r := &review{
		cfg: DefaultConfig(),
	}

	vote := make([]config.Vote, 1)
	vote[0] = config.Vote{
		Approval:    "+1",
		Disapproval: "-1",
		Label:       "Code-Review",
	}

	r.cfg.Reviews = make([]config.Review, 1)
	r.cfg.Reviews[0] = config.Review{
		Host: "127.0.0.1",
		Name: NAME_GERRIT,
		Pass: "D/uccEPCcItsY3Cti4unrkS/zsyW65MZBrEsiHiXpg",
		Port: 8080,
		User: "admin",
		Vote: vote,
	}

	r.cfg.Hash = ""
	r.cfg.Name = NAME_INVALID

	_, err := r.Init()
	assert.NotEqual(t, nil, err)

	r.cfg.Hash = HASH_GERRIT
	r.cfg.Name = NAME_INVALID

	_, err = r.Init()
	assert.NotEqual(t, nil, err)

	r.cfg.Hash = HASH_GERRIT
	r.cfg.Name = NAME_GERRIT

	project, err := r.Init()
	assert.Equal(t, nil, err)

	buf := make([]proto.Format, 0)

	err = r.Run(buf)
	assert.Equal(t, nil, err)

	buf = make([]proto.Format, 1)
	buf[0] = proto.Format{
		Details: "Disapproved",
		File:    "AndroidManifest.xml",
		Line:    1,
		Type:    proto.TYPE_ERROR,
	}

	err = r.Run(buf)
	assert.Equal(t, nil, err)

	err = r.Deinit(project)
	assert.Equal(t, nil, err)
}
