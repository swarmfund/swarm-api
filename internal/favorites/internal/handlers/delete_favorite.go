package handlers

import (
	"net/http"

	"strconv"

	"github.com/go-chi/chi"
	. "github.com/go-ozzo/ozzo-validation"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/swarmfund/api/internal/api/handlers"
	"gitlab.com/swarmfund/api/internal/api/movetoape"
	"gitlab.com/swarmfund/api/internal/types"
	"gitlab.com/swarmfund/go/doorman"
)

type DeleteFavoriteRequest struct {
	Owner types.Address `json:"-"`
	ID    int64         `json:"-"`
}

func NewDeleteFavoriteRequest(r *http.Request) (DeleteFavoriteRequest, error) {
	request := DeleteFavoriteRequest{
		Owner: types.Address(chi.URLParam(r, "address")),
	}
	rawid := chi.URLParam(r, "favorite")
	id, err := strconv.ParseInt(rawid, 0, 10)
	if err != nil {
		return request, Errors{
			"id": err,
		}
	}
	request.ID = id
	return request, request.Validate()
}

func (r DeleteFavoriteRequest) Validate() error {
	return ValidateStruct(&r,
		Field(&r.ID, Required),
		Field(&r.Owner, Required),
	)
}

func DeleteFavorite(w http.ResponseWriter, r *http.Request) {
	request, err := NewDeleteFavoriteRequest(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	if err := handlers.Doorman(r, doorman.SignerOf(string(request.Owner))); err != nil {
		movetoape.RenderDoormanErr(w, err)
		return
	}

	if err := FavoritesQ(r).Delete(request.Owner, request.ID); err != nil {
		handlers.Log(r).WithError(err).Error("failed to delete favorite")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	w.WriteHeader(204)
}
