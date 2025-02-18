package generator

import (
	"go/ast"
	"go/parser"
	"go/token"
	"log"

	"golang.org/x/tools/go/ast/astutil"
)

func AnalyzeStruct(filePath string, targetStruct string) TargetStruct {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		log.Fatalf("Failed to parse file: %v", err)
	}

	packageName := ""
	targetFields := []TargetField{}
	astutil.Apply(node, nil, func(c *astutil.Cursor) bool {
		n := c.Node()

		// パッケージ名を取得
		packageName = node.Name.Name

		typeSpec, ok := n.(*ast.TypeSpec)
		if !ok {
			return true
		}

		// 指定した構造体以外はスキップ
		if typeSpec.Name.Name != targetStruct {
			return true
		}

		structType, ok := typeSpec.Type.(*ast.StructType)
		if !ok {
			return true
		}

		// 構造体のフィールド名と型名を取得
		for _, field := range structType.Fields.List {
			for _, name := range field.Names {
				typeStr := getFieldType(field.Type)
				targetFields = append(targetFields, TargetField{
					Name: name.Name,
					Type: typeStr,
				})
			}
		}
		return true
	})

	gs := TargetStruct{
		PackageName: packageName,
		StructName:  targetStruct,
		Fields:      targetFields,
	}
	return gs
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
