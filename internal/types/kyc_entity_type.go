package types

//go:generate jsonenums -tprefix=false -transform=snake -type=KYCEntityType
type KYCEntityType int32

const (
	KYCEntityTypeIndividual KYCEntityType = 1 << iota
	KYCEntityTypeSyndicate
)
