package config

type TFA struct {
	//*Base
	Dev bool
}

//
//func (t *TFA) DefineConfigStructure() {
//	t.bindEnv("dev")
//
//	t.setDefault("dev", false)
//}
//
//func (t *TFA) Init() {
//	t.Dev = t.getBool("dev")
//}
