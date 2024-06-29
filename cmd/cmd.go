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
	"time"

	"github.com/alecthomas/kingpin/v2"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"

	"github.com/devops-lintflow/lintflow/config"
	"github.com/devops-lintflow/lintflow/flow"
	"github.com/devops-lintflow/lintflow/lint"
	"github.com/devops-lintflow/lintflow/review"
)

const (
	Timeout = 120 * time.Second
)

var (
	app        = kingpin.New("lintflow", "Lint Flow").Version(config.Version + "-build-" + config.Build)
	codeReview = app.Flag("code-review", "Code review (bitbucket|gerrit|gitee|github|gitlab)").Required().String()
	commitHash = app.Flag("commit-hash", "Commit hash (SHA-1)").Required().String()
	configFile = app.Flag("config-file", "Config file (.yml)").Required().String()
)

func Run(ctx context.Context) error {
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

	log.Println("flow running")

	if err := runFlow(ctx, c, r, l); err != nil {
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

	c.Review = cfg.Spec.Review

	return review.New(c), nil
}

func initLint(cfg *config.Config) (lint.Lint, error) {
	c := lint.DefaultConfig()
	if c == nil {
		return nil, errors.New("failed to config")
	}

	c.Lints = cfg.Spec.Lints

	return lint.New(c), nil
}

func runFlow(ctx context.Context, c *config.Config, r review.Review, l lint.Lint) error {
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

	timeout, err := setTimeout(c.Spec.Flow.Timeout)
	if err != nil {
		return errors.Wrap(err, "failed to set timeout")
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	if err = f.Run(ctx, *commitHash); err != nil {
		return errors.Wrap(err, "failed to run flow")
	}

	return nil
}

func setTimeout(timeout string) (time.Duration, error) {
	var t time.Duration
	var err error

	if timeout != "" {
		t, err = time.ParseDuration(timeout)
		if err != nil {
			return 0, errors.Wrap(err, "failed to parse duration")
		}
	} else {
		t = Timeout
	}

	return t, nil
}
