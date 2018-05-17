package config

import (
	"bytes"
	"sync"
	"testing"

	"github.com/spf13/viper"
)

func ConfigHelper(t *testing.T, raw string) ViperConfig {
	t.Helper()

	r := bytes.NewReader([]byte(raw))
	v := viper.New()
	v.SetConfigType("yaml")

	err := v.ReadConfig(r)
	if err != nil {
		t.Fatal(err)
	}

	return ViperConfig{
		Viper:   v,
		RWMutex: &sync.RWMutex{},
	}
}
