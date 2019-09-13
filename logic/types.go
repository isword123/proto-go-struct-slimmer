package logic

import (
	"bytes"
	"fmt"
)

type Object struct {
	Name string
	Fields []Field
}

func (o *Object)Export() []byte {
	buff := new(bytes.Buffer)

	buff.WriteString(fmt.Sprintf("type %s struct {\n", o.Name))

	for _, field := range o.Fields {
		buff.Write(field.Export())
	}

	buff.WriteString("}\n\n")

	return buff.Bytes()
}

type Field struct {
	Name string
	Type string
	IsPointer bool
	IsArr bool
	IsArrSubPointer bool
	Tag string
}

func (f *Field)Export() []byte {
	buff := new(bytes.Buffer)

	buff.WriteString("\t")

	if f.IsArr {
		eleT := ""
		if f.IsArrSubPointer {
			eleT = "*" + f.Type
		} else {
			eleT = f.Type
		}

		buff.WriteString(fmt.Sprintf("%s []%s", f.Name, eleT))
	} else if f.IsPointer {
		buff.WriteString(fmt.Sprintf("%s *%s", f.Name, f.Type))
	} else {
		buff.WriteString(fmt.Sprintf("%s %s", f.Name, f.Type))
	}

	if len(f.Tag) > 0 {
		buff.WriteString(fmt.Sprintf(" `%s`", f.Tag))
	}

	buff.WriteString("\n")

	return buff.Bytes()
}

func BuildTransMethod() {

}
