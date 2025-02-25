package targetstruct

import (
	"fmt"
	"strings"
)

type TargetStruct struct {
	packageName string
	structName  string
	fields      []TargetField
}

func NewTaregtStruct(packageName, structName string, fields []TargetField) *TargetStruct {
	return &TargetStruct{
		packageName: packageName,
		structName:  structName,
		fields:      fields,
	}
}

func (t *TargetStruct) PackageName() string {
	return t.packageName
}

func (t *TargetStruct) StructName() string {
	return t.structName
}

func (t *TargetStruct) Fields() []TargetField {
	return t.fields
}

func (t *TargetStruct) GetBuilderName() string {
	return fmt.Sprintf("%sBuilder", t.structName)
}

type TargetField struct {
	name     string
	typeName string
}

func NewTargetField(name, typeName string) *TargetField {
	return &TargetField{
		name:     name,
		typeName: typeName,
	}
}

func (f *TargetField) Name() string {
	return f.name
}

func (f *TargetField) TypeName() string {
	return f.typeName
}

func (f *TargetField) GetSetterName() string {
	if f.name == "id" {
		return "SetID"
	}
	return "Set" + strings.ToUpper(f.name[:1]) + f.name[1:]
}
