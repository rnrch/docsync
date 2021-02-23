// Copyright 2021 rnrch
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package version

import (
	"bytes"
	"html/template"
	"runtime"
	"strings"
)

var (
	version   string
	revision  string
	buildDate string
	goVersion = runtime.Version()
	os        = runtime.GOOS
	arch      = runtime.GOARCH
)

var versionInfoTmpl = `
{{.program}} version  {{.version}}
  Git commit:     {{.revision}}
  Build date:     {{.buildDate}}
  Go version:     {{.goVersion}}
  OS/Arch:        {{.OS}}/{{.Arch}}
`

func Info(program string) string {
	m := map[string]string{
		"program":   program,
		"version":   version,
		"revision":  revision,
		"buildDate": buildDate,
		"goVersion": goVersion,
		"OS":        os,
		"Arch":      arch,
	}
	t := template.Must(template.New("version").Parse(versionInfoTmpl))

	var buf bytes.Buffer
	if err := t.ExecuteTemplate(&buf, "version", m); err != nil {
		panic(err)
	}
	return strings.TrimSpace(buf.String())
}
