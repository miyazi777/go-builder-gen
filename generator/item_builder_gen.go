package generator

import (
	"github.com/dave/jennifer/jen"
)

func GenerateItemBuilder() *jen.File {
	f := jen.NewFile("model")

	// ItemBuilder struct
	f.Type().Id("ItemBuilder").Struct(
		jen.Id("d").Id("Item"),
	)

	// NewItemBuilder function
	f.Func().Id("NewItemBuilder").Params().Op("*").Id("ItemBuilder").Block(
		jen.Return(jen.Op("&").Id("ItemBuilder").Block()),
	)

	// SetID method
	f.Func().Params(
		jen.Id("b").Op("*").Id("ItemBuilder"),
	).Id("SetID").Params(
		jen.Id("id").Id("OptInt64"),
	).Op("*").Id("ItemBuilder").Block(
		jen.Id("b").Dot("d").Dot("id").Op("=").Id("id"),
		jen.Return(jen.Id("b")),
	)

	// SetCommunicationID method
	f.Func().Params(
		jen.Id("b").Op("*").Id("ItemBuilder"),
	).Id("SetCommunicationID").Params(
		jen.Id("communicationID").Int64(),
	).Op("*").Id("ItemBuilder").Block(
		jen.Id("b").Dot("d").Dot("communicationID").Op("=").Id("communicationID"),
		jen.Return(jen.Id("b")),
	)

	// SetAssessmentOfferID method
	f.Func().Params(
		jen.Id("b").Op("*").Id("ItemBuilder"),
	).Id("SetAssessmentOfferID").Params(
		jen.Id("assessmentOfferID").Int64(),
	).Op("*").Id("ItemBuilder").Block(
		jen.Id("b").Dot("d").Dot("assessmentOfferID").Op("=").Id("assessmentOfferID"),
		jen.Return(jen.Id("b")),
	)

	// ReadBuild method
	f.Func().Params(
		jen.Id("b").Op("*").Id("ItemBuilder"),
	).Id("ReadBuild").Params().Op("*").Id("Item").Block(
		jen.Return(jen.Op("&").Id("b").Dot("d")),
	)

	// Clone method with source Item
	f.Func().Params(
		jen.Id("b").Op("*").Id("ItemBuilder"),
	).Id("Clone").Params(
		jen.Id("src").Op("*").Id("Item"),
	).Op("*").Id("Item").Block(
		jen.Id("b").Dot("d").Dot("id").Op("=").Id("src").Dot("id"),
		jen.Id("b").Dot("d").Dot("communicationID").Op("=").Id("src").Dot("communicationID"),
		jen.Id("b").Dot("d").Dot("assessmentOfferID").Op("=").Id("src").Dot("assessmentOfferID"),
		jen.Return(jen.Op("&").Id("b").Dot("d")),
	)

	return f
}
