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

package cmd

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestInitConfig(t *testing.T) {
	var err error

	_, err = initConfig("invalid.yml")
	assert.NotEqual(t, nil, err)

	_, err = initConfig("../tests/invalid.yml")
	assert.NotEqual(t, nil, err)

	_, err = initConfig("../tests/config.yml")
	assert.Equal(t, nil, err)
}

func TestInitReview(t *testing.T) {
	c, err := initConfig("../tests/config.yml")
	assert.Equal(t, nil, err)

	_, err = initReview(c)
	assert.Equal(t, nil, err)
}

func TestInitLint(t *testing.T) {
	c, err := initConfig("../tests/config.yml")
	assert.Equal(t, nil, err)

	_, err = initLint(c)
	assert.Equal(t, nil, err)
}

func TestSetTimeout(t *testing.T) {
	timeout, err := setTimeout("")
	assert.Equal(t, nil, err)
	assert.Equal(t, Timeout, timeout)

	timeout, err = setTimeout("1s")
	assert.Equal(t, nil, err)
	assert.Equal(t, 1*time.Second, timeout)

	timeout, err = setTimeout("10m")
	assert.Equal(t, nil, err)
	assert.Equal(t, 10*time.Minute, timeout)

	timeout, err = setTimeout("100h")
	assert.Equal(t, nil, err)
	assert.Equal(t, 100*time.Hour, timeout)
}
