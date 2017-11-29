package config

//type Core struct {
//	*Base
//	AccountManagerKey string
//}
//
//func (c *Core) DefineConfigStructure() {
//	c.bindEnv("account_manager_key")
//}
//
//func (c *Core) Init() error {
//	var err error
//	c.AccountManagerKey, err = c.getNonEmptyString("account_manager_key")
//	if err != nil {
//		return err
//	}
//	kp, err := keypair.Parse(c.AccountManagerKey)
//	if err != nil {
//		return errors.New("Could not parse Account Manager Key")
//	}
//	_, ok := kp.(*keypair.Full)
//	if !ok {
//		return errors.New("Account Manager Key must be a Secret Key, not Public Key")
//	}
//
//	return nil
//}
