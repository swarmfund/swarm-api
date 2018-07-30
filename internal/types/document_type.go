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
	DocumentTypeKYCSelfie
	DocumentTypeGeneral //never use it, it needs to init default value of allowed media types
	DocumentTypeGeneralPublic
	DocumentTypeGeneralPrivate
	DocumentTypeAssetPhoto
)

func IsPublicDocument(t DocumentType) bool {
	return t != DocumentTypeKYCIdDocument && t != DocumentTypeKYCPoa && t != DocumentTypeKYCSelfie && t != DocumentTypeGeneralPrivate
}
