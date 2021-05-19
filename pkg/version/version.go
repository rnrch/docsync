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
	"encoding/json"
	"runtime"

	"github.com/rnrch/rlog"
)

// Build information. Populated at build-time.
var (
	Version   string
	Revision  string
	Branch    string
	BuildDate string
	GoVersion = runtime.Version()
	OS        = runtime.GOOS
	Arch      = runtime.GOARCH
)

// Info provides the iterable version information.
var Info = map[string]string{
	"version":   Version,
	"revision":  Revision,
	"branch":    Branch,
	"buildDate": BuildDate,
	"goVersion": GoVersion,
	"os":        OS,
	"arch":      Arch,
}

func String() string {
	v, err := json.MarshalIndent(Info, "", "  ")
	if err != nil {
		rlog.Error(err, "marshal version info")
		return ""
	}
	return string(v)
}
