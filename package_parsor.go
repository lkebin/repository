package main

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"strings"

	"github.com/sirupsen/logrus"
	"golang.org/x/tools/go/packages"
)

type RepositoryTypes struct {
	Name       string
	TypeParams []*types.TypeParam
	Methods    []*types.Func
}

func ParseTypes(tps, patterns []string, tags []string) []RepositoryTypes {
	cfg := &packages.Config{
		Mode:       packages.NeedName | packages.NeedTypes | packages.NeedTypesInfo | packages.NeedSyntax | packages.NeedImports | packages.NeedDeps,
		Tests:      false,
		BuildFlags: []string{fmt.Sprintf("-tags=%s", strings.Join(tags, " "))},
		Logf:       logrus.Infof,
	}
	pkgs, err := packages.Load(cfg, patterns...)
	if err != nil {
		logrus.Fatal(err)
	}
	if len(pkgs) != 1 {
		logrus.Fatalf("error: %d packages matching %v", len(pkgs), strings.Join(patterns, " "))
	}

	var rts []RepositoryTypes

	// ast.Print(pkgs[0].Fset, pkgs[0].Syntax[2])

	for _, t := range tps {
		for _, file := range pkgs[0].Syntax {
			ast.Inspect(file, func(n ast.Node) bool {
				decl, ok := n.(*ast.GenDecl)
				if !ok || decl.Tok != token.TYPE {
					return true
				}

				for _, spec := range decl.Specs {
					ts := spec.(*ast.TypeSpec)
					it, ok := ts.Type.(*ast.InterfaceType)
					if !ok {
						continue
					}

					if ts.Name.Name != t {
						continue
					}

					rt := RepositoryTypes{
						Name: ts.Name.Name,
					}

					if ts.TypeParams != nil {
						for _, tp := range ts.TypeParams.List {
							rt.TypeParams = append(rt.TypeParams, pkgs[0].TypesInfo.Defs[tp.Names[0]].Type().(*types.TypeParam))
						}
					}

					logrus.Debug(ts.Name)
					for _, field := range it.Methods.List {
						switch field.Type.(type) {
						case *ast.FuncType:
							logrus.Debug("\tFuncType", field.Names[0], "\t", "Params", field.Type.(*ast.FuncType).Params.List, "\t", "Results", field.Type.(*ast.FuncType).Results.List)
							logrus.Debug("\t\tTypeInfo", pkgs[0].TypesInfo.Defs[field.Names[0]])
							rt.Methods = append(rt.Methods, pkgs[0].TypesInfo.Defs[field.Names[0]].(*types.Func))
						case *ast.IndexExpr:
							switch field.Type.(*ast.IndexExpr).X.(type) {
							case *ast.SelectorExpr:
								logrus.Debug("\tIndexExpr", field.Type.(*ast.IndexExpr).X.(*ast.SelectorExpr).Sel)
								logrus.Debug("\t\tTypeInfo", pkgs[0].TypesInfo.Instances[field.Type.(*ast.IndexExpr).X.(*ast.SelectorExpr).Sel])
								rt.Methods = append(rt.Methods, parseTypeEmbedding(pkgs[0].TypesInfo.Instances[field.Type.(*ast.IndexExpr).X.(*ast.SelectorExpr).Sel])...)
							case *ast.Ident:
								logrus.Debug("\tIndexExpr", field.Type.(*ast.IndexExpr).X.(*ast.Ident).Name)
								logrus.Debug("\t\tTypeInfo", pkgs[0].TypesInfo.Instances[field.Type.(*ast.IndexExpr).X.(*ast.Ident)])
								rt.Methods = append(rt.Methods, parseTypeEmbedding(pkgs[0].TypesInfo.Instances[field.Type.(*ast.IndexExpr).X.(*ast.Ident)])...)
							}
						case *ast.IndexListExpr:
							switch field.Type.(*ast.IndexListExpr).X.(type) {
							case *ast.SelectorExpr:
								logrus.Debug("\tIndexListExpr", field.Type.(*ast.IndexListExpr).X.(*ast.SelectorExpr).X, field.Type.(*ast.IndexListExpr).X.(*ast.SelectorExpr).Sel)
								logrus.Debug("\t\tTypesInfo", pkgs[0].TypesInfo.Instances[field.Type.(*ast.IndexListExpr).X.(*ast.SelectorExpr).Sel])
								rt.Methods = append(rt.Methods, parseTypeEmbedding(pkgs[0].TypesInfo.Instances[field.Type.(*ast.IndexListExpr).X.(*ast.SelectorExpr).Sel])...)
							case *ast.Ident:
								logrus.Debug("\tIndexListExpr", field.Type.(*ast.IndexListExpr).X.(*ast.Ident).Name)
								logrus.Debug("\t\tTypesInfo", pkgs[0].TypesInfo.Instances[field.Type.(*ast.IndexListExpr).X.(*ast.Ident)])
								rt.Methods = append(rt.Methods, parseTypeEmbedding(pkgs[0].TypesInfo.Instances[field.Type.(*ast.IndexListExpr).X.(*ast.Ident)])...)
							}
						}
					}

					rts = append(rts, rt)
				}

				return true
			})
		}

		for _, v := range rts[0].Methods {
			for i := 0; i < v.Type().(*types.Signature).Params().Len(); i++ {
				fmt.Println(v.Name(), types.TypeString(v.Type().(*types.Signature), nil))
			}
		}
	}

	return rts
}

func parseTypeEmbedding(instance types.Instance) []*types.Func {
	return parseTypeInterface(instance.Type)
}

func parseTypeInterface(t types.Type) []*types.Func {
	var funcs []*types.Func
	switch tt := t.(type) {
	case *types.Named:
		return parseTypeInterface(tt.Underlying())
	case *types.Interface:
		for i := 0; i < tt.NumMethods(); i++ {
			logrus.Debug("\tMethod", tt.Method(i).Name(), tt.Method(i).Pkg())
			logrus.Debug("\t\tParams", types.TypeString(tt.Method(i).Type().(*types.Signature).Params(), nil))
			for p := 0; p < tt.Method(i).Type().(*types.Signature).Params().Len(); p++ {
				logrus.Debug("\t\t\tParams type:", tt.Method(i).Type().(*types.Signature).Params().At(p).Type())
			}

			funcs = append(funcs, tt.Method(i))
		}
	}
	return funcs
}
