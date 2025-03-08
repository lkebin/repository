package repository

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"strings"

	"github.com/sirupsen/logrus"
	"golang.org/x/tools/go/packages"
)

type RepositorySpecs struct {
	Name             string
	Pkg              *types.Package
	TypeParams       []*types.TypeParam
	EmbeddedTypeArgs []types.Type
	Methods          []*types.Func
}

type ModelSpecs struct {
	Name       string
	Pkg        *types.Package
	TypeParams []*types.TypeParam
	Struct     *types.Struct
	Type       types.Type
}

func ParseModel(tps, patterns []string, tags []string) []ModelSpecs {
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

	var rts []ModelSpecs

	for _, t := range tps {
		for _, file := range pkgs[0].Syntax {
			ast.Inspect(file, func(n ast.Node) bool {
				decl, ok := n.(*ast.GenDecl)
				if !ok || decl.Tok != token.TYPE {
					return true
				}

				for _, spec := range decl.Specs {
					ts := spec.(*ast.TypeSpec)
					if ts.Name.Name != t {
						continue
					}

					rt := ModelSpecs{
						Name: ts.Name.Name,
						Pkg:  pkgs[0].Types,
					}

					if ts.TypeParams != nil {
						for _, tp := range ts.TypeParams.List {
							rt.TypeParams = append(rt.TypeParams, pkgs[0].TypesInfo.Defs[tp.Names[0]].Type().(*types.TypeParam))
						}
					}

					logrus.Debug(ts.Name)
					rt.Struct = pkgs[0].TypesInfo.Defs[ts.Name].Type().Underlying().(*types.Struct)
					rt.Type = pkgs[0].TypesInfo.Defs[ts.Name].Type()

					rts = append(rts, rt)
				}

				return true
			})
		}
	}

	return rts
}

func ParseRepository(tps, patterns []string, tags []string) []RepositorySpecs {
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

	var rts []RepositorySpecs

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

					rt := RepositorySpecs{
						Name: ts.Name.Name,
						Pkg:  pkgs[0].Types,
					}

					if ts.TypeParams != nil {
						for _, tp := range ts.TypeParams.List {
							rt.TypeParams = append(rt.TypeParams, pkgs[0].TypesInfo.Defs[tp.Names[0]].Type().(*types.TypeParam))
						}
					}

					interfaceType := pkgs[0].TypesInfo.Defs[ts.Name].Type().(*types.Named).Underlying().(*types.Interface)
					for i := 0; i < interfaceType.NumEmbeddeds(); i++ {
						ta := interfaceType.EmbeddedType(i).(*types.Named).TypeArgs()
						for ii := 0; ii < ta.Len(); ii++ {
							rt.EmbeddedTypeArgs = append(rt.EmbeddedTypeArgs, ta.At(ii))
						}
					}

					logrus.Debug(ts.Name)
					for _, field := range it.Methods.List {
						switch field.Type.(type) {
						case *ast.FuncType:
							logrus.Debug("FuncType: ", field.Names[0])
							logrus.Debug("  TypesInfo: ", pkgs[0].TypesInfo.Defs[field.Names[0]])
							rt.Methods = append(rt.Methods, pkgs[0].TypesInfo.Defs[field.Names[0]].(*types.Func))
						case *ast.IndexExpr:
							switch field.Type.(*ast.IndexExpr).X.(type) {
							case *ast.SelectorExpr:
								logrus.Debug("IndexExpr: ", field.Type.(*ast.IndexExpr).X.(*ast.SelectorExpr).Sel)
								logrus.Debug("  TypeInfo: ", pkgs[0].TypesInfo.Instances[field.Type.(*ast.IndexExpr).X.(*ast.SelectorExpr).Sel].Type.String())
								rt.Methods = append(rt.Methods, parseTypeEmbedding(pkgs[0].TypesInfo.Instances[field.Type.(*ast.IndexExpr).X.(*ast.SelectorExpr).Sel])...)
							case *ast.Ident:
								logrus.Debug("IndexExpr: ", field.Type.(*ast.IndexExpr).X.(*ast.Ident).Name)
								logrus.Debug("  TypesInfo: ", pkgs[0].TypesInfo.Instances[field.Type.(*ast.IndexExpr).X.(*ast.Ident)].Type.String())
								rt.Methods = append(rt.Methods, parseTypeEmbedding(pkgs[0].TypesInfo.Instances[field.Type.(*ast.IndexExpr).X.(*ast.Ident)])...)
							}
						case *ast.IndexListExpr:
							switch field.Type.(*ast.IndexListExpr).X.(type) {
							case *ast.SelectorExpr:
								logrus.Debug("IndexListExpr: ", field.Type.(*ast.IndexListExpr).X.(*ast.SelectorExpr).X, field.Type.(*ast.IndexListExpr).X.(*ast.SelectorExpr).Sel)
								logrus.Debug("  TypesInfo: ", pkgs[0].TypesInfo.Instances[field.Type.(*ast.IndexListExpr).X.(*ast.SelectorExpr).Sel].Type.String())
								rt.Methods = append(rt.Methods, parseTypeEmbedding(pkgs[0].TypesInfo.Instances[field.Type.(*ast.IndexListExpr).X.(*ast.SelectorExpr).Sel])...)
							case *ast.Ident:
								logrus.Debug("IndexListExpr: ", field.Type.(*ast.IndexListExpr).X.(*ast.Ident).Name)
								logrus.Debug("  TypesInfo: ", pkgs[0].TypesInfo.Instances[field.Type.(*ast.IndexListExpr).X.(*ast.Ident)].Type.String())
								rt.Methods = append(rt.Methods, parseTypeEmbedding(pkgs[0].TypesInfo.Instances[field.Type.(*ast.IndexListExpr).X.(*ast.Ident)])...)
							}
						}
					}

					rts = append(rts, rt)
				}

				return true
			})
		}
	}

	return rts
}

func parseTypeEmbedding(instance types.Instance) []*types.Func {
	return parseTypeInterface(instance.Type)
}

// Test
func parseTypeInterface(t types.Type) []*types.Func {
	var funcs []*types.Func
	switch tt := t.(type) {
	case *types.Named:
		return parseTypeInterface(tt.Underlying())
	case *types.Interface:
		for i := 0; i < tt.NumMethods(); i++ {
			logrus.Debug("  Method: ", tt.Method(i).Name(), " ", tt.Method(i).Pkg().Path())
			logrus.Debug("    Params: ", types.TypeString(tt.Method(i).Type().(*types.Signature).Params(), nil))
			for p := 0; p < tt.Method(i).Type().(*types.Signature).Params().Len(); p++ {
				logrus.Debug("      Param type: ", tt.Method(i).Type().(*types.Signature).Params().At(p).Type())
			}
			logrus.Debug("    Result: ", types.TypeString(tt.Method(i).Type().(*types.Signature).Results(), nil))
			for p := 0; p < tt.Method(i).Type().(*types.Signature).Results().Len(); p++ {
				logrus.Debug("      Result type: ", tt.Method(i).Type().(*types.Signature).Results().At(p).Type())
			}

			funcs = append(funcs, tt.Method(i))
		}
	}
	return funcs
}
