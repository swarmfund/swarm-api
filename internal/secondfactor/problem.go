package secondfactor

// FactorRequiredErr implements error interface and provides additional info about factor verification token
type FactorRequiredErr struct {
	token string
	meta  map[string]interface{}
}

func (e FactorRequiredErr) Error() string {
	return "second factor required"
}

func (e FactorRequiredErr) Token() string {
	return e.token
}

func (e FactorRequiredErr) Meta() *map[string]interface{} {
	return &e.meta
}
