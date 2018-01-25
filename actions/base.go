package actions

import (
	"encoding/json"
	"net/http"

	gctx "github.com/goji/context"
	"github.com/zenazn/goji/web"
	"gitlab.com/swarmfund/api/render"
	"gitlab.com/swarmfund/api/render/problem"
	"golang.org/x/net/context"
)

// Base is a helper struct you can use as part of a custom action via
// composition.
//
// TODO: example usage
type Base struct {
	Ctx     context.Context
	GojiCtx web.C
	W       http.ResponseWriter
	R       *http.Request
	Err     error

	Signer string

	isSetup  bool
	IsSigned bool

	isJson   bool
	jsonBody map[string]string
}

// Prepare established the common attributes that get used in nearly every
// action.  "Child" actions may override this method to extend action, but it
// is advised you also call this implementation to maintain behavior.
func (base *Base) Prepare(c web.C, w http.ResponseWriter, r *http.Request) {
	base.Ctx = gctx.FromC(c)
	base.GojiCtx = c
	base.W = w
	base.R = r
}

func (base *Base) JsonValue(name string) string {
	if !base.isJson {
		return ""
	}

	if base.jsonBody == nil {
		decoder := json.NewDecoder(base.R.Body)
		err := decoder.Decode(&base.jsonBody)
		if err != nil {
			base.jsonBody = map[string]string{}
		}
	}

	value, ok := base.jsonBody[name]
	if !ok {
		return ""
	}
	return value
}

// Execute trigger content negottion and the actual execution of one of the
// action's handlers.
func (base *Base) Execute(action interface{}) {
	contentType := render.Negotiate(base.Ctx, base.R)

	switch contentType {
	case render.MimeHal, render.MimeJSON:
		action, ok := action.(JSON)

		if !ok {
			goto NotAcceptable
		}

		action.JSON()

		if base.Err != nil {
			problem.Render(base.Ctx, base.W, base.Err)
			return
		}
	case render.MimeRaw:
		action, ok := action.(Raw)

		if !ok {
			goto NotAcceptable
		}

		action.Raw()

		if base.Err != nil {
			problem.Render(base.Ctx, base.W, base.Err)
			return
		}
	default:
		goto NotAcceptable
	}
	return

NotAcceptable:
	problem.Render(base.Ctx, base.W, problem.NotAcceptable)
	return
}

// Do executes the provided func iff there is no current error for the action.
// Provides a nicer way to invoke a set of steps that each may set `action.Err`
// during execution
func (base *Base) Do(fns ...func()) {
	for _, fn := range fns {
		if base.Err != nil {
			return
		}

		fn()
	}
}
