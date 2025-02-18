package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/dave/jennifer/jen"
)

type Generator struct {
	TargetFilePath string
	TargetStruct   TargetStruct
}

func NewGenerator(filePath string, targetStruct string) *Generator {
	// 構造体を解析する
	ts := AnalyzeStruct(filePath, targetStruct)

	return &Generator{
		TargetFilePath: filePath,
		TargetStruct:   ts,
	}
}

func (g *Generator) Generate() error {
	f := g.generateBuilder()

	outputFile := filepath.Join(g.getTargetDir(), g.getBuilderFileName())

	err := os.MkdirAll(g.getTargetDir(), 0755)
	if err != nil {
		fmt.Printf("Failed to create directory: %v\n", err)
		return err
	}

	if err := f.Save(outputFile); err != nil {
		fmt.Printf("Failed to save file: %v\n", err)
		return err
	}
	return nil
}

func (g *Generator) generateBuilder() *jen.File {
	ts := g.TargetStruct
	// ファイルオブジェクト作成
	f := jen.NewFile(ts.PackageName)

	// Builder struct作成
	f.Type().Id(ts.GetBuilderName()).Struct(
		jen.Id("d").Id(ts.StructName),
	)

	// NewBuilder function作成
	f.Func().Id(ts.GetNewBuilderName()).Params().Op("*").Id(ts.GetBuilderName()).Block(
		jen.Return(jen.Op("&").Id(ts.GetBuilderName()).Block()),
	)

	// setter作成
	for _, field := range ts.Fields {
		f.Func().Params(
			jen.Id("b").Op("*").Id(ts.GetBuilderName()),
		).Id(field.GetSetterName()).Params(
			jen.Id(field.Name).Id(field.Type),
		).Op("*").Id(ts.GetBuilderName()).Block(
			jen.Id("b").Dot("d").Dot(field.Name).Op("=").Id(field.Name),
			jen.Return(jen.Id("b")),
		)
	}

	// readbuild作成
	f.Func().Params(
		jen.Id("b").Op("*").Id(ts.GetBuilderName()),
	).Id("ReadBuild").Params().Op("*").Id(ts.StructName).Block(
		jen.Return(jen.Op("&").Id("b").Dot("d")),
	)

	// Clone method with source Item
	f.Func().Params(
		jen.Id("b").Op("*").Id(ts.GetBuilderName()),
	).Id("Clone").Params(
		jen.Id("src").Op("*").Id(ts.StructName),
	).Op("*").Id(ts.StructName).Block(ts.GetMoveFieldStatement()...)

	return f
}

func (g *Generator) getTargetDir() string {
	return filepath.Dir(g.TargetFilePath)
}

func (g *Generator) getBuilderFileName() string {
	baseNameWithExt := filepath.Base(g.TargetFilePath)
	baseName := strings.TrimSuffix(baseNameWithExt, filepath.Ext(baseNameWithExt))
	return fmt.Sprintf("%s_builder.go", baseName)
}

type TargetStruct struct {
	PackageName string
	StructName  string
	Fields      []TargetField
}

func (t *TargetStruct) GetBuilderName() string {
	return fmt.Sprintf("%sBuilder", t.StructName)
}

func (t *TargetStruct) GetNewBuilderName() string {
	return fmt.Sprintf("New%s", t.GetBuilderName())
}

func (t *TargetStruct) GetMoveFieldStatement() []jen.Code {
	codes := []jen.Code{}
	for _, field := range t.Fields {
		codes = append(codes, jen.Id("b").Dot("d").Dot(field.Name).Op("=").Id("src").Dot(field.Name))
	}
	codes = append(codes, jen.Return(jen.Op("&").Id("b").Dot("d")))
	return codes
}

type TargetField struct {
	Name string
	Type string
}

func (f *TargetField) GetSetterName() string {
	return "Set" + strings.ToUpper(f.Name[:1]) + f.Name[1:]
}
