package sso

import (
	"encoding/base64"
	"net/http"

	"net/url"

	"fmt"

	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/swarmfund/api/internal/clienturl"
)

type (
	SSOHandlerRequest struct {
		Signature string
		Nonce     string
		ReturnURL string
		ClientURL string
	}
)

func NewSSOHandlerRequest(r *http.Request) (SSOHandlerRequest, error) {
	values := r.URL.Query()
	request := SSOHandlerRequest{
		Signature: values.Get("sig"),
		ClientURL: values.Get("client_url"),
	}
	// parse sso payload
	{
		payload, err := base64.StdEncoding.DecodeString(values.Get("sso"))
		if err != nil {
			return request, errors.Wrap(err, "failed to decode payload")
		}
		fmt.Println(string(payload))
		values, err := url.ParseQuery(string(payload))
		if err != nil {
			return request, errors.Wrap(err, "failed to parse payload")
		}
		request.Nonce = values.Get("nonce")
		request.ReturnURL = values.Get("return_sso_url")
	}
	return request, nil
}

func SSOReceiver(w http.ResponseWriter, r *http.Request) {
	request, err := NewSSOHandlerRequest(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	payload := clienturl.DiscourseSSO(request.Signature, request.ReturnURL, request.Nonce)

	encoded, err := payload.Encode()

	// TODO proper url build
	http.Redirect(w, r, fmt.Sprintf("%s/%s", request.ClientURL, encoded), http.StatusTemporaryRedirect)
}
