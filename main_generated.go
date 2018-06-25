package api

import (
	"github.com/zenazn/goji/web"
	"net/http"
)

// ServeHTTPC is a method for web.Handler
func (action NotFoundAction) ServeHTTPC(c web.C, w http.ResponseWriter, r *http.Request) {
	ap := &action.Action
	ap.Prepare(c, w, r)
	action.Log = action.Log.WithField("action", "NotFoundAction")
	ap.Execute(&action)
}

// ServeHTTPC is a method for web.Handler
func (action ParticipantsAction) ServeHTTPC(c web.C, w http.ResponseWriter, r *http.Request) {
	ap := &action.Action
	ap.Prepare(c, w, r)
	action.Log = action.Log.WithField("action", "ParticipantsAction")
	ap.Execute(&action)
}
