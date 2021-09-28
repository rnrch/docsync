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

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/jessevdk/go-flags"
	"github.com/rnrch/docsync/pkg/version"
	"github.com/rnrch/rlog"
)

var defaultTemplate = `{{ define "folder" }}
{{ .Depth | depthToHeader }} {{ .Name }}
{{ range .Files }}
[{{ .Name }}]({{ .Path }})
{{ end }}
{{- range .SubFolders }}
{{- template "folder" . }}
{{- end }}
{{- end }}

{{- template "folder" . -}}
`

type File struct {
	Name string
	Path string
}

type Folder struct {
	Depth      int
	Path       string
	Name       string
	SubFolders []Folder
	Files      []File
}

type Options struct {
	Templates []string `long:"template" short:"t" description:"template files, the main template should be the first one"`
	Include   []string `long:"include" short:"i" description:"files to be included. If a file name does not match any of the patterns specified, it is ignored."`
	Exclude   []string `long:"exclude" short:"e" description:"files and folders to be excluded. Priority higher than include flag."`
	Output    string   `long:"output" short:"o" description:"output file name" default:"output.md"`
	Directory string   `long:"directory" short:"d" description:"directory to process"`
	Version   bool     `long:"version" short:"v" description:"show version info"`
}

func ParseFlags() Options {
	options := Options{}
	parser := flags.NewParser(&options, flags.Default)
	if _, err := parser.Parse(); err != nil {
		code := 1
		if fe, ok := err.(*flags.Error); ok {
			if fe.Type == flags.ErrHelp {
				code = 0
			}
		}
		os.Exit(code)
	}
	if options.Version {
		fmt.Println(version.String())
		os.Exit(0)
	}
	return options
}

func main() {
	options := ParseFlags()
	if options.Directory == "" {
		pwd, err := os.Getwd()
		if err != nil {
			rlog.Error(err, "get working directory")
			os.Exit(1)
		}
		options.Directory = pwd
	}
	rlog.Info("Parsed options", "options", options)

	f, err := processFolder(options.Directory, options.Include, options.Exclude, 1)
	if err != nil {
		rlog.Error(err, "process folder", "folder", options.Directory)
		os.Exit(1)
	}

	err = write(f, options.Templates, options.Output)
	if err != nil {
		rlog.Error(err, "write to file", "folder", f, "templates", options.Templates, "output", options.Output)
		os.Exit(1)
	}
}

func processFolder(folder string, include []string, exclude []string, depth int) (Folder, error) {
	f := Folder{Path: folder, Name: filepath.Base(folder), Depth: depth}
	contents, err := ioutil.ReadDir(folder)
	if err != nil {
		return f, err
	}
	for _, content := range contents {
		if !content.IsDir() {
			if contains(exclude, content.Name()) {
				continue
			}
			if !contains(include, content.Name()) {
				continue
			}
			n := strings.TrimSuffix(content.Name(), filepath.Ext(content.Name()))
			f.Files = append(f.Files, File{
				Name: n,
				Path: filepath.Join(folder, content.Name()),
			})
			continue
		}
		if contains(exclude, content.Name()) {
			continue
		}
		sub, err := processFolder(filepath.Join(folder, content.Name()), include, exclude, depth+1)
		if err != nil {
			return f, err
		}
		f.SubFolders = append(f.SubFolders, sub)
	}
	return f, nil
}

func contains(set []string, value string) bool {
	for _, s := range set {
		match, err := filepath.Match(s, value)
		if err != nil {
			rlog.Error(err, "match path", "pattern", s, "value", value)
			continue
		}
		if match {
			return true
		}
	}
	return false
}

func write(folder Folder, templates []string, output string) error {
	fm := template.FuncMap{
		"depthToHeader": depthToHeader,
		"processPwd":    processPwd,
	}
	var tmpl *template.Template
	var err error
	if len(templates) == 0 {
		tmpl, err = template.New("tmpl").Funcs(fm).Parse(defaultTemplate)
	} else {
		tmpl, err = template.New(filepath.Base(templates[0])).Funcs(fm).ParseFiles(templates...)
	}
	if err != nil {
		return err
	}
	out, err := os.Create(output)
	if err != nil {
		return err
	}
	defer out.Close()
	return tmpl.Execute(out, folder)
}

func depthToHeader(depth int) string {
	return strings.Repeat("#", depth)
}

func processPwd(name string) string {
	if name == "." {
		pwd, err := os.Getwd()
		if err != nil {
			log.Printf("get wd err %v", err)
			return name
		}
		return filepath.Base(pwd)
	}
	return name
}
