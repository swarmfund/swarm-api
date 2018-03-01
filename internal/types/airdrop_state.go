package types

//go:generate jsonenums -type=AirdropState -tprefix=false -transform=snake
type AirdropState int32

const (
	AirdropStateNotEligible AirdropState = 1 << iota
	AirdropStateEligible
	AirdropStateClaimed
)
