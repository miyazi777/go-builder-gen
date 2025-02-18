package main

import (
	//"flag"
	"fmt"
	"os"
	"path/filepath"
	"test1/analyze"
	"test1/generator"
)

func main() {
	//	structName := flag.String("struct", "", "target struct name")
	//	sourceFilePath := flag.String("path", "", "source file path")
	//
	//	flag.Parse()
	//
	//	fmt.Printf("Target struct: %s\n", *structName)
	//	fmt.Printf("Source path: %s\n", *sourceFilePath)

	filePath := "./model/item.go"
	structName := "Item"
	gs := analyze.AnalyzeStruct(filePath, structName)

	fmt.Printf("PackageName: %s\n", gs.PackageName)
	fmt.Printf("StructName: %s\n", gs.StructName)
	for _, field := range gs.GenFields {
		fmt.Printf("Field: %s, Type: %s\n", field.Name, field.Type)
	}

	gen(gs)
}

func gen(gs generator.GenStruct) {
	f := gs.Generate()
	outputDir := "./output"
	outputFile := filepath.Join(outputDir, "item_builder_generated_test.go")

	err := os.MkdirAll(outputDir, 0755)
	if err != nil {
		fmt.Printf("Failed to create directory: %v\n", err)
		return
	}

	if err := f.Save(outputFile); err != nil {
		fmt.Printf("Failed to save file: %v\n", err)
		return
	}
}
