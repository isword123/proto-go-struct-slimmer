package main

import (
	"flag"
	"fmt"
	"github.com/isword123/proto-go-struct-slimmer/logic"
)

var (
	goFile = flag.String("f", "file path", "go file path")
	dstDir = flag.String("d", "dst file dir", "")
)

func main() {
	flag.Parse()

	if len(*goFile) == 0 || len(*dstDir) == 0 {
		fmt.Println("No go file path specified", *goFile, *dstDir)
		return
	}

	parser := new(logic.ProtoGoParser)
	ok := parser.ParseAndSave(*goFile, *dstDir)
	if !ok {
		fmt.Println("Parse and save failed")
		return
	}
}
