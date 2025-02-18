package main

import (
	//"flag"

	"fmt"
	"os"
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
	err := gen.Generate()
	if err != nil {
		fmt.Errorf("Failed to generate: %v", err)
		os.Exit(1)
	}
	os.Exit(0)
}
