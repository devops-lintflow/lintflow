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

package runtime

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/goleak"
)

func TestRunRuntime(t *testing.T) {
	defer goleak.VerifyNone(t)

	op := func(ctx context.Context, req interface{}) interface{} {
		return nil
	}

	var req []interface{}

	_, err := Run(context.Background(), op, req)
	assert.Equal(t, nil, err)
}
