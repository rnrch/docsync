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
	"path"
	"strings"
)

func layerToHeader(layer int) string {
	return strings.Repeat("#", layer)
}

func listFileName(s string) string {
	name := path.Base(s)
	return strings.Join([]string{"-", name}, " ")
}

func folderTOC(folder Folder) string {
	whitespace := strings.Repeat(" ", (folder.Layer-1)*2)
	return whitespace + "- [" + folder.Name + "](#" + strings.ReplaceAll(folder.Name, " ", "-") + ")"
}

func fileLinks(s string) string {
	name := path.Base(s)
	name = strings.Split(name, ".")[0]
	filename := strings.Split(s, folder)[1]
	filename = strings.TrimLeft(filename, "/")
	return "- [" + name + "](" + strings.ReplaceAll(filename, " ", "%20") + ")"
}
