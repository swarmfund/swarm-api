package api

const (
	TFAActionLogin                   = "login"
	TFAActionUpdateGoogleTOTPBackend = "update_google_totp_backend"
)

type TFA struct {
	ID int64 `db:"id"`
	// TODO
	BackendID int64 `db:"backend"`
	// TODO remove
	OTPData  []byte `db:"otp_data"`
	Token    string `db:"token"`
	Verified bool   `db:"verified"`
}
