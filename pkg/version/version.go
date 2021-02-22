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
