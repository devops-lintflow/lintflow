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
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	root = "../tests/gerrit-2021-03-06/21/c5d3440911e06ed4fc60252bd89e7756f9ae67ee"
)

func TestMarshal(t *testing.T) {
	var buf []string
	var l lint

	_, err := l.marshal(root, buf)
	assert.NotEqual(t, nil, err)

	buf = []string{"foo.base64"}
	_, err = l.marshal(root, buf)
	assert.NotEqual(t, nil, err)
}
