/*
Copyright © 2020 rnrch <rnrch@outlook.com>
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"text/template"
)

var (
	tmpls       stringSlice
	output      string
	folder      string
	excludes    stringSlice
	excludesAll stringSlice
)

type stringSlice []string

func (s *stringSlice) String() string {
	return strings.Join(*s, ",")
}

func (s *stringSlice) Set(value string) error {
	*s = append(*s, value)
	return nil
}

// Folder represents a directory
type Folder struct {
	Name    string
	Folders []Folder
	Files   []string
	Layer   int
}

// Options is the data we render on the template.
type Options struct {
	Data        *Folder
	ExcludeRoot bool
}

func init() {
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Printf("failed to get current directory, err: %s \n", err)
		return
	}
	flag.Var(&tmpls, "f", "name of the template files")
	flag.Var(&excludes, "e", "name of the subfolders to be excluded, use absolute path, otherwise means the folder under root folder")
	flag.Var(&excludesAll, "ea", "name of subfolders to be excluded, use relative path")
	flag.StringVar(&output, "o", "output.md", "name of the output file")
	flag.StringVar(&folder, "d", pwd, "absolute path of the directory to process")
	flag.Parse()
	if tmpls == nil {
		fmt.Print("no template specified \n")
		return
	}
	if !path.IsAbs(folder) {
		folder = path.Join(pwd, folder)
	}
	for i, e := range excludes {
		if !path.IsAbs(e) {
			excludes[i] = path.Join(folder, e)
		}
	}
}

func main() {
	t, err := processDir(folder, 1)
	if err != nil {
		fmt.Print("failed to process dir")
		return
	}
	d := Options{
		Data:        t,
		ExcludeRoot: inExcludes(folder, excludes),
	}
	writeToFile(d)
}

func inExcludes(dir string, excludes []string) bool {
	for _, s := range excludes {
		if dir == s {
			return true
		}
	}
	return false
}

func inExcludesAll(dir string, excludesAll []string) bool {
	for _, s := range excludesAll {
		if dir == s {
			return true
		}
	}
	return false
}

func processDir(dir string, layer int) (*Folder, error) {
	contents, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Printf("can not read dir, err: %s \n", err)
		return nil, err
	}
	res := &Folder{
		Folders: []Folder{},
		Files:   []string{},
		Name:    path.Base(dir),
		Layer:   layer,
	}
	for _, content := range contents {
		if content.IsDir() {
			dirname := path.Join(dir, content.Name())
			if inExcludes(dirname, excludes) {
				continue
			}
			if inExcludesAll(content.Name(), excludesAll) {
				continue
			}
			subfolder, _ := processDir(dirname, layer+1)
			res.Folders = append(res.Folders, *subfolder)
			continue
		}
		if !inExcludes(path.Join(dir, content.Name()), excludes) && !inExcludesAll(content.Name(), excludesAll) {
			res.Files = append(res.Files, path.Join(dir, content.Name()))
		}
	}
	return res, nil
}

func writeToFile(o Options) {
	t := template.New(path.Base(tmpls[0]))
	processStyles(t)
	t, err := t.ParseFiles(tmpls...)
	if err != nil {
		fmt.Printf("can not parse templates, err: %s \n", err)
		return
	}

	f, err := os.Create(output)
	defer f.Close()
	if err != nil {
		fmt.Printf("can not create file %s, err: %s \n", output, err)
		return
	}

	err = t.Execute(f, o)
	if err != nil {
		fmt.Printf("can not exec template, err: %s \n", err)
		return
	}
}

func processStyles(t *template.Template) {
	fm := template.FuncMap{
		"folderTOC":     folderTOC,
		"layerToHeader": layerToHeader,
		"listFileName":  listFileName,
		"fileLinks":     fileLinks,
	}
	t.Funcs(fm)
}
