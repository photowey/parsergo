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

package parser

import (
	"github.com/photowey/parsergo/astx"
	"github.com/photowey/parsergo/loader"
)

type Parser interface {
	Parse(pkg *loader.Package) *astx.AstSpec
	StructParser
	InterfaceParser
	MethodParser
	FuncParser
}

type StructParser interface {
	ParseStructs(aw *astx.Astx) *astx.PackageSpec
}

type InterfaceParser interface {
	ParseInterfaces(aw *astx.Astx, ps *astx.PackageSpec)
}

type MethodParser interface {
	ParseMethods(aw *astx.Astx, ps *astx.PackageSpec)
}

type FuncParser interface {
	ParseFuncs(aw *astx.Astx, ps *astx.PackageSpec)
}
