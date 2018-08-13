package config

import "net/url"

func (v *ViperConfig) KYCIndex() *url.URL {
	u, _ := url.Parse("http://localhost:7006")
	return u
}
