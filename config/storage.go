package config

import "time"

type Storage struct {
	DisableStorage   bool
	AccessKey        string
	SecretKey        string
	Host             string
	ForceSSL         bool
	FormDataExpire   time.Duration
	MinContentLength int64
	MaxContentLength int64
	ObjectCreateARN  string
	Listener         struct {
		BrokerURL    string
		Exchange     string
		ExchangeType string
		BindingKey   string
	}
}

//
//func (s *Storage) DefineConfigStructure() {
//	for _, key := range []string{"disable", "access_key", "secret_key", "host", "force_ssl",
//		"form_data_expire", "min_content_length", "max_content_length", "object_create_arn",
//		"listener_broker_url", "listener_exchange", "listener_exchange_type", "listener_binding_key"} {
//		s.bindEnv(key)
//	}
//	s.setDefault("disable", false)
//	s.setDefault("force_ssl", false)
//	s.setDefault("form_data_expire", 3600*time.Second)
//	s.setDefault("min_content_length", 1)
//	s.setDefault("max_content_length", 1024*1024*50)
//	s.setDefault("listener_exchange_type", "direct")
//}
//
//func (s *Storage) Init() error {
//	var err error
//
//	s.DisableStorage = s.getBool("disable")
//	if s.DisableStorage {
//		return nil
//	}
//
//	s.AccessKey, err = s.getNonEmptyString("access_key")
//	if err != nil {
//		return err
//	}
//
//	s.SecretKey, err = s.getNonEmptyString("secret_key")
//	if err != nil {
//		return err
//	}
//
//	s.Host, err = s.getNonEmptyString("host")
//	if err != nil {
//		return err
//	}
//
//	s.ForceSSL = s.getBool("force_ssl")
//
//	formDataExpire := s.getInt("form_data_expire")
//	s.FormDataExpire = time.Duration(formDataExpire) * time.Second
//
//	s.MinContentLength = int64(s.getInt("min_content_length"))
//
//	s.MaxContentLength = int64(s.getInt("max_content_length"))
//
//	s.ObjectCreateARN, err = s.getNonEmptyString("object_create_arn")
//	if err != nil {
//		return err
//	}
//
//	s.Listener.BindingKey, err = s.getNonEmptyString("listener_binding_key")
//	if err != nil {
//		return err
//	}
//
//	s.Listener.BrokerURL, err = s.getNonEmptyString("listener_broker_url")
//	if err != nil {
//		return err
//	}
//
//	s.Listener.Exchange, err = s.getNonEmptyString("listener_exchange")
//	if err != nil {
//		return err
//	}
//
//	s.Listener.ExchangeType = s.getString("listener_exchange_type")
//	if err != nil {
//		return err
//	}
//
//	return nil
//}
