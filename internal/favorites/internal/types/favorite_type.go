package types

//go:generate jsonenums -tprefix=false -transform=snake -type=FavoriteType
type FavoriteType int32

const (
	FavoriteTypeSale FavoriteType = 1 << iota
	FavoriteTypeAssetPair
	FavoriteTypeSettings
)
