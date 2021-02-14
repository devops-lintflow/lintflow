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
	"io/ioutil"
	"log"
	"os"

	"github.com/pkg/errors"
	"gopkg.in/alecthomas/kingpin.v2"
	"gopkg.in/yaml.v3"

	"github.com/craftslab/lintflow/config"
)

var (
	app        = kingpin.New("lintflow", "Lint Flow").Version(config.Version + "-build-" + config.Build)
	configFile = app.Flag("config-file", "Config file (.yml)").Required().String()
)

func Run() {
	kingpin.MustParse(app.Parse(os.Args[1:]))

	c, err := initConfig(*configFile)
	if err != nil {
		log.Fatalf("failed to init config: %v", err)
	}

	log.Println("flow running")

	if err := runFlow(c); err != nil {
		log.Fatalf("failed to run flow: %v", err)
	}

	log.Println("flow exiting")
}

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

func runFlow(_ *config.Config) error {
	return nil
}
