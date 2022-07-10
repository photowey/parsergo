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

package astx

import (
	"go/token"
)

type AstSpec struct {
	ID      string
	Name    string
	PkgPath string
	Pkgs    []*PackageSpec
}

type PackageSpec struct {
	Pkg        string
	Alias      string
	Structs    []*StructSpec
	Interfaces []*InterfaceSpec
	Funcs      []*FuncSpec
}

type StructSpec struct {
	Pkg         string
	Alias       string
	Name        string
	Type        token.Pos
	Comments    []string
	Fields      []*FieldSpec
	Methods     []*MethodSpec
	Annotations []*Annotation
}

type FieldSpec struct {
	Struct string
	Name   string
	Type   string
	Ptr    bool
	Tags   []*TagSpec
}

type InterfaceSpec struct {
	Pkg         string
	Name        string
	Type        token.Pos
	Comments    []string
	Methods     []*MethodSpec
	Annotations []*Annotation
}

type MethodSpec struct {
	Pkg      string
	Struct   string
	Name     string
	Comments []string
	Params   []*ParamSpec
	Returns  []*ReturnSpec
}

type FuncSpec struct {
	Pkg      string
	Name     string
	Comments []string
	Params   []*ParamSpec
	Returns  []*ReturnSpec
}

type ParamSpec struct {
	Pkg      string
	FuncName string
	Name     string
	Ptr      bool
	Type     string
}

type ReturnSpec struct {
	Pkg      string
	FuncName string
	Name     string
	Type     string
}

type Annotation struct {
	Pkg    string
	Name   string
	Values string // maybe json ?
}

type TagSpec struct {
	Field string
	Tags  []*Tag
}
type Tag struct {
	Name  string
	Key   string
	Value string
}
