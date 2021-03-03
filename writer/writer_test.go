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

package writer

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/craftslab/lintflow/proto"
)

var (
	fileContent = proto.Format{
		File:    "name",
		Line:    1,
		Type:    proto.TypeError,
		Details: "text",
	}
)

func TestWriterJson(t *testing.T) {
	w := &writer{
		cfg: DefaultConfig(),
	}

	w.cfg.Name = "output.json"
	w.data = make([]proto.Format, 2)
	w.data[0] = fileContent
	w.data[1] = fileContent

	err := w.writeJson()
	defer func(name string) { _ = os.Remove(name) }(w.cfg.Name)

	assert.Equal(t, nil, err)
}

func TestWriteTxt(t *testing.T) {
	w := &writer{
		cfg: DefaultConfig(),
	}

	w.cfg.Name = "output.txt"
	w.data = make([]proto.Format, 2)
	w.data[0] = fileContent
	w.data[1] = fileContent

	err := w.writeTxt()
	defer func(name string) { _ = os.Remove(name) }(w.cfg.Name)

	assert.Equal(t, nil, err)
}

func TestWriteXlsx(t *testing.T) {
	w := &writer{
		cfg: DefaultConfig(),
	}

	w.cfg.Name = "output.xlsx"
	w.data = make([]proto.Format, 2)
	w.data[0] = fileContent
	w.data[1] = fileContent

	err := w.writeXlsx()
	defer func(name string) { _ = os.Remove(name) }(w.cfg.Name)

	assert.Equal(t, nil, err)
}
