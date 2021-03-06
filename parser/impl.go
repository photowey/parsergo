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
	"fmt"
	"go/ast"
	"strings"

	"github.com/photowey/parsergo/astx"
	"github.com/photowey/parsergo/loader"
	"github.com/photowey/parsergo/pkg/regexpx"
)

type parser struct{}

func (psr parser) Parse(pkg *loader.Package) *astx.AstSpec {
	pkgs := make([]*astx.PackageSpec, 0, len(pkg.CompiledGoFiles))
	for _, cf := range pkg.CompiledGoFiles {
		aw := astx.NewAstx(cf, pkg)
		if aw.Ast.Comments == nil {
			continue
		}

		ps := psr.ParseStructs(aw)
		if len(ps.Structs) == 0 {
			continue
		}

		pkgs = append(pkgs, ps)
	}

	return &astx.AstSpec{
		ID:      pkg.ID,
		Name:    pkg.Name,
		PkgPath: pkg.PkgPath,
		Pkgs:    pkgs,
	}
}

func (psr parser) ParseStructs(aw *astx.Astx) *astx.PackageSpec {
	ps := populatePackageSpec(aw)
	for _, d := range aw.Ast.Decls {
		switch decl := d.(type) {
		case *ast.GenDecl:
		SPEC:
			for _, spec := range decl.Specs {
				switch specVal := spec.(type) {
				case *ast.TypeSpec:
					if st, ok := specVal.Type.(*ast.StructType); ok {
						if decl.Doc == nil {
							continue SPEC
						}
						ss := psr.parserStruct(aw, decl, specVal, st)

						ss.Type = st.Struct
						ps.Structs = append(ps.Structs, ss)
					}
				}
			}
		}
	}

	psr.ParseMethods(aw, ps)
	psr.ParseAnnotations(aw, ps)

	return ps
}

func (psr parser) ParseInterfaces(aw *astx.Astx, ps *astx.PackageSpec) {
	for _, d := range aw.Ast.Decls {
		switch decl := d.(type) {
		case *ast.GenDecl:
		SPEC:
			for _, spec := range decl.Specs {
				switch specVal := spec.(type) {
				case *ast.TypeSpec:
					if it, ok := specVal.Type.(*ast.InterfaceType); ok {
						if decl.Doc == nil {
							continue SPEC
						}
						comments := make([]string, 0, len(decl.Doc.List))
						if decl.Doc != nil {
							for _, comment := range decl.Doc.List {
								comments = append(comments, comment.Text)
							}
						}

						is := &astx.InterfaceSpec{
							Pkg:         ps.Pkg,
							Name:        specVal.Name.String(),
							Comments:    comments,
							Methods:     make([]*astx.MethodSpec, 0),
							Annotations: make([]*astx.Annotation, 0),
						}
						is.Type = it.Interface
						ps.Interfaces = append(ps.Interfaces, is)
					}
				}
			}
		}
	}
}

func (psr parser) ParseMethods(aw *astx.Astx, ps *astx.PackageSpec) {
	for _, d := range aw.Ast.Decls {
		switch funcDecl := d.(type) {
		case *ast.FuncDecl:
			comments := make([]string, 0)
			if funcDecl.Doc != nil {
				for _, comment := range funcDecl.Doc.List {
					comments = append(comments, comment.Text)
				}
			}
			if funcDecl.Recv != nil {
				for _, field := range funcDecl.Recv.List {

					for _, spec := range ps.Structs {
						structName := spec.Name

						stn := ""
						switch ft := field.Type.(type) {
						case *ast.Ident:
							// func (x Xxx) MethodName(...) {}
							stn = ft.Name
						case *ast.StarExpr:
							// func (x *Xxx) MethodName(...) {}
							switch xt := ft.X.(type) {
							case *ast.Ident:
								stn = xt.Name
							case *ast.SelectorExpr:
								// do nothing
							}
						}

						if structName == stn {
							ms := &astx.MethodSpec{
								Pkg:      ps.Pkg,
								Struct:   structName,
								Name:     funcDecl.Name.String(),
								Comments: comments,
								Params:   make([]*astx.ParamSpec, 0),
								Returns:  make([]*astx.ReturnSpec, 0),
							}

							psr.handleParams(funcDecl, ps, ms)
							psr.handleResults(funcDecl, ps, ms)

							spec.Methods = append(spec.Methods, ms)
						}
					}
				}
			}
		}
	}
}

func (psr parser) ParseFuncs(aw *astx.Astx, ps *astx.PackageSpec) {
	for _, d := range aw.Ast.Decls {
		switch decl := d.(type) {
		case *ast.FuncDecl:
			comments := make([]string, 0)
			if decl.Doc != nil {
				for _, comment := range decl.Doc.List {
					comments = append(comments, comment.Text)
				}
			}
			if decl.Recv == nil {
				fs := &astx.FuncSpec{
					Pkg:      aw.Pkg,
					Name:     decl.Name.String(),
					Comments: comments,
					Params:   make([]*astx.ParamSpec, 0),
					Returns:  make([]*astx.ReturnSpec, 0),
				}
				ps.Funcs = append(ps.Funcs, fs)
			}
		}
	}
}

func (psr parser) ParseAnnotations(aw *astx.Astx, ps *astx.PackageSpec) {
	for _, spec := range ps.Structs {
		comments := spec.Comments
		for _, comment := range comments {
			// `// @Service`
			// `// @Service("helloService")`
			// `// @ComponentScan({"path":"github.com/photowey/parsergo/tests","excludes":["github.com/photowey/parsergo/tests/structx"]})`
			cmt := comment[3:]
			if ok := strings.HasPrefix(cmt, "@"); ok {
				anno := &astx.Annotation{
					Pkg:   spec.Pkg,
					Alias: spec.Alias,
					Anno:  comment,
				}
				anno.Name = regexpx.RegexpExtract(`^@(?P<annotation>[\S]+)[\(.*\)]?`, cmt, "$annotation")
				if strings.Contains(cmt, "(") && strings.Contains(cmt, ")") {
					// `// @Service("helloService")`
					anno.Values = regexpx.RegexpExtract(`^@.*\((?P<value>[\S]+)\)`, cmt, "$value")
				}

				spec.Annotations = append(spec.Annotations, anno)
			}
		}
	}
}

func (psr parser) parserStruct(aw *astx.Astx, decl *ast.GenDecl, specVal *ast.TypeSpec, st *ast.StructType) *astx.StructSpec {
	comments := make([]string, 0, len(decl.Doc.List))
	if decl.Doc != nil {
		for _, comment := range decl.Doc.List {
			comments = append(comments, comment.Text)
		}
	}

	ss := &astx.StructSpec{
		Pkg:         aw.Pkg,
		Alias:       aw.Package.Name,
		Name:        specVal.Name.String(),
		Comments:    comments,
		Fields:      make([]*astx.FieldSpec, 0),
		Methods:     make([]*astx.MethodSpec, 0),
		Annotations: make([]*astx.Annotation, 0),
	}

	if fields := st.Fields; fields != nil && fields.List != nil {
		for _, field := range fields.List {
			fs := &astx.FieldSpec{
				Struct: specVal.Name.String(),
				Name:   field.Names[0].Name,
				Tags:   make([]*astx.TagSpec, 0),
			}

			// handle field's type
			psr.handleFieldType(field, fs)
			// handle field's tag
			psr.handleFieldTag(field, fs)

			ss.Fields = append(ss.Fields, fs)
		}
	}

	return ss
}

func (psr parser) handleFieldType(field *ast.Field, fs *astx.FieldSpec) {
	switch expr := field.Type.(type) {
	case *ast.Ident: // Xxx Yyy `k:"v"`
		fs.Type = expr.Name
	case *ast.StarExpr: // Xxx *Yyy `k:"v"`
		switch starExpr := expr.X.(type) {
		case *ast.SelectorExpr:
			if sx, oke := starExpr.X.(*ast.Ident); oke {
				fs.Type = fmt.Sprintf("*%s.%s", sx.Name, starExpr.Sel.Name)
			} else {
				fs.Type = fmt.Sprintf("*%s", starExpr.Sel.Name)
			}
		case *ast.Ident:
			fs.Type = fmt.Sprintf("*%s", starExpr.Name)
		}

		fs.Ptr = true
	case *ast.SelectorExpr: //  Xxx yyy.Zzz `k:"v"`
		if sx, oke := expr.X.(*ast.Ident); oke {
			fs.Type = fmt.Sprintf("%s.%s", sx.Name, expr.Sel.Name)
		} else {
			fs.Type = expr.Sel.Name
		}
	}
}

func (psr parser) handleFieldTag(field *ast.Field, fs *astx.FieldSpec) {
	if fieldTag := field.Tag; fieldTag != nil {
		tagValue := fieldTag.Value               // `xxx:"xv" yyy:"yv"`
		tagValue = tagValue[1 : len(tagValue)-1] // xxx:"xv" yyy:"yv"
		ts := &astx.TagSpec{
			Field: fs.Name,
			Tags:  make([]*astx.Tag, 0),
		}
		tvs := strings.Split(tagValue, " ")
		for _, tv := range tvs {
			kvs := strings.Split(tv, ":")
			k := kvs[0]                    // xxx | yyy
			v := kvs[1][1 : len(kvs[1])-1] // xv | yv
			tag := &astx.Tag{
				Name:  fieldTag.Value,
				Key:   k,
				Value: v,
			}
			ts.Tags = append(ts.Tags, tag)
		}

		fs.Tags = append(fs.Tags, ts)
	}
}

func (psr parser) handleResults(funcDecl *ast.FuncDecl, ps *astx.PackageSpec, ms *astx.MethodSpec) {
	hasResults := funcDecl.Type != nil && funcDecl.Type.Results != nil && funcDecl.Type.Results.List != nil
	if hasResults {
		for _, rvt := range funcDecl.Type.Results.List {
			rs := &astx.ReturnSpec{
				Pkg:      ps.Pkg,
				FuncName: funcDecl.Name.String(),
				Ptr:      false,
			}
			if names := rvt.Names; names != nil {
				rs.Name = rvt.Names[0].Name
			}

			switch rvtType := rvt.Type.(type) {
			case *ast.Ident:
				rs.Type = rvtType.Name
			case *ast.StarExpr:
				switch xt := rvtType.X.(type) {
				case *ast.Ident:
					rs.Type = fmt.Sprintf("*%s", xt.Name)
				case *ast.SelectorExpr:
					if x, okx := xt.X.(*ast.Ident); okx {
						rs.Type = fmt.Sprintf("*%s.%s", x.Name, xt.Sel.Name)
					} else {
						rs.Type = fmt.Sprintf("*%s", xt.Sel.Name)
					}
				}
				rs.Ptr = true
			}

			ms.Returns = append(ms.Returns, rs)
		}
	}
}

func (psr parser) handleParams(funcDecl *ast.FuncDecl, ps *astx.PackageSpec, ms *astx.MethodSpec) {
	hasParams := funcDecl.Type != nil && funcDecl.Type.Params != nil && funcDecl.Type.Params.List != nil
	if hasParams {
		for _, param := range funcDecl.Type.Params.List {
			for _, pn := range param.Names {
				pms := &astx.ParamSpec{
					Pkg:      ps.Pkg,
					FuncName: funcDecl.Name.String(),
					Name:     pn.Name,
					Ptr:      false,
				}

				switch expr := param.Type.(type) {
				case *ast.Ident:
					pms.Type = expr.Name
				case *ast.StarExpr:
					switch starExpr := expr.X.(type) {
					case *ast.Ident:
						pms.Type = starExpr.Name
					case *ast.SelectorExpr:
						if x, okx := starExpr.X.(*ast.Ident); okx {
							pt := fmt.Sprintf("*%s.%s", x.Name, starExpr.Sel.Name)
							pms.Type = pt
						} else {
							pt := fmt.Sprintf("*%s", starExpr.Sel.Name)
							pms.Type = pt
						}
					}
					pms.Ptr = true
				case *ast.SelectorExpr:
					if x, okx := expr.X.(*ast.Ident); okx {
						pms.Type = fmt.Sprintf("%s.%s", x.Name, expr.Sel.Name)
					} else {
						pms.Type = expr.Sel.Name
					}
				}

				ms.Params = append(ms.Params, pms)
			}
		}
	}
}

func NewParser() Parser {
	return &parser{}
}

func Parse(pkg *loader.Package) *astx.AstSpec {
	return _parser_.Parse(pkg)
}

func populatePackageSpec(aw *astx.Astx) *astx.PackageSpec {
	ps := &astx.PackageSpec{
		Pkg:        aw.Package.PkgPath,
		Alias:      aw.Package.Name,
		Structs:    make([]*astx.StructSpec, 0),
		Interfaces: make([]*astx.InterfaceSpec, 0),
		Funcs:      make([]*astx.FuncSpec, 0),
	}

	return ps
}
