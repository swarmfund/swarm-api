package types

//go:generate jsonenums -tprefix=false -transform=snake -type=LogEventType
type LogEventType int32

const (
	LogEventTypeUserCreated LogEventType = 1 << iota
	LogEventTypeSuccessfulTFA
	LogEventTypeUnsuccessfulTFA
	LogEventTypeCreateTFABackend
	LogEventTypeUpdateWalletFactor
	LogEventTypeDeleteWalletFactor
	LogEventTypeLoginSuccessful
)
