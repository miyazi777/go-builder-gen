package main

import (
	//"flag"

	"flag"
	"fmt"
	"os"
	"test1/file"
	"test1/generator"
)

func main() {
	structName := flag.String("struct", "", "target struct name")
	sourceFilePath := flag.String("path", "", "source file path")

	flag.Parse()

	// check builder file
	builderFile := file.NewBuilderFile(*sourceFilePath)
	if builderFile.IsExistFile() {
		if err := builderFile.IsInvalidComment(); err != nil {
			fmt.Printf("Failed to check comment: %v", err)
			os.Exit(1)
		}
	}

	// analyze source file
	sourceFile := file.NewSourceFile(*sourceFilePath, *structName)

	// generate builder
	gen := generator.NewGenerator(sourceFile, builderFile)
	gen.Generate()

	os.Exit(0)
}
