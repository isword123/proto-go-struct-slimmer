package logic

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"reflect"
	"strings"
)

type ProtoGoParser struct {
	file *ast.File
	packageName string
}

func (pp *ProtoGoParser)Parse(filePath string) bool {
	fs := token.NewFileSet()
	parsedFile, err := parser.ParseFile(fs, filePath, nil, parser.ParseComments)

	if err != nil {
		log.Println("Parse proto go file failed", err)
		return false
	}

	pp.file = parsedFile

	return true
}

func (pp *ProtoGoParser)getPackageName() string {
	if len(pp.packageName) == 0 {
		return "hello"
	}

	return pp.packageName
}

func (pp *ProtoGoParser)PrintStructs() {
	bufs := new(bytes.Buffer)

	bufs.WriteString(fmt.Sprintf("package %s\n\n", pp.getPackageName()))

	for _, decl := range pp.file.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}

		for _, spec := range genDecl.Specs {
			tSpec, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}

			structExp, ok1 := tSpec.Type.(*ast.StructType)
			if !ok1 {
				continue
			}

			fmt.Println("struct", tSpec.Name)

			// 不是公开的数据结构，不处理
			if !ast.IsExported(tSpec.Name.Name) {
				continue
			}


			bufs.WriteString(fmt.Sprintf("type %s struct {\n", tSpec.Name))

			for _, field := range structExp.Fields.List {
				if len(field.Names) <= 0 {
					continue
				}

				fmt.Println("filed type---", reflect.TypeOf(field.Type))

				fieldName := field.Names[0].Name
				if strings.HasPrefix(fieldName, "XXX_") {
					continue
				}

				if ident, ok := field.Type.(*ast.Ident); ok {
					fmt.Println("names", field.Names, "type", ident.Name, "tag", field.Tag)
					bufs.WriteString(fmt.Sprintf("\t%s %s\n", field.Names[0].Name, ident.Name))
				}

				fmt.Println("filed type---", reflect.TypeOf(field.Type), "field name ---", field.Names[0].Name)

				// *ast.ArrayType field name --- Titles
				if arrI, ok := field.Type.(*ast.ArrayType); ok {
					if eleI, ok := arrI.Elt.(*ast.StarExpr); ok {
						if detailI, ok := eleI.X.(*ast.Ident); ok {
							bufs.WriteString(fmt.Sprintf("\t%s []*%s\n", field.Names[0].Name, detailI.Name))
						} else {
							fmt.Println("wrong identifier type", field.Names[0].Name, reflect.TypeOf(eleI.X))
						}
					} else if eleI, ok := arrI.Elt.(*ast.Ident); ok {
						bufs.WriteString(fmt.Sprintf("\t%s []%s\n", field.Names[0].Name, eleI.Name))
					} else {
						fmt.Println("wrong identifier type", field.Names[0].Name, reflect.TypeOf(arrI.Elt))
						continue
					}

				}
			}

			bufs.WriteString("}\n\n")

			fmt.Println("------")
		}
	}

	fmt.Println("structs\n\n", bufs.String())
}