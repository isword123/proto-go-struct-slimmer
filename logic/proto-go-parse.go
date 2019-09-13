package logic

import (
	"bytes"
	"fmt"
	"github.com/isword123/proto-go-struct-slimmer/models"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

type ProtoGoParser struct {
	file *ast.File
	packageName string
	fileBaseName string
	modName string
}

func (pp *ProtoGoParser)Parse(filePath string) bool {
	fs := token.NewFileSet()
	parsedFile, err := parser.ParseFile(fs, filePath, nil, parser.ParseComments)

	if err != nil {
		log.Println("Parse proto go file failed", err)
		return false
	}

	pp.file = parsedFile

	pp.fileBaseName = strings.TrimSuffix(filepath.Base(filePath), ".pb.go")
	fmt.Println("file base name", filePath, pp.fileBaseName)

	pp.packageName = parsedFile.Name.Name + "_trans"

	return true
}

func (pp *ProtoGoParser)getPackageName() string {
	if len(pp.packageName) == 0 {
		return "hello"
	}

	return pp.packageName
}

func (pp *ProtoGoParser) getFileBaseName() string {
	return pp.fileBaseName
}

func (pp *ProtoGoParser) GetStructsBytes() []byte {
	bufs := new(bytes.Buffer)

	bufs.WriteString(fmt.Sprintf("package %s\n\n", pp.getPackageName()))

	for _, decl := range pp.file.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}

		for _, spec := range genDecl.Specs {

			// 处理 const 字段
			if vSpec, ok := spec.(*ast.ValueSpec); ok && genDecl.Tok == token.CONST {
				pp.getConstDefs(vSpec, bufs)
				continue
			}

			tSpec, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}

			// for type definitions
			if it, ok := tSpec.Type.(*ast.Ident); ok {
				pp.getConstTypeDefs(tSpec ,it, bufs)
				continue
			}

			structExp, ok1 := tSpec.Type.(*ast.StructType)
			if !ok1 {
				continue
			}

			// 不是公开的数据结构，不处理
			if !ast.IsExported(tSpec.Name.Name) {
				continue
			}

			structName := tSpec.Name.Name
			bufs.WriteString(fmt.Sprintf("type %s struct {\n", structName))

			for _, field := range structExp.Fields.List {
				if len(field.Names) <= 0 {
					continue
				}

				fieldSet := false
				fieldName := field.Names[0].Name
				if strings.HasPrefix(fieldName, "XXX_") {
					continue
				}

				if ident, ok := field.Type.(*ast.Ident); ok {
					if !models.IsExcludeInDasAgDota2(structName, fieldName) {
						bufs.WriteString(fmt.Sprintf("\t%s %s", fieldName, ident.Name))
						fieldSet = true
					}
				} else if arrI, ok := field.Type.(*ast.ArrayType); ok {
					// *ast.ArrayType field name --- Titles
					if eleI, ok := arrI.Elt.(*ast.StarExpr); ok {
						if detailI, ok := eleI.X.(*ast.Ident); ok {
							bufs.WriteString(fmt.Sprintf("\t%s []*%s", fieldName, detailI.Name))
							fieldSet = true
						} else {
							fmt.Println("wrong identifier type", fieldName, reflect.TypeOf(eleI.X))
						}
					} else if eleI, ok := arrI.Elt.(*ast.Ident); ok {
						bufs.WriteString(fmt.Sprintf("\t%s []%s", fieldName, eleI.Name))
						fieldSet = true
					} else {
						fmt.Println("wrong identifier type", fieldName, reflect.TypeOf(arrI.Elt))
						continue
					}

				} else if starI, ok := field.Type.(*ast.StarExpr); ok {
					detailI, ok := starI.X.(*ast.Ident)
					if ok {
						if !models.IsExcludeInDasAgDota2(structName, fieldName) {
							bufs.WriteString(fmt.Sprintf("\t%s *%s", fieldName, detailI.Name))
							fieldSet = true
						}
					} else {
						fmt.Println("Unknown star type", fieldName, starI.X)
					}
				} else {
					fmt.Println("Unknown type", fieldName, field.Type)
				}

				if !fieldSet {
					continue
				}

				if field.Tag != nil {
					jsonTag, ok := pp.parseJSONTag(field.Tag.Value)
					if ok {
						bufs.WriteString(fmt.Sprintf(" `%s`\n", jsonTag))
					}
				} else {
					bufs.WriteString("\n")
				}
			}

			bufs.WriteString("}\n\n")
		}
	}

	return bufs.Bytes()
}

func (pp *ProtoGoParser) parseJSONTag(srcTag string) (string, bool) {
	index := strings.Index(srcTag, "json:")

	if index < 0 {
		return "", false
	}

	return srcTag[index:len(srcTag) - 1], true
}

func (pp *ProtoGoParser) saveNewCode(bs []byte, dir string) bool {
	fileName := filepath.Join(dir, fmt.Sprintf("%s.go", pp.getFileBaseName()))

	err := ioutil.WriteFile(fileName, bs, os.ModePerm)
	if err != nil {
		fmt.Println("Save new code failed", fileName, err.Error())
		return false
	}

	return true
}

func (pp *ProtoGoParser) getConstTypeDefs(ts *ast.TypeSpec,ident *ast.Ident, buff *bytes.Buffer) {
	buff.WriteString(fmt.Sprintf("type %s %s \n\n", ts.Name, ident.Name))
}

func (pp *ProtoGoParser) getConstDefs(vs *ast.ValueSpec, buff *bytes.Buffer) {
	typ := ""

	if vs.Type == nil && len(vs.Values) > 0 {
		ce, ok := vs.Values[0].(*ast.CallExpr)
		if !ok {
			return
		}
		id, ok := ce.Fun.(*ast.Ident)
		if !ok {
			return
		}
		typ = id.Name
	} else if vs.Type != nil {
		ident, ok := vs.Type.(*ast.Ident)
		if !ok {
			return
		}

		typ = ident.Name
	}

	if len(typ) == 0 {
		return
	}

	if len(vs.Names) == 0 || len(vs.Values) == 0 {
		return
	}

	name := vs.Names[0]
	val, ok := vs.Values[0].(*ast.BasicLit)
	if !ok {
		return
	}

	buff.WriteString(fmt.Sprintf("const %s %s = %s\n\n", name.Name, typ, val.Value))
}

func (pp *ProtoGoParser) ParseAndSave(filePath string, dir string) bool {
	pp.Parse(filePath)
	bs := pp.GetStructsBytes()
	return pp.saveNewCode(bs, dir)
}
