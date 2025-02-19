package model

//go:generate go run ../main.go -struct Item -path=./$GOFILE
type Item struct {
	id                OptInt64 // ID
	communicationID   int64    // コミュニケーションID
	assessmentOfferID int64    // 査定依頼ID
}

type OptInt64 struct {
	value int64
}
