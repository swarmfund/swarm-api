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
	DocumentTypeKYCIdDocument
	DocumentTypeKYCPoa
)

func IsPublicDocument(t DocumentType) bool {
	return t != DocumentTypeKYCIdDocument && t != DocumentTypeKYCPoa
}
