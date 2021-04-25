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
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/craftslab/lintflow/config"
)

// nolint: funlen
// nolint: goconst
func TestFilter(t *testing.T) {
	filter := config.Filter{
		Include: config.Include{
			Extension: []string{".java", ".xml"},
			File:      []string{"message"},
			Repo:      []string{""},
		},
	}

	f := flow{nil, &filter}

	ret := f.match(&filter, "", ".ext")
	assert.Equal(t, false, ret)

	ret = f.match(&filter, "", "foo.ext")
	assert.Equal(t, false, ret)

	ret = f.match(&filter, "", ".java")
	assert.Equal(t, true, ret)

	ret = f.match(nil, "", ".java")
	assert.Equal(t, true, ret)

	ret = f.match(&filter, "", "foo.java")
	assert.Equal(t, true, ret)

	ret = f.match(nil, "", "foo.java")
	assert.Equal(t, true, ret)

	ret = f.match(&filter, "", "message")
	assert.Equal(t, true, ret)

	ret = f.match(nil, "", "message")
	assert.Equal(t, true, ret)

	filter.Include.Repo = []string{"alpha", "beta"}

	ret = f.match(&filter, "", "message")
	assert.Equal(t, true, ret)

	ret = f.match(nil, "", "message")
	assert.Equal(t, true, ret)

	ret = f.match(&filter, "foo", "message")
	assert.Equal(t, false, ret)

	ret = f.match(&filter, "alpha", "message")
	assert.Equal(t, true, ret)

	ret = f.match(nil, "alpha", "message")
	assert.Equal(t, true, ret)
}
