package logic

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
)

type ProtoGoParser struct {
	file *ast.File
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

func (pp *ProtoGoParser)PrintStructs() {
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

			for _, field := range structExp.Fields.List {
				fmt.Println(field.Names, field.Type, field.Tag)
			}

			fmt.Println("------")
		}

		//fmt.Println(decl)
	}
}