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
	"context"
	"encoding/json"
	"io"
	"math"
	"os"
	"path/filepath"
	"strconv"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/devops-lintflow/lintflow/config"
	"github.com/devops-lintflow/lintflow/format"
)

type Lint interface {
	Run(context.Context, string, string, []string, func(*config.Filter, string, string) bool) (map[string][]format.Report, error)
}

type Config struct {
	Lints []config.Lint
}

type lint struct {
	cfg *Config
}

func New(cfg *Config) Lint {
	return &lint{
		cfg: cfg,
	}
}

func DefaultConfig() *Config {
	return &Config{}
}

// nolint:gosec
func (l *lint) Run(ctx context.Context, root, repo string, files []string,
	match func(*config.Filter, string, string) bool) (map[string][]format.Report, error) {
	helper := func(filter *config.Filter, files []string) []string {
		var buf []string
		for _, item := range files {
			if match(filter, repo, item) {
				buf = append(buf, item)
			}
		}
		return buf
	}

	type result struct {
		data map[string][]format.Report
		err  error
	}

	bypass := true
	ch := make(chan result, len(l.cfg.Lints))

	for i := range l.cfg.Lints {
		buf := helper(&l.cfg.Lints[i].Filter, files)
		if len(buf) != 0 {
			bypass = false
		}
		go func(ctx context.Context, f []string, v config.Lint) {
			if len(f) != 0 {
				m, e := l.marshal(root, f)
				if e != nil {
					ch <- result{nil, errors.Wrap(e, "failed to marshal")}
					return
				}
				r, e := l.routine(ctx, v.Host, v.Port, m)
				if e != nil {
					ch <- result{nil, errors.Wrap(e, "failed to routine")}
					return
				}
				ch <- result{r, nil}
			} else {
				ch <- result{map[string][]format.Report{}, nil}
			}
		}(ctx, buf, l.cfg.Lints[i])
	}

	if bypass {
		return nil, nil
	}

	ret := map[string][]format.Report{}

	for range l.cfg.Lints {
		r := <-ch
		if r.err != nil {
			return nil, r.err
		}
		if len(r.data) != 0 {
			for k, v := range r.data {
				if len(v) == 0 {
					continue
				}
				if _, ok := ret[k]; !ok {
					ret[k] = v
				} else {
					ret[k] = append(ret[k], v...)
				}
			}
		}
	}

	return ret, nil
}

func (l *lint) marshal(root string, data []string) ([]byte, error) {
	helper := func(name string) (string, error) {
		fi, err := os.Open(name)
		if err != nil {
			return "", errors.Wrap(err, "failed to open")
		}
		defer func() { _ = fi.Close() }()
		buf, err := io.ReadAll(fi)
		if err != nil {
			return "", errors.Wrap(err, "failed to readall")
		}
		return string(buf), nil
	}

	var err error
	buf := map[string]string{}

	for _, val := range data {
		if val == "" {
			err = errors.New("invalid data")
			break
		}
		buf[val], err = helper(filepath.Join(root, val))
		if err != nil {
			break
		}
	}

	if len(buf) == 0 {
		return nil, errors.New("invalid data")
	}

	if err != nil {
		return nil, errors.Wrap(err, "failed to read")
	}

	ret, err := json.Marshal(buf)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal")
	}

	return ret, nil
}

func (l *lint) routine(ctx context.Context, host string, port int, data []byte) (map[string][]format.Report, error) {
	helper := func(data string) (map[string][]format.Report, error) {
		var buf map[string][]format.Report
		if err := json.Unmarshal([]byte(data), &buf); err != nil {
			return nil, errors.Wrap(err, "failed to unmarshal")
		}
		return buf, nil
	}

	// nolint: staticcheck
	conn, err := grpc.DialContext(ctx, host+":"+strconv.Itoa(port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(math.MaxInt32), grpc.MaxCallSendMsgSize(math.MaxInt32)))
	if err != nil {
		return nil, errors.Wrap(err, "failed to dial")
	}
	defer func() { _ = conn.Close() }()

	client := NewLintProtoClient(conn)

	reply, err := client.SendLint(ctx, &LintRequest{Message: string(data)})
	if err != nil {
		return nil, errors.Wrap(err, "failed to send")
	}

	buf, err := helper(reply.GetMessage())
	if err != nil {
		return nil, errors.Wrap(err, "failed to get")
	}

	return buf, nil
}
