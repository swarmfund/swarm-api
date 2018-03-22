package handlers

import (
	"net/http"
	"strconv"

	"encoding/json"

	"github.com/go-chi/chi"
	. "github.com/go-ozzo/ozzo-validation"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/api/internal/api/resources"
	"gitlab.com/swarmfund/api/internal/api/urlval"
	"gitlab.com/swarmfund/api/internal/types"
	"gitlab.com/swarmfund/go/doorman"
)

type (
	BlobIndexRequest struct {
		Address       types.Address
		Relationships map[string]string
		Type          *types.BlobType
	}

	BlobIndexFilter struct {
		Page          uint64            `url:"page"`
		Address       types.Address     `url:"address"`
		Type          *uint64           `url:"type"`
		Relationships map[string]string `url:"relationships"`
	}
)

func NewBlobIndexFilter(r *http.Request) (BlobIndexFilter, error) {
	filters := BlobIndexFilter{
		Page:    1,
		Address: types.Address(chi.URLParam(r, "address")),
	}

	if err := urlval.Decode(r.URL.Query(), &filters); err != nil {
		return filters, errors.Wrap(err, "failed to populate")
	}
	return filters, filters.Validate()
}

func (r BlobIndexFilter) Validate() error {
	return ValidateStruct(&r,
		Field(&r.Page, Min(uint64(1))),
		Field(&r.Type, Min(uint64(1))),
	)
}

func NewBlobIndexRequest(r *http.Request) (BlobIndexRequest, error) {
	request := BlobIndexRequest{
		Address:       types.Address(chi.URLParam(r, "address")),
		Relationships: make(map[string]string),
	}
	values := r.URL.Query()
	{
		if values.Get("type") != "" {
			raw, err := strconv.ParseInt(values.Get("type"), 0, 32)
			if err != nil {
				return request, err
			}
			tpe := types.BlobType(raw)
			request.Type = &tpe
		}
		values.Del("type")
	}
	for k := range values {
		request.Relationships[k] = values.Get(k)
	}
	return request, request.Validate()
}

func (r BlobIndexRequest) Validate() error {
	// TODO implement
	return nil
}

func BlobIndex(w http.ResponseWriter, r *http.Request) {
	request, err := NewBlobIndexRequest(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	filters := map[string]string{}
	for k := range r.URL.Query() {
		filters[k] = r.URL.Query().Get(k)
	}

	q := BlobQ(r).
		ByOwner(request.Address).
		ByRelationships(request.Relationships)

	if request.Type != nil {
		q = q.ByType(*request.Type)
	}

	records, err := q.Select()
	if err != nil {
		Log(r).WithError(err).Error("failed to get blobs")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	var isAllowed bool
	err = Doorman(r, doorman.SignerOf(string(request.Address)), doorman.SignerOf(CoreInfo(r).GetMasterAccountID()))
	if err == nil {
		isAllowed = true
	}

	var response struct {
		Data []resources.Blob `json:"data"`
	}

	response.Data = make([]resources.Blob, 0)

	for _, record := range records {
		// sorry
		if types.IsPublicBlob(record.Type) {
			response.Data = append(response.Data, resources.NewBlob(&record))
		} else if isAllowed {
			response.Data = append(response.Data, resources.NewBlob(&record))
		}
	}

	json.NewEncoder(w).Encode(response)

}
