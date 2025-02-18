package generator

import (
	"fmt"
	"strings"

	"github.com/dave/jennifer/jen"
)

type GenStruct struct {
	PackageName string
	StructName  string
	GenFields   []Field
}

type Field struct {
	Name string
	Type string
}

func (f *Field) GetSetterName() string {
	return "Set" + strings.ToUpper(f.Name[:1]) + f.Name[1:]
}

func NewGenStruct() *GenStruct {
	return &GenStruct{
		PackageName: "model",
		StructName:  "Item",
		GenFields: []Field{
			{Name: "id", Type: "OptInt64"},
			{Name: "communicationID", Type: "int64"},
			{Name: "assessmentOfferID", Type: "int64"},
		},
	}
}

func (g *GenStruct) GetBuilderName() string {
	return fmt.Sprintf("%sBuilder", g.StructName)
}

func (g *GenStruct) GetNewBuilderName() string {
	return fmt.Sprintf("New%s", g.GetBuilderName())
}

func (g *GenStruct) GetMoveFieldStatement() []jen.Code {
	codes := []jen.Code{}
	for _, field := range g.GenFields {
		codes = append(codes, jen.Id("b").Dot("d").Dot(field.Name).Op("=").Id("src").Dot(field.Name))
	}
	codes = append(codes, jen.Return(jen.Op("&").Id("b").Dot("d")))
	return codes
}

func (g *GenStruct) Generate() *jen.File {
	// ファイルオブジェクト作成
	f := jen.NewFile(g.PackageName)

	// Builder struct作成
	f.Type().Id(g.GetBuilderName()).Struct(
		jen.Id("d").Id(g.StructName),
	)

	// NewBuilder function作成
	f.Func().Id(g.GetNewBuilderName()).Params().Op("*").Id(g.GetBuilderName()).Block(
		jen.Return(jen.Op("&").Id(g.GetBuilderName()).Block()),
	)

	// setter作成
	for _, field := range g.GenFields {
		f.Func().Params(
			jen.Id("b").Op("*").Id(g.GetBuilderName()),
		).Id(field.GetSetterName()).Params(
			jen.Id(field.Name).Id(field.Type),
		).Op("*").Id(g.GetBuilderName()).Block(
			jen.Id("b").Dot("d").Dot(field.Name).Op("=").Id(field.Name),
			jen.Return(jen.Id("b")),
		)
	}

	// readbuild作成
	f.Func().Params(
		jen.Id("b").Op("*").Id(g.GetBuilderName()),
	).Id("ReadBuild").Params().Op("*").Id(g.StructName).Block(
		jen.Return(jen.Op("&").Id("b").Dot("d")),
	)

	// Clone method with source Item
	f.Func().Params(
		jen.Id("b").Op("*").Id(g.GetBuilderName()),
	).Id("Clone").Params(
		jen.Id("src").Op("*").Id(g.StructName),
	).Op("*").Id(g.StructName).Block(g.GetMoveFieldStatement()...)

	return f
}
