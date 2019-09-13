package logic

import (
	"bytes"
	"fmt"
	"strings"
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

func (o *Object)BuildTransMethod() []byte {
	buff := new(bytes.Buffer)

	funcName := fmt.Sprintf("Trans%s", o.Name)
	srcType := fmt.Sprintf("*%s.%s", srcPkgName, o.Name)

	// 单个
	buff.WriteString(fmt.Sprintf("func %s(src %s) *%s {\n", funcName, srcType, o.Name))

	buff.WriteString("\tif src == nil {\n")
	buff.WriteString("\t\treturn nil\n")
	buff.WriteString("\t}\n")

	buff.WriteString(fmt.Sprintf("\tvar dst %s\n", o.Name))

	for _, field := range o.Fields {
		buff.WriteString(fmt.Sprintf("\tdst.%s = %s\n",field.Name, field.BuildAssignSt("src")))
	}

	buff.WriteString(fmt.Sprintf("\treturn &dst\n"))
	buff.WriteString("}\n\n")

	// 数组
	funcArrName := funcName + "Arr"
	buff.WriteString(fmt.Sprintf("func %s(sources []%s) []*%s{\n", funcArrName, srcType, o.Name))
	buff.WriteString(fmt.Sprintf("\tvar dsts []*%s\n", o.Name))

	buff.WriteString("\tfor _, src := range sources {\n")
	buff.WriteString(fmt.Sprintf("\t\tdsts = append(dsts, %s(src))\n", funcName))
	buff.WriteString("\t}\n")

	buff.WriteString("\treturn dsts\n")
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

func (f *Field)BuildAssignSt(prefix string) []byte {
	buff := new(bytes.Buffer)

	firstTypeStr := f.Type[0:1]
	if firstTypeStr != strings.ToUpper(firstTypeStr) {
		// 如果是原始值类型
		buff.WriteString(fmt.Sprintf("%s.%s", prefix, f.Name))
	} else {
		// 如果是 struct
		if f.IsArr {
			funcName := "Trans" + f.Type + "Arr"
			buff.WriteString(fmt.Sprintf("%s(%s.%s)", funcName, prefix, f.Name))
		} else if f.IsPointer {
			funcName := "Trans" + f.Type
			buff.WriteString(fmt.Sprintf("%s(%s.%s)", funcName, prefix, f.Name))
		} else {
			funcName := "Trans" + f.Type
			buff.WriteString(fmt.Sprintf("%s(%s.%s)", funcName, prefix, f.Name))
		}
	}

	return buff.Bytes()
}

type Const struct {
	Name string
	Type string
	Values map[string]string
}

func (c *Const)AddVal(key, val string) {
	if c.Values == nil {
		c.Values = make(map[string]string)
	}
	c.Values[key] = val
}

func (c *Const)Export() []byte {
	buff := new(bytes.Buffer)

	buff.WriteString(fmt.Sprintf("type %s %s\n", c.Name, c.Type))

	if len(c.Values) > 0 {
		buff.WriteString("const (\n")
		for key, val := range c.Values {
			buff.WriteString(fmt.Sprintf("\t%s %s = %s\n", key, c.Name, val))
		}
		buff.WriteString(")\n")
	}

	return buff.Bytes()
}

func (c *Const)BuildTransMethod() []byte {
	buff := new(bytes.Buffer)

	funcName := "Trans" + c.Name
	srcType := srcPkgName + "." + c.Name

	buff.WriteString(fmt.Sprintf("func %s(src %s) %s {\n", funcName, srcType, c.Name))
	buff.WriteString(fmt.Sprintf("\treturn %s(src)\n", c.Name))
	buff.WriteString("}\n")

	return buff.Bytes()
}
