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
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"strings"
	"text/template"

	"github.com/jessevdk/go-flags"
	"github.com/rnrch/rlog"
)

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
	Templates []string `long:"template" short:"t" description:"template files" required:"true"`
	Ignore    []string `long:"ignore" short:"i" description:"ignore files"`
	Output    string   `long:"output" short:"o" description:"output file name" default:"output.md"`
	Directory string   `long:"directory" short:"d" description:"directory to process"`
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
	return options
}

func main() {
	rlog.SwtichMode(rlog.Development)
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

	f, err := processFolder(options.Directory, options.Ignore, 1)
	if err != nil {
		rlog.Error(err, "process folder", "folder", options.Directory, "ignore", options.Ignore)
		os.Exit(1)
	}

	err = write(f, options.Templates, options.Output)
	if err != nil {
		rlog.Error(err, "write to file", "folder", f, "templates", options.Templates, "output", options.Output)
		os.Exit(1)
	}
}

func processFolder(folder string, ignore []string, depth int) (Folder, error) {
	f := Folder{Path: folder, Name: path.Base(folder), Depth: depth}
	contents, err := ioutil.ReadDir(folder)
	if err != nil {
		return f, err
	}
	for _, content := range contents {
		if contains(ignore, content.Name()) {
			continue
		}
		if !content.IsDir() {
			f.Files = append(f.Files, File{
				Name: content.Name(),
				Path: path.Join(folder, content.Name()),
			})
			continue
		}
		sub, err := processFolder(path.Join(folder, content.Name()), ignore, depth+1)
		if err != nil {
			return f, err
		}
		f.SubFolders = append(f.SubFolders, sub)
	}
	return f, nil
}

func contains(set []string, value string) bool {
	for _, s := range set {
		match, err := regexp.MatchString(s, value)
		if err != nil {
			rlog.Error(err, "regexp match")
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
	}
	t, err := template.New(path.Base(templates[0])).Funcs(fm).ParseFiles(templates...)
	if err != nil {
		return err
	}
	out, err := os.Create(output)
	defer out.Close()
	if err != nil {
		return err
	}
	return t.Execute(out, folder)
}

func depthToHeader(depth int) string {
	return strings.Repeat("#", depth)
}
