package handlers

import (
	"net/http"

	"encoding/json"

	. "github.com/go-ozzo/ozzo-validation"
	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/swarmfund/api/internal/api/movetoape"
	"gitlab.com/swarmfund/api/internal/api/resources"
	"gitlab.com/swarmfund/api/internal/api/urlval"
	"gitlab.com/swarmfund/go/doorman"
)

type (
	WalletsIndexResponse struct {
		Data  WalletsIndexData   `json:"data"`
		Links urlval.FilterLinks `json:"links"`
	}
	WalletsIndexData    []resources.WalletData
	WalletsIndexFitlers struct {
		Page  uint64  `url:"page"`
		State *uint64 `url:"state"`
	}
)

func NewWalletsFilters(r *http.Request) (WalletsIndexFitlers, error) {
	filters := WalletsIndexFitlers{
		Page: 1,
	}
	if err := urlval.Decode(r.URL.Query(), &filters); err != nil {
		return filters, errors.Wrap(err, "failed to populate")
	}
	return filters, filters.Validate()
}

func (r WalletsIndexFitlers) Validate() error {
	return ValidateStruct(&r,
		Field(&r.Page, Min(uint64(1))),
		Field(&r.State, Min(uint64(1))),
	)
}

func WalletsIndex(w http.ResponseWriter, r *http.Request) {
	filters, err := NewWalletsFilters(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	if err := Doorman(r, doorman.SignerOf(CoreInfo(r).GetMasterAccountID())); err != nil {
		movetoape.RenderDoormanErr(w, err)
		return
	}

	q := WalletQ(r).Page(filters.Page)

	if filters.State != nil {
		q = q.ByState(*filters.State)
	}

	wallets, err := q.Select()
	if err != nil {
		Log(r).WithError(err).Error("failed to get wallets")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	response := WalletsIndexResponse{
		Data: make(WalletsIndexData, 0, len(wallets)),
	}
	for _, wallet := range wallets {
		response.Data = append(response.Data, resources.NewWallet(&wallet).Data)
	}
	response.Links = urlval.Encode(r, filters)

	json.NewEncoder(w).Encode(&response)
}
