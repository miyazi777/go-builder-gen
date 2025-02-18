package main

import (
	//"flag"

	"test1/generator"
)

func main() {
	// TODO: 引数を受け取る
	//	structName := flag.String("struct", "", "target struct name")
	//	sourceFilePath := flag.String("path", "", "source file path")
	//
	//	flag.Parse()
	//
	//	fmt.Printf("Target struct: %s\n", *structName)
	//	fmt.Printf("Source path: %s\n", *sourceFilePath)

	gen := generator.NewGenerator("./model/item.go", "Item")
	gen.Generate()
}
