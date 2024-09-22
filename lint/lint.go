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
	Run(context.Context, string, string, []string, string, string,
		func(*config.Filter, string, string) bool) (map[string][]format.Report, error)
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

// nolint:gocyclo
func (l *lint) Run(ctx context.Context, root, repo string, files []string, meta, patch string,
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
		go func(ctx context.Context, lint config.Lint, files []string, meta, patch string) {
			if len(files) != 0 {
				req, err := l.encode(lint.Name, root, files, meta, patch)
				if err != nil {
					ch <- result{nil, errors.Wrap(err, "failed to encode")}
					return
				}
				ret, err := l.routine(ctx, lint.Host, lint.Port, req)
				if err != nil {
					ch <- result{nil, errors.Wrap(err, "failed to routine")}
					return
				}
				rep, err := l.decode(ret)
				if err != nil {
					ch <- result{nil, errors.Wrap(err, "failed to decode")}
					return
				}
				ch <- result{rep, nil}
			} else {
				ch <- result{map[string][]format.Report{}, nil}
			}
		}(ctx, l.cfg.Lints[i], buf, meta, patch)
	}

	if bypass {
		return nil, nil
	}

	ret := map[string][]format.Report{}

	for range l.cfg.Lints {
		rep := <-ch
		if rep.err != nil {
			return nil, rep.err
		}
		if len(rep.data) != 0 {
			for name, reports := range rep.data {
				if _, ok := ret[name]; !ok {
					ret[name] = reports
				} else {
					ret[name] = append(ret[name], reports...)
				}
			}
		}
	}

	return ret, nil
}

func (l *lint) routine(ctx context.Context, host string, port int, request *LintRequest) (*LintReply, error) {
	conn, err := grpc.NewClient(host+":"+strconv.Itoa(port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(math.MaxInt32), grpc.MaxCallSendMsgSize(math.MaxInt32)))
	if err != nil {
		return nil, errors.Wrap(err, "failed to dial")
	}

	defer func() {
		_ = conn.Close()
	}()

	client := NewLintProtoClient(conn)

	reply, err := client.SendLint(ctx, request)
	if err != nil {
		return nil, errors.Wrap(err, "failed to send lint")
	}

	return reply, nil
}

func (l *lint) encode(lint, root string, files []string, meta, patch string) (*LintRequest, error) {
	var err error

	helper := func(path string) ([]byte, error) {
		fi, err := os.Open(path)
		if err != nil {
			return nil, errors.Wrap(err, "failed to open")
		}
		defer func() {
			_ = fi.Close()
		}()
		buf, err := io.ReadAll(fi)
		if err != nil {
			return nil, errors.Wrap(err, "failed to read")
		}
		return buf, nil
	}

	request := &LintRequest{
		Name: lint,
	}

	request.LintFiles = []*LintFile{}

	for _, item := range files {
		buf := LintFile{Path: item}
		buf.Content, err = helper(filepath.Join(root, item))
		if err != nil {
			break
		}
		request.LintFiles = append(request.LintFiles, &buf)
	}

	if err != nil {
		return nil, errors.New("invalid file")
	}

	request.LintMeta = &LintMeta{}
	request.LintMeta.Path = meta

	request.LintMeta.Content, err = helper(filepath.Join(root, meta))
	if err != nil {
		return nil, errors.New("invalid meta")
	}

	request.LintPatch = &LintPatch{}
	request.LintPatch.Path = patch

	request.LintPatch.Content, err = helper(filepath.Join(root, patch))
	if err != nil {
		return nil, errors.New("invalid patch")
	}

	return request, nil
}

func (l *lint) decode(reply *LintReply) (map[string][]format.Report, error) {
	name := reply.GetName()
	if name == "" {
		return map[string][]format.Report{}, nil
	}

	buf := map[string][]format.Report{
		name: {},
	}

	reports := reply.GetLintReports()

	for i := range reports {
		b := format.Report{
			File:    reports[i].GetFile(),
			Line:    int(reports[i].GetLine()),
			Type:    reports[i].GetType(),
			Details: reports[i].GetDetails(),
		}
		buf[name] = append(buf[name], b)
	}

	return buf, nil
}
