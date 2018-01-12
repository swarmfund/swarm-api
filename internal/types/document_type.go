package types

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
	DocumentTypeTokenTerms
	DocumentTypeTokenMetrics
)

func IsPublicDocument(t DocumentType) bool {
	switch t {
	case DocumentTypeAssetLogo, DocumentTypeFundLogo, DocumentTypeFundDocument, DocumentTypeNavReport, DocumentTypeTokenTerms, DocumentTypeTokenMetrics:
		return true
	default:
		return false
	}
}
