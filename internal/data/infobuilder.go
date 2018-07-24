package data

import (
	"gitlab.com/tokend/go/xdrbuild"
)

type Infobuilder func(info Info) *xdrbuild.Transaction
