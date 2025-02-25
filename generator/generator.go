package generator

import (
	"bytes"
	"fmt"
	"os"
	"regexp"
	"test1/constant"
	"test1/file"
	"test1/targetstruct"

	"github.com/dave/jennifer/jen"
)

type Generator struct {
	sourceFile  *file.SourceFile
	builderFile *file.BuilderFile
}

func NewGenerator(sourceFile *file.SourceFile, builderFile *file.BuilderFile) *Generator {
	return &Generator{sourceFile: sourceFile, builderFile: builderFile}
}

func (g *Generator) Generate() error {
	// analyze
	ts := g.sourceFile.Analyze()

	// generate
	f := g.generate(ts)

	// convert
	buf := new(bytes.Buffer)
	if err := f.Render(buf); err != nil {
		return err
	}
	code := g.convert(buf.String())

	// output
	err := g.output(code)
	if err != nil {
		return err
	}
	return nil
}

func (g *Generator) removeBeforeCommonLine(code string) string {
	// 正規表現を使用して// Code generatedから始まる行以降のすべてのコードを抽出
	regex := regexp.MustCompile(`(?m)^\s*\/\/ Code generated.*$[\s\S]*`)
	matches := regex.FindString(code)
	return matches
}

func (g *Generator) generate(ts targetstruct.TargetStruct) *jen.File {
	// ファイルオブジェクト作成
	f := jen.NewFile(ts.PackageName())

	// コメント追加
	f.Comment(constant.COMMENT)

	// Builder struct作成
	f.Type().Id(ts.GetBuilderName()).Struct(
		jen.Id("d").Id(ts.StructName()),
	)

	// setter作成
	for _, field := range ts.Fields() {
		f.Func().Params(
			jen.Id("b").Op("*").Id(ts.GetBuilderName()),
		).Id(field.GetSetterName()).Params(
			jen.Id(field.Name()).Id(field.TypeName()),
		).Op("*").Id(ts.GetBuilderName()).Block(
			jen.Id("b").Dot("d").Dot(field.Name()).Op("=").Id(field.Name()),
			jen.Return(jen.Id("b")),
		)
	}

	// readBuild作成
	f.Func().Params(
		jen.Id("b").Op("*").Id(ts.GetBuilderName()),
	).Id("ReadBuild").Params().Op("*").Id(ts.StructName()).Block(
		jen.Return(jen.Op("&").Id("b").Dot("d")),
	)

	// Clone method with source Item
	f.Func().Params(
		jen.Id("b").Op("*").Id(ts.GetBuilderName()),
	).Id("Clone").Params(
		jen.Id("src").Op("*").Id(ts.StructName()),
	).Op("*").Id(ts.GetBuilderName()).Block(g.GetMoveFieldStatement(ts)...)

	// コメント追加
	f.Comment(constant.COMMENT)

	return f
}

func (g *Generator) GetMoveFieldStatement(ts targetstruct.TargetStruct) []jen.Code {
	codes := []jen.Code{}
	for _, field := range ts.Fields() {
		codes = append(codes, jen.Id("b").Dot("d").Dot(field.Name()).Op("=").Id("src").Dot(field.Name()))
	}
	codes = append(codes, jen.Return(jen.Id("b")))
	return codes
}

func (g *Generator) convert(code string) string {
	// 新規作成時は何もしない
	if !g.builderFile.IsExistFile() {
		return code
	}

	// 既存ファイルの最初のコメントよりも前の部分を取得
	beforeLines := g.builderFile.GetBeforeCommentLines()

	// 生成コードの最初のコメント以降の部分を取得
	filteredCode := g.removeBeforeCommonLine(code)

	// 既存ファイルの２つ目のコメント以降を取得
	afterLines := g.builderFile.GetAfterCommentLines()

	return beforeLines + filteredCode + afterLines
}

func (g *Generator) output(code string) error {
	err := os.MkdirAll(g.builderFile.GetTargetDir(), 0755)
	if err != nil {
		fmt.Printf("Failed to create directory: %v\n", err)
		return err
	}

	file, err := os.Create(g.builderFile.GetTargetFilePath())
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(bytes.NewBufferString(code).Bytes())
	if err != nil {
		return err
	}
	return nil
}
