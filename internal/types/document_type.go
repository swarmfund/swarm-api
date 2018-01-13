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
	// TODO revert once past demo
	return true
}
