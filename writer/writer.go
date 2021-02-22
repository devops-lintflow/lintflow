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
	"io/ioutil"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/pkg/errors"

	"github.com/craftslab/lintflow/proto"
)

const (
	SEP = ","
)

var (
	sheetName = time.Now().Local().Format(time.RFC3339)
)

type Writer interface {
	Run([]proto.Format) error
}

type Config struct {
	Name string
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

func (w *writer) Run(data []proto.Format) error {
	var err error
	w.data = data

	if strings.HasSuffix(w.cfg.Name, ".json") {
		err = w.writeJson()
	} else if strings.HasSuffix(w.cfg.Name, ".txt") {
		err = w.writeTxt()
	} else if strings.HasSuffix(w.cfg.Name, ".xlsx") {
		err = w.writeXlsx()
	} else {
		return errors.New("invalid suffix")
	}

	return err
}

func (w *writer) writeJson() error {
	buf := make(map[string][]proto.Format)
	buf[sheetName] = w.data

	b, err := json.Marshal(buf)
	if err != nil {
		return errors.Wrap(err, "failed to marshal")
	}

	if err := ioutil.WriteFile(w.cfg.Name, b, 0600); err != nil {
		return errors.Wrap(err, "failed to write")
	}

	return nil
}

func (w *writer) writeTxt() error {
	helper := func() []string {
		var buf []string

		var h []string
		r := reflect.TypeOf(proto.Format{})
		for i := 0; i < r.NumField(); i++ {
			h = append(h, r.Field(i).Tag.Get("json"))
		}
		buf = append(buf, strings.Join(h, SEP))

		for _, val := range w.data {
			var c []string
			v := reflect.ValueOf(val)
			for i := 0; i < v.NumField(); i++ {
				f := v.Field(i)
				switch f.Kind() {
				case reflect.String:
					c = append(c, f.String())
				case reflect.Int:
					c = append(c, strconv.FormatInt(f.Int(), 10))
				}
			}
			buf = append(buf, strings.Join(c, SEP))
		}

		return buf
	}

	f, err := os.Create(w.cfg.Name)
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

func (w *writer) writeXlsx() error {
	helper := func() ([]interface{}, [][]interface{}) {
		var head []interface{}
		var data [][]interface{}

		r := reflect.TypeOf(proto.Format{})
		for i := 0; i < r.NumField(); i++ {
			head = append(head, strings.ToUpper(r.Field(i).Tag.Get("json")))
		}

		for _, val := range w.data {
			var c []interface{}
			v := reflect.ValueOf(val)
			for i := 0; i < v.NumField(); i++ {
				f := v.Field(i)
				switch f.Kind() {
				case reflect.String:
					c = append(c, f.String())
				case reflect.Int:
					c = append(c, strconv.FormatInt(f.Int(), 10))
				}
			}
			data = append(data, c)
		}

		return head, data
	}

	f := excelize.NewFile()

	write := func(row, col string, style, data interface{}) error {
		s, err := f.NewStyle(style)
		if err != nil {
			return errors.Wrap(err, "failed to new style")
		}

		if err := f.SetCellStyle(sheetName, "A"+row, col+row, s); err != nil {
			return errors.Wrap(err, "failed to set style")
		}

		if err := f.SetSheetRow(sheetName, "A"+row, data); err != nil {
			return errors.Wrap(err, "failed to set row")
		}

		return nil
	}

	index := f.NewSheet(sheetName)

	if err := f.SetPanes(sheetName, `{"freeze":true,"split":true,"y_split":1}`); err != nil {
		return errors.Wrap(err, "failed to set pane")
	}

	head, data := helper()

	style := `{"alignment":{"horizontal":"center","vertical":"center"},"font":{"bold":true}}`
	if err := write("1", "D", style, &head); err != nil {
		return errors.Wrap(err, "failed to write head")
	}

	style = `{"alignment":{"horizontal":"center","vertical":"center"},"font":{"bold":false}}`
	offset := 2
	for index := range data {
		if err := write(strconv.Itoa(index+offset), "D", style, &data[index]); err != nil {
			return errors.Wrap(err, "failed to write data")
		}
	}

	f.SetActiveSheet(index)

	if err := f.SaveAs(w.cfg.Name); err != nil {
		return errors.Wrap(err, "failed to save file")
	}

	return nil
}
