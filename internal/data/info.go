package data

import "gitlab.com/tokend/horizon-connector"

type Info interface {
	Info() (*horizon.Info, error)
}
