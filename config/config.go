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

package config

type Config struct {
	ApiVersion string   `yaml:"apiVersion"`
	Kind       string   `yaml:"kind"`
	MetaData   MetaData `yaml:"metadata"`
	Spec       Spec     `yaml:"spec"`
}

type MetaData struct {
	Name string `yaml:"name"`
}

type Spec struct {
	Flow   Flow   `yaml:"flow"`
	Lints  []Lint `yaml:"lint"`
	Review Review `yaml:"review"`
}

type Flow struct {
	Timeout string `yaml:"timeout"`
}

type Lint struct {
	Name   string `yaml:"name"`
	Host   string `yaml:"host"`
	Port   int    `yaml:"port"`
	Filter Filter `yaml:"filter"`
	Vote   string `yaml:"vote"`
}

type Filter struct {
	Include Include `yaml:"include"`
}

type Include struct {
	Extensions []string `yaml:"extension"`
	Files      []string `yaml:"file"`
	Repos      []string `yaml:"repo"`
}

type Review struct {
	Name  string `yaml:"name"`
	Url   string `yaml:"url"`
	User  string `yaml:"user"`
	Pass  string `yaml:"pass"`
	Votes []Vote `yaml:"vote"`
}

type Vote struct {
	Label       string `yaml:"label"`
	Approval    string `yaml:"approval"`
	Disapproval string `yaml:"disapproval"`
	Message     string `yaml:"message"`
}

var (
	Build   string
	Version string
)

func New() *Config {
	return &Config{}
}
