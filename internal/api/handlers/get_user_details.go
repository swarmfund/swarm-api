package handlers

import (
	"encoding/json"
	"net/http"

	. "github.com/go-ozzo/ozzo-validation"
	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/swarmfund/api/db2/api"
	"gitlab.com/swarmfund/api/internal/types"
)

type (
	GetDetailsRequest struct {
		Addresses []string `json:"addresses"`
	}

	GetDetailsResponse struct {
		Users map[types.Address]api.User `json:"users"`
	}
)

func NewGetDetailsRequest(r *http.Request) (GetDetailsRequest, error) {
	var request GetDetailsRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return request, errors.Wrap(err, "failed to unmarshal")
	}

	return request, nil
}

func (r GetDetailsRequest) Validate() error {
	return ValidateStruct(&r,
		Field(&r.Addresses, Required),
	)
}

func GetUsersDetails(w http.ResponseWriter, r *http.Request) {
	request, err := NewGetDetailsRequest(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	users, err := UsersQ(r).ByAddresses(request.Addresses)

	if err != nil {
		Log(r).WithError(err).Error("Failed to get users")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	var response GetDetailsResponse
	for _, user := range users {
		response.Users[user.Address] = user
	}

	json.NewEncoder(w).Encode(response)

	return
}
