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
	"gitlab.com/swarmfund/api/internal/api/urlval"
	"gitlab.com/swarmfund/api/internal/favorites/internal/resources"
	"gitlab.com/swarmfund/api/internal/favorites/internal/types"
	types2 "gitlab.com/swarmfund/api/internal/types"
	"gitlab.com/swarmfund/go/doorman"
)

type (
	FavoritesIndexResponse struct {
		Data  FavoritesIndexResponseData `json:"data"`
		Links urlval.FilterLinks         `json:"links"`
	}
	FavoritesIndexResponseData []resources.Favorite
	FavoriteFilters            struct {
		Page uint64  `url:"page"`
		Type *uint64 `url:"type"`
	}
)

func NewFavoriteFilters(r *http.Request) (FavoriteFilters, error) {
	filters := FavoriteFilters{
		Page: 1,
	}
	if err := urlval.Decode(r.URL.Query(), &filters); err != nil {
		return filters, errors.Wrap(err, "failed to populate")
	}
	return filters, filters.Validate()
}

func (r FavoriteFilters) Validate() error {
	return ValidateStruct(&r,
		Field(&r.Page, Min(uint64(1))),
		Field(&r.Type, Min(uint64(1))),
	)
}

func FavoriteIndex(w http.ResponseWriter, r *http.Request) {
	// TODO validate owner address
	owner := types2.Address(chi.URLParam(r, "address"))

	filters, err := NewFavoriteFilters(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	if err := handlers.Doorman(r, doorman.SignerOf(string(owner))); err != nil {
		movetoape.RenderDoormanErr(w, err)
		return
	}

	q := FavoritesQ(r).Page(filters.Page).ByOwner(owner)

	if filters.Type != nil {
		q = q.ByType(types.FavoriteType(*filters.Type))
	}

	records, err := q.Select()
	if err != nil {
		handlers.Log(r).WithError(err).Error("failed to get favorites")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	response := FavoritesIndexResponse{
		Data: make(FavoritesIndexResponseData, 0, len(records)),
	}
	for _, record := range records {
		response.Data = append(response.Data, resources.NewFavorite(record))
	}

	response.Links = urlval.Encode(r, filters)
	json.NewEncoder(w).Encode(&response)
}
