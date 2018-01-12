package types

import "fmt"

//go:generate jsonenums -tprefix=false -transform=snake -type=DocumentType
type DocumentType int32

const (
	DocumentTypeAssetLogo DocumentType = 1 << iota
	DocumentTypeFundLogo
	DocumentTypeFundDocument
	DocumentTypeNavReport
	DocumentTypeAlpha
	DocumentTypeBravo
	DocumentTypeCharlie
	DocumentTypeDelta
)

func IsPublicDocument(t DocumentType) bool {
	switch t {
	case DocumentTypeAssetLogo, DocumentTypeFundLogo, DocumentTypeFundDocument, DocumentTypeNavReport:
		return true
	default:
		panic(fmt.Errorf("unknown document type %s", t))
	}
}
