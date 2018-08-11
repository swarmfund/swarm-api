package config

import "net/url"

func (v *ViperConfig) KYCIndex() *url.URL {
	u, _ := url.Parse("http://kycprovider:7006")
	return u
}
