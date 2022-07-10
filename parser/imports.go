/*
Copyright The Kubernetes Authors.

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

package parser

import (
	"fmt"
	"path"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/photowey/parsergo/loader"
	"github.com/photowey/parsergo/sets"
)

type Imports struct {
	byPath  sets.StringMap
	byAlias sets.StringMap

	pkg *loader.Package
}

func (its *Imports) ImportSpecs() []string {
	res := make([]string, 0, len(its.byPath))
	for importPath, alias := range its.byPath {
		pkg := its.pkg.Imports()[importPath]
		if pkg != nil && pkg.Name == alias {
			res = append(res, fmt.Sprintf("%q", importPath))
		} else {
			res = append(res, fmt.Sprintf("%s %q", alias, importPath))
		}
	}
	return res
}

func (its *Imports) NeedImport(importPath string) string {
	if ind := strings.LastIndex(importPath, "/vendor/"); ind != -1 {
		importPath = importPath[ind+8:/* len("/vendor/") */ ]
	}
	alias, exists := its.byPath[importPath]
	if exists {
		return alias
	}
	restPath, nextWord := path.Split(importPath)
	for otherPath, exists := "", true; exists && otherPath != importPath; otherPath, exists = its.byAlias[alias] {
		if restPath == "" {
			alias += "x"
		}
		for firstRune, runeLen := utf8.DecodeRuneInString(nextWord); unicode.IsDigit(firstRune); firstRune, runeLen = utf8.DecodeRuneInString(nextWord) {
			nextWord = nextWord[runeLen:]
		}
		nextWord = strings.Map(func(r rune) rune {
			if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' {
				return r
			}
			return '_'
		}, nextWord)

		alias = nextWord + alias
		if len(restPath) > 0 {
			restPath, nextWord = path.Split(restPath[:len(restPath)-1] /* chop off final slash */)
		}
	}

	its.byPath[importPath] = alias
	its.byAlias[alias] = importPath

	return alias
}
