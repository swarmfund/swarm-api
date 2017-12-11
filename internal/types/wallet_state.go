package types

type WalletState int32

const (
	WalletStateNotVerified WalletState = 1 << iota
	WalletStateVerified
)
