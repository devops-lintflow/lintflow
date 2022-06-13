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
	"context"
	"io"
	"log"
	"os"

	"github.com/pkg/errors"
	"gopkg.in/alecthomas/kingpin.v2"
	"gopkg.in/yaml.v3"

	"github.com/craftslab/lintflow/config"
	"github.com/craftslab/lintflow/flow"
	"github.com/craftslab/lintflow/lint"
	"github.com/craftslab/lintflow/review"
	"github.com/craftslab/lintflow/writer"
)

var (
	app        = kingpin.New("lintflow", "Lint Flow").Version(config.Version + "-build-" + config.Build)
	codeReview = app.Flag("code-review", "Code review (bitbucket|gerrit|gitee|github|gitlab)").Required().String()
	commitHash = app.Flag("commit-hash", "Commit hash (SHA-1)").Required().String()
	configFile = app.Flag("config-file", "Config file (.yml)").Required().String()
	outputFile = app.Flag("output-file", "Output file (.json|.txt|.xlsx)").Default().String()
)

func Run() error {
	kingpin.MustParse(app.Parse(os.Args[1:]))

	c, err := initConfig(*configFile)
	if err != nil {
		return errors.Wrap(err, "failed to init config")
	}

	r, err := initReview(c)
	if err != nil {
		return errors.Wrap(err, "failed to init review")
	}

	l, err := initLint(c)
	if err != nil {
		return errors.Wrap(err, "failed to init lint")
	}

	w, err := initWriter(c)
	if err != nil {
		return errors.Wrap(err, "failed to init writer")
	}

	log.Println("flow running")

	if err := runFlow(c, r, l, w); err != nil {
		return errors.Wrap(err, "failed to run flow")
	}

	log.Println("flow exiting")

	return nil
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

	buf, err := io.ReadAll(fi)
	if err != nil {
		return c, errors.Wrap(err, "failed to readall")
	}

	if err := yaml.Unmarshal(buf, c); err != nil {
		return c, errors.Wrap(err, "failed to unmarshal")
	}

	return c, nil
}

func initReview(cfg *config.Config) (review.Review, error) {
	c := review.DefaultConfig()
	if c == nil {
		return nil, errors.New("failed to config")
	}

	c.Name = *codeReview
	c.Reviews = cfg.Spec.Review

	return review.New(c), nil
}

func initLint(cfg *config.Config) (lint.Lint, error) {
	c := lint.DefaultConfig()
	if c == nil {
		return nil, errors.New("failed to config")
	}

	c.Lints = cfg.Spec.Lint

	return lint.New(c), nil
}

func initWriter(_ *config.Config) (writer.Writer, error) {
	c := writer.DefaultConfig()
	if c == nil {
		return nil, errors.New("failed to config")
	}

	return writer.New(c), nil
}

func runFlow(c *config.Config, r review.Review, l lint.Lint, w writer.Writer) error {
	cfg := flow.DefaultConfig()
	if cfg == nil {
		return errors.New("failed to config flow")
	}

	cfg.Config = *c
	cfg.Lint = l
	cfg.Review = r

	f := flow.New(context.Background(), cfg)
	if f == nil {
		return errors.New("failed to new flow")
	}

	buf, err := f.Run(*commitHash)
	if err != nil {
		return errors.Wrap(err, "failed to run flow")
	}

	if *outputFile != "" {
		if _, err = os.Stat(*outputFile); err == nil {
			return errors.New("file already exists")
		}
		if len(buf) != 0 {
			if err = w.Run(*outputFile, buf); err != nil {
				return errors.Wrap(err, "failed to run writer")
			}
		} else {
			return nil
		}
	}

	return nil
}
