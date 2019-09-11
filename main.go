package main

import (
	"flag"
	"fmt"
	"github.com/isword123/proto-go-struct-slimmer/logic"
)

var (
	goFile = flag.String("f", "file path", "go file path")
)

func main() {
	flag.Parse()

	fmt.Println("Hello, proto-go-struct-slimmer")
	parser := new(logic.ProtoGoParser)
	ok := parser.Parse(*goFile)
	if !ok {
		return
	}

	parser.PrintStructs()
}
