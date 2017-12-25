package types

//go:generate jsonenums -tprefix=false -transform=snake -type=BlobType
type BlobType int32

const (
	BlobTypeAssetDescription BlobType = 1 << iota
)
