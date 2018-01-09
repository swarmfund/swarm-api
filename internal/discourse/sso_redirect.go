package discourse

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"

	"encoding/json"

	"github.com/go-ozzo/ozzo-validation"
	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/swarmfund/api/internal/api/handlers"
)

type (
	SSORedirectRequest struct {
		Address   string
		Nonce     string
		Signature string
		ReturnURL string
	}
	SSORedirectResponse struct {
		Location string
	}
)

func NewSSORedirectRequest(r *http.Request) (SSORedirectRequest, error) {
	values := r.URL.Query()
	request := SSORedirectRequest{
		Address:   values.Get("address"),
		Nonce:     values.Get("nonce"),
		ReturnURL: values.Get("return"),
		Signature: values.Get("signature"),
	}

	return request, nil
}

func (r SSORedirectRequest) Validate() error {
	// TODO implement
	return nil
}

func SSORedirect(w http.ResponseWriter, r *http.Request) {
	request, err := NewSSORedirectRequest(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	request = request

	// TODO validate signature

	// TODO check allowed

	user, err := handlers.UsersQ(r).ByAddress(request.Address)
	if err != nil {
		handlers.Log(r).WithError(err).Error("failed to get user")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	if user == nil {
		ape.RenderErr(w, problems.BadRequest(validation.Errors{
			"address": errors.New("user does not exist"),
		})...)
		return
	}

	payload := url.Values{}
	payload.Set("nonce", request.Nonce)
	payload.Set("email", user.Email)
	payload.Set("external_id", string(user.Address))

	encodedPayload := payload.Encode()
	b64 := base64.StdEncoding.EncodeToString([]byte(encodedPayload))
	escaped := url.QueryEscape(b64)

	hash := hmac.New(sha256.New, []byte("super-sekrit"))
	hash.Write([]byte(b64))
	sig := fmt.Sprintf("%x", hash.Sum(nil))

	// TODO build url properly
	response := SSORedirectResponse{
		Location: fmt.Sprintf("%s?sso=%s&sig=%s", request.ReturnURL, escaped, sig),
	}
	json.NewEncoder(w).Encode(&response)
}
