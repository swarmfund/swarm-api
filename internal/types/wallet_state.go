package types

//go:generate jsonenums -type=WalletState -tprefix=false -transform=snake
type WalletState int32

const (
	WalletStateNotVerified WalletState = 1 << iota
	WalletStateVerified
)
