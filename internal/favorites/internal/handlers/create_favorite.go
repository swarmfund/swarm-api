package handlers

import (
	"net/http"

	"encoding/json"

	"github.com/go-chi/chi"
	. "github.com/go-ozzo/ozzo-validation"
	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/swarmfund/api/internal/api/handlers"
	"gitlab.com/swarmfund/api/internal/api/movetoape"
	"gitlab.com/swarmfund/api/internal/favorites/internal/data"
	"gitlab.com/swarmfund/api/internal/favorites/internal/resources"
	"gitlab.com/swarmfund/api/internal/types"
	"gitlab.com/swarmfund/go/doorman"
)

type CreateFavoriteRequest struct {
	Owner types.Address `json:"-"`
	resources.Favorite
}

func NewCreateFavoriteRequest(r *http.Request) (CreateFavoriteRequest, error) {
	request := CreateFavoriteRequest{
		Owner: types.Address(chi.URLParam(r, "address")),
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return request, errors.Wrap(err, "failed to unmarshal")
	}
	return request, request.Validate()
}

func (r CreateFavoriteRequest) Validate() error {
	return ValidateStruct(&r,
		Field(&r.Owner, Required),
		Field(&r.Favorite, Required),
	)
}

func CreateFavorite(w http.ResponseWriter, r *http.Request) {
	request, err := NewCreateFavoriteRequest(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	if err := handlers.Doorman(r, doorman.SignatureOf(string(request.Owner))); err != nil {
		movetoape.RenderDoormanErr(w, err)
		return
	}

	favorite := data.Favorite{
		Owner: request.Owner,
		Type:  request.Data.Type,
		Key:   request.Data.Attributes.Key,
	}

	if err := FavoritesQ(r).Create(favorite); err != nil {
		if err == data.ErrOwnerViolated {
			ape.RenderErr(w, problems.NotFound())
			return
		}
		handlers.Log(r).WithError(err).Error("failed to create favorite")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	w.WriteHeader(204)
}
