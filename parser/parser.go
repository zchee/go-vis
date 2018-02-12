// Copyright 2018 The go-vis Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package parser

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"os"
	"strconv"
	"strings"

	dotast "gonum.org/v1/gonum/graph/formats/dot/ast"
)

type namedType struct {
	Ident *ast.Ident
	Type  ast.Expr
}

type pkgTypes map[string]map[string]namedType

func InspectDir(path string) (pkgTypes, error) {
	ptypes := make(pkgTypes)

	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, path, func(n os.FileInfo) bool { return true }, 0)
	if err != nil {
		return nil, err
	}

	for name, pkg := range pkgs {
		ptypes[name] = make(map[string]namedType)
		for _, f := range pkg.Files {
			ast.Inspect(f, func(n ast.Node) bool {
				switch n := n.(type) {
				case *ast.CommentGroup, *ast.Comment: // skip comment
					return false
				case *ast.TypeSpec:
					ptypes[name][n.Name.Name] = namedType{
						Ident: n.Name,
						Type:  n.Type,
					}
				}
				return true
			})
		}
	}

	return ptypes, nil
}

func makeLabel(vals ...string) *dotast.Attr {
	return &dotast.Attr{
		Key: "label",
		Val: "\"{" + escape(strings.Join(vals, " ")) + "}\"",
	}
}

func makeInterface(builder strings.Builder, name string, t *ast.InterfaceType) *dotast.Attr {
	builder.WriteString(fmt.Sprintf("%s interface|", name))
	for i, f := range t.Methods.List {
		if i > 0 {
			fmt.Fprintf(&builder, `|`)
		}
		fmt.Fprintf(&builder, `<f%d>`, i)
		// a,b,c Type
		for ii, n := range f.Names {
			fmt.Fprintf(&builder, "%s", n.Name)
			if ii > 0 {
				fmt.Fprintf(&builder, `,`)
			}
		}
		if len(f.Names) > 0 {
			fmt.Fprintf(&builder, ` `)
		}
		fmt.Fprintf(&builder, `%s`, toString(f.Type))
	}

	return &dotast.Attr{
		Key: "label",
		Val: "\"{" + escape(builder.String()) + "}\"",
	}
}

func makeStruct(builder strings.Builder, name string, t *ast.StructType) *dotast.Attr {
	builder.WriteString(fmt.Sprintf("%s|", name))
	for i, f := range t.Fields.List {
		if i > 0 {
			fmt.Fprintf(&builder, `|`)
		}
		fmt.Fprintf(&builder, `<f%d>`, i)
		// a,b,c Type
		for ii, n := range f.Names {
			fmt.Fprintf(&builder, "%s", n.Name)
			if ii > 0 {
				fmt.Fprintf(&builder, `,`)
			}
		}
		if len(f.Names) > 0 {
			fmt.Fprintf(&builder, ` `)
		}
		fmt.Fprintf(&builder, `%s`, toString(f.Type))
	}

	return &dotast.Attr{
		Key: "label",
		Val: "\"{" + escape(builder.String()) + "}\"",
	}
}

// func makePort(name string)

func Render(out io.Writer, ptypes pkgTypes) {
	graph := &dotast.Graph{
		ID:       strconv.Quote("go-vis"),
		Directed: true,
	}

	var builder strings.Builder
	for pkg, typs := range ptypes {
		if strings.HasSuffix(pkg, "_test") {
			continue
		}

		subgraph := &dotast.Subgraph{
			ID:    pkg,
			Stmts: []dotast.Stmt{makeLabel(pkg)},
		}

		stmt := []dotast.Stmt{}
		for _, typ := range typs {
			builder.Reset()

			switch t := typ.Type.(type) {
			case *ast.Ident:
				attrLabel := makeLabel(typ.Ident.Name, t.Name)
				attr := &dotast.Attr{Key: "shape", Val: "ellipse"}
				nodeStmt := &dotast.NodeStmt{
					Node: &dotast.Node{
						ID: typ.Ident.Name,
					},
					Attrs: []*dotast.Attr{attr, attrLabel},
				}
				stmt = append(stmt, nodeStmt)
			case *ast.SelectorExpr:
				attrLabel := makeLabel(typ.Ident.Name, toString(t))
				attr := &dotast.Attr{Key: "shape", Val: "ellipse"}
				nodeStmt := &dotast.NodeStmt{
					Node: &dotast.Node{
						ID: typ.Ident.Name,
					},
					Attrs: []*dotast.Attr{attr, attrLabel},
				}
				stmt = append(stmt, nodeStmt)
			case *ast.ChanType:
				attrLabel := makeLabel(typ.Ident.Name, toString(t))
				attr := &dotast.Attr{Key: "shape", Val: "box"}
				nodeStmt := &dotast.NodeStmt{
					Node: &dotast.Node{
						ID: typ.Ident.Name,
					},
					Attrs: []*dotast.Attr{attr, attrLabel},
				}
				stmt = append(stmt, nodeStmt)
			case *ast.FuncType:
				attrLabel := makeLabel(typ.Ident.Name, toString(t))
				attr := &dotast.Attr{Key: "shape", Val: "rectangle"}
				nodeStmt := &dotast.NodeStmt{
					Node: &dotast.Node{
						ID: typ.Ident.Name,
					},
					Attrs: []*dotast.Attr{attr, attrLabel},
				}
				stmt = append(stmt, nodeStmt)
			case *ast.ArrayType:
				attrLabel := makeLabel(typ.Ident.Name, toString(t))
				attr := &dotast.Attr{Key: "shape", Val: "rectangle"}
				nodeStmt := &dotast.NodeStmt{
					Node: &dotast.Node{
						ID: typ.Ident.Name,
					},
					Attrs: []*dotast.Attr{attr, attrLabel},
				}
				stmt = append(stmt, nodeStmt)
			case *ast.MapType:
				attrLabel := makeLabel(typ.Ident.Name, toString(t))
				attr := &dotast.Attr{Key: "shape", Val: "rectangle"}
				nodeStmt := &dotast.NodeStmt{
					Node: &dotast.Node{
						ID: typ.Ident.Name,
					},
					Attrs: []*dotast.Attr{attr, attrLabel},
				}
				stmt = append(stmt, nodeStmt)
			case *ast.InterfaceType:
				interfaceLabel := makeInterface(builder, typ.Ident.Name, t)
				attr := &dotast.Attr{Key: "shape", Val: "Mrecord"}
				nodeStmt := &dotast.NodeStmt{
					Node: &dotast.Node{
						ID: typ.Ident.Name,
					},
					Attrs: []*dotast.Attr{attr, interfaceLabel},
				}
				stmt = append(stmt, nodeStmt)
			case *ast.StructType:
				structLabel := makeStruct(builder, typ.Ident.Name, t)
				attr := &dotast.Attr{Key: "shape", Val: "record"}
				nodeStmt := &dotast.NodeStmt{
					Node: &dotast.Node{
						ID: typ.Ident.Name,
					},
					Attrs: []*dotast.Attr{attr, structLabel},
				}
				stmt = append(stmt, nodeStmt)
			}
		}

		for _, etype := range typs {
			switch t := etype.Type.(type) {
			case *ast.FuncType:
				for i, typ := range dependsOn(t) {
					if _, ok := typs[typ]; ok {
						from := &dotast.Node{
							ID: strconv.Quote(etype.Ident.Name),
							Port: &dotast.Port{
								ID: fmt.Sprintf("f%d", i),
							},
						}
						to := &dotast.Edge{
							Directed: true,
							Vertex: &dotast.Node{
								ID: strconv.Quote(typ),
								Port: &dotast.Port{
									ID: fmt.Sprintf("f%d", i),
								},
							},
						}
						edgeStmt := &dotast.EdgeStmt{
							From: from,
							To:   to,
						}
						stmt = append(stmt, edgeStmt)
					}
				}
			case *ast.ChanType:
				for i, typ := range dependsOn(t) {
					if _, ok := typs[typ]; ok {
						from := &dotast.Node{
							ID: strconv.Quote(etype.Ident.Name),
							Port: &dotast.Port{
								ID: fmt.Sprintf("f%d", i),
							},
						}
						to := &dotast.Edge{
							Directed: true,
							Vertex: &dotast.Node{
								ID: strconv.Quote(typ),
							},
						}
						edgeStmt := &dotast.EdgeStmt{
							From: from,
							To:   to,
						}
						stmt = append(stmt, edgeStmt)
					}
				}
			case *ast.InterfaceType:
				for i, f := range t.Methods.List {
					from := &dotast.Node{
						ID: strconv.Quote(etype.Ident.Name),
						Port: &dotast.Port{
							ID: fmt.Sprintf("f%d", i),
						},
					}
					for _, typ := range dependsOn(f.Type) {
						if _, ok := typs[typ]; ok {
							to := &dotast.Edge{
								Directed: true,
								Vertex: &dotast.Node{
									ID: strconv.Quote(typ),
								},
							}
							edgeStmt := &dotast.EdgeStmt{
								From: from,
								To:   to,
							}
							stmt = append(stmt, edgeStmt)
						}
					}
				}
			case *ast.StructType:
				for i, f := range t.Fields.List {
					from := &dotast.Node{
						ID: strconv.Quote(etype.Ident.Name),
						Port: &dotast.Port{
							ID: fmt.Sprintf("f%d", i),
						},
					}
					for _, typ := range dependsOn(f.Type) {
						if _, ok := typs[typ]; ok {
							to := &dotast.Edge{
								Directed: true,
								Vertex: &dotast.Node{
									ID: strconv.Quote(typ),
								},
							}
							edgeStmt := &dotast.EdgeStmt{
								From: from,
								To:   to,
							}
							stmt = append(stmt, edgeStmt)
						}
					}
				}

			}
		}

		subgraph.Stmts = append(subgraph.Stmts, stmt...)
		graph.Stmts = append(graph.Stmts, subgraph)
	}

	file := &dotast.File{
		Graphs: []*dotast.Graph{graph},
	}
	fmt.Fprintf(out, file.String())
}

func escape(s string) string {
	for _, ch := range " '`[]{}()*" {
		s = strings.Replace(s, string(ch), `\`+string(ch), -1)
	}

	return s
}

func toString(n interface{}) string {
	switch t := n.(type) {
	case nil:
		return "nil"
	case *ast.Ident:
		return t.Name
	case *ast.SelectorExpr:
		return toString(t.X) + "." + toString(t.Sel)
	case *ast.Object:
		return t.Name
	case *ast.StarExpr:
		return `*` + toString(t.X)
	case *ast.InterfaceType:
		// TODO:
		return `interface{}`
	case *ast.MapType:
		return `map[` + toString(t.Key) + `]` + toString(t.Value)
	case *ast.ChanType:
		return `chan ` + toString(t.Value)
	case *ast.StructType:
		// TODO:
		return `struct {}` //+ toString(t.)
	case *ast.Ellipsis:
		return `...` + toString(t.Elt)
	case *ast.Field:
		// ignoring names
		return toString(t.Type)

	case *ast.FuncType:
		var buf bytes.Buffer
		fmt.Fprint(&buf, `func(`)
		if t.Params != nil && len(t.Params.List) > 0 {
			for i, p := range t.Params.List {
				if i > 0 {
					fmt.Fprint(&buf, `, `)
				}
				fmt.Fprint(&buf, toString(p))
			}
		}
		fmt.Fprint(&buf, `)`)

		if t.Results != nil && len(t.Results.List) > 0 {
			fmt.Fprint(&buf, ` (`)
			for i, r := range t.Results.List {
				if i > 0 {
					fmt.Fprint(&buf, `, `)
				}
				fmt.Fprint(&buf, toString(r))
			}
			fmt.Fprint(&buf, `)`)
		}

		return buf.String()
	case *ast.ArrayType:
		return `[]` + toString(t.Elt)
	default:
		return fmt.Sprintf("%#v", n)
	}
}

// dependsOn collect all the type names node n depends on
func dependsOn(n interface{}) []string {
	switch t := n.(type) {
	case nil:
		return nil
	case *ast.Ident:
		return []string{t.Name}
	case *ast.SelectorExpr:
		return []string{toString(t.X) + "." + t.Sel.Name}
	case *ast.Object:
		return []string{t.Name}
	case *ast.Field:
		return dependsOn(t.Type)
	case *ast.StarExpr:
		return dependsOn(t.X)
	case *ast.MapType:
		return append(dependsOn(t.Key), dependsOn(t.Value)...)
	case *ast.ChanType:
		return dependsOn(t.Value)
	case *ast.InterfaceType:
		if t.Methods == nil {
			return nil
		}
		var types []string
		for _, v := range t.Methods.List {
			types = append(types, dependsOn(v.Type)...)
		}
		return types
	case *ast.StructType:
		var types []string
		for _, v := range t.Fields.List {
			types = append(types, dependsOn(v.Type)...)
		}
		return types
	case *ast.FuncType:
		var types []string

		if t.Params != nil {
			for _, v := range t.Params.List {
				types = append(types, dependsOn(v.Type)...)
			}
		}

		if t.Results != nil {
			for _, v := range t.Results.List {
				types = append(types, dependsOn(v.Type)...)
			}
		}

		return types

	case *ast.ArrayType:
		return dependsOn(t.Elt)
	default:
		return []string{fmt.Sprintf("%#v", n)}
	}
}
