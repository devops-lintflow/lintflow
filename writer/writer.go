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
	"bufio"
	"encoding/json"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/devops-lintflow/lintflow/proto"
)

const (
	base = 10
	perm = 0600
	sep  = ","
)

var (
	sheetName = time.Now().Local().Format(time.RFC3339)
)

type Writer interface {
	Run(string, []proto.Format) error
}

type Config struct {
}

type writer struct {
	cfg  *Config
	data []proto.Format
}

func New(cfg *Config) Writer {
	return &writer{
		cfg: cfg,
	}
}

func DefaultConfig() *Config {
	return &Config{}
}

func (w *writer) Run(name string, data []proto.Format) error {
	var err error
	w.data = data

	if strings.HasSuffix(name, ".json") {
		err = w.writeJson(name)
	} else if strings.HasSuffix(name, ".txt") {
		err = w.writeTxt(name)
	} else {
		return errors.New("invalid suffix")
	}

	return err
}

func (w *writer) writeJson(name string) error {
	buf := make(map[string][]proto.Format)
	buf[sheetName] = w.data

	b, err := json.Marshal(buf)
	if err != nil {
		return errors.Wrap(err, "failed to marshal")
	}

	if err := os.WriteFile(name, b, perm); err != nil {
		return errors.Wrap(err, "failed to write")
	}

	return nil
}

func (w *writer) writeTxt(name string) error {
	helper := func() []string {
		var buf []string

		var h []string
		r := reflect.TypeOf(proto.Format{})
		for i := 0; i < r.NumField(); i++ {
			h = append(h, r.Field(i).Tag.Get("json"))
		}
		buf = append(buf, strings.Join(h, sep))

		for _, val := range w.data {
			var c []string
			v := reflect.ValueOf(val)
			for i := 0; i < v.NumField(); i++ {
				f := v.Field(i)
				switch f.Kind() {
				case reflect.String:
					c = append(c, f.String())
				case reflect.Int:
					c = append(c, strconv.FormatInt(f.Int(), base))
				}
			}
			buf = append(buf, strings.Join(c, sep))
		}

		return buf
	}

	f, err := os.Create(name)
	if err != nil {
		return errors.Wrap(err, "failed to create")
	}
	defer func() {
		_ = f.Close()
	}()

	b := bufio.NewWriter(f)

	if _, err := b.WriteString(strings.Join(helper(), "\n")); err != nil {
		return errors.Wrap(err, "failed to write")
	}
	defer func() {
		_ = b.Flush()
	}()

	return nil
}
