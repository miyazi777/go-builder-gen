package model

type ItemBuilder struct {
	d Item
}

func NewItemBuilder() *ItemBuilder {
	return &ItemBuilder{}
}

func (b *ItemBuilder) SetID(id OptInt64) *ItemBuilder {
	b.d.id = id
	return b
}

func (b *ItemBuilder) SetCommunicationID(communicationID int64) *ItemBuilder {
	b.d.communicationID = communicationID
	return b
}

func (b *ItemBuilder) SetAssessmentOfferID(assessmentOfferID int64) *ItemBuilder {
	b.d.assessmentOfferID = assessmentOfferID
	return b
}

func (b *ItemBuilder) ReadBuild() *Item {
	return &b.d
}

func (b *ItemBuilder) Clone(src *Item) *Item {
	b.d.id = src.id
	b.d.communicationID = src.communicationID
	b.d.assessmentOfferID = src.assessmentOfferID
	return &b.d
}
