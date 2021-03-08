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

package proto

// Prototype
// {
//   "lint": [
//     {
//       "file": "name",
//       "line": 1,
//       "type": "Error",
//       "details": "text"
//     }
//   ]
// }

const (
	Base64Content = ".base64"
	Base64Patch   = "diff.base64"
)

const (
	TypeError = "Error"
	TypeInfo  = "Info"
	TypeWarn  = "Warning"
)

type Format struct {
	File    string `json:"file"`
	Line    int    `json:"line"`
	Type    string `json:"type"`
	Details string `json:"details"`
}
