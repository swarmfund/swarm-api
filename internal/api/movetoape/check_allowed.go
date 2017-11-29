package movetoape

import (
	"net/http"

	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/swarmfund/go/signcontrol"
)

func RenderDoormanErr(w http.ResponseWriter, err error) {
	switch err {
	case signcontrol.ErrNotAllowed, signcontrol.ErrNotSigned:
		ape.RenderErr(w, problems.NotAllowed())
	case nil:
		panic("expected not nil error")
	default:
		ape.RenderErr(w, problems.InternalError())
	}
}
