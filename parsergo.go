/*
 * Copyright © 2022 photowey (photowey@gmail.com)
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package parsergo

import (
	"github.com/photowey/parsergo/astx"
	"github.com/photowey/parsergo/loader"
	"github.com/photowey/parsergo/parser"
)

var _ PackageScanner = (*scanner)(nil)

type Scanner interface {
	Scan() []*astx.AstSpec
}

type PackageScanner interface {
	Scanner
	ScanPackages(rootPaths ...string) []*astx.AstSpec
}

type scanner struct {
	Paths []string
}

func (scr *scanner) Scan() []*astx.AstSpec {
	return scr.ScanPackages(scr.Paths...)
}

func (scr *scanner) ScanPackages(rootPaths ...string) []*astx.AstSpec {
	paths := toSlice(rootPaths...)
	if len(paths) == 0 {
		paths = append(paths, "./...")
	}

	roots, err := loader.LoadRoots(paths...)
	if err != nil {
		panic(err)
	}

	ass := make([]*astx.AstSpec, 0, len(roots))

	for _, root := range roots {
		as := parser.Parse(root)
		ass = append(ass, as)
	}

	return ass
}

func NewScanner(rootPaths ...string) PackageScanner {
	return &scanner{
		Paths: toSlice(rootPaths...),
	}
}

func toSlice(rootPaths ...string) []string {
	return rootPaths
}
