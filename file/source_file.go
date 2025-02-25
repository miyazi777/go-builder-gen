package file

import (
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	ts "test1/targetstruct"

	"golang.org/x/tools/go/ast/astutil"
)

type SourceFile struct {
	path       string
	structName string
}

func NewSourceFile(path, structName string) *SourceFile {
	return &SourceFile{path: path, structName: structName}
}

func (s *SourceFile) Analyze() ts.TargetStruct {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, s.path, nil, parser.ParseComments)
	if err != nil {
		log.Fatalf("Failed to parse file: %v", err)
	}

	packageName := ""
	targetFields := []ts.TargetField{}
	astutil.Apply(node, nil, func(c *astutil.Cursor) bool {
		n := c.Node()

		// パッケージ名を取得
		packageName = node.Name.Name

		typeSpec, ok := n.(*ast.TypeSpec)
		if !ok {
			return true
		}

		// 指定した構造体以外はスキップ
		if typeSpec.Name.Name != s.structName {
			return true
		}

		structType, ok := typeSpec.Type.(*ast.StructType)
		if !ok {
			return true
		}

		// 構造体のフィールド名と型名を取得
		for _, field := range structType.Fields.List {
			for _, name := range field.Names {
				typeName := getFieldType(field.Type)
				targetFields = append(targetFields, *ts.NewTargetField(name.Name, typeName))
			}
		}
		return true
	})

	return *ts.NewTaregtStruct(packageName, s.structName, targetFields)
}

func getFieldType(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.StarExpr:
		return "*" + getFieldType(t.X)
	case *ast.ArrayType:
		return "[]" + getFieldType(t.Elt)
	case *ast.MapType:
		return "map[" + getFieldType(t.Key) + "]" + getFieldType(t.Value)
	case *ast.SelectorExpr:
		return getFieldType(t.X) + "." + t.Sel.Name
	default:
		return "unknown"
	}
}
