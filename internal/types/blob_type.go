package types

//go:generate jsonenums -tprefix=false -transform=snake -type=BlobType
type BlobType int32

const (
	BlobTypeAssetDescription BlobType = 1 << iota
	BlobTypeFundOverview
	BlobTypeFundUpdate
	BlobTypeNavUpdate
	BlobTypeFundDocument
	BlobTypeAlpha
	BlobTypeBravo
	BlobTypeCharlie
	BlobTypeDelta
	BlobTypeTokenTerms
	BlobTypeTokenMetrics
	BlobTypeKYCForm
	BlobTypeKYCIdDocument
	BlobTypeKYCPoa
)

func IsPublicBlob(t BlobType) bool {
	return t <= BlobTypeTokenMetrics
}
