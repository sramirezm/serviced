// Copyright 2014 The Serviced Authors.
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

// Package agent implements a service that runs on a serviced node. It is
// responsible for ensuring that a particular node is running the correct services
// and reporting the state and health of those services back to the master
// serviced.

package utils

import (
	"os"
	"path"
	"runtime"
	"strings"
)

//ServiceDHome gets the home location of serviced by looking at the enviornment
func ServiceDHome() string {
	return os.Getenv("SERVICED_HOME")
}

//LocalDir gets the absolute path to a particular directory under ServiceDHome
// if SERVICED_HOME is not defined then we use the location of the caller
func LocalDir(p string) string {
	homeDir := ServiceDHome()
	if len(homeDir) == 0 {
		_, filename, _, _ := runtime.Caller(1)
		homeDir = strings.Replace(path.Dir(filename), "utils", "", 1)
	}
	return path.Join(homeDir, p)
}

// ResourcesDir points to internal services resources directory
func ResourcesDir() string {
	return LocalDir("isvcs/resources")
}
