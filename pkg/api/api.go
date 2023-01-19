// Copyright 2023 The envd Authors
//
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

package api

import (
	"io/fs"
	"log"
	"sort"
	"strconv"
	"strings"

	"github.com/cockroachdb/errors"
	"github.com/tensorchord/envd/envd"
	"github.com/urfave/cli/v2"
)

const StableAPI = "v0"

// APIOptions contains "stable" / "latest" and all actual api version: "v0", "v1", ...
var APIOptions map[string]string = map[string]string{"stable": StableAPI}

func init() {
	files, err := fs.ReadDir(envd.ApiStubs(), ".")
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	keys := []string{}
	for _, f := range files {
		APIOptions[f.Name()] = f.Name()
		keys = append(keys, f.Name())
	}
	sort.Slice(keys[:], func(i, j int) bool {
		left, _ := strconv.Atoi(strings.ReplaceAll(keys[i], "v", ""))
		right, _ := strconv.Atoi(strings.ReplaceAll(keys[j], "v", ""))
		return left > right
	})
	APIOptions["latest"] = keys[0]
}

func ArgValidator(clicontext *cli.Context, v string) error {
	_, ok := APIOptions[v]
	if ok {
		return nil
	}
	keys := make([]string, len(APIOptions))
	i := 0
	for k := range APIOptions {
		keys[i] = k
		i++
	}
	return errors.Errorf(`Argument syntax only allows %v, found "%v"`, keys, v)
}
