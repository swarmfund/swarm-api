package handlers

import (
	"net/http"

	"encoding/json"

	"net/url"

	"github.com/go-chi/chi"
	. "github.com/go-ozzo/ozzo-validation"
	"github.com/spf13/cast"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/api/internal/api/resources"
	"gitlab.com/swarmfund/api/internal/api/urlval"
	"gitlab.com/swarmfund/api/internal/types"
	"gitlab.com/tokend/go/doorman"
)

type (
	BlobIndexRequest struct {
		Filter        BlobIndexFilter
		Address       types.Address
		Relationships map[string]string
	}

	BlobIndexFilter struct {
		Page uint64  `url:"page"`
		Type *uint64 `url:"type"`
	}
)

func NewBlobIndexRequest(r *http.Request) (BlobIndexRequest, error) {
	values := r.URL.Query()
	filter, err := NewBlobIndexFilter(values)
	if err != nil {
		return BlobIndexRequest{}, err
	}

	request := BlobIndexRequest{
		Filter:        filter,
		Address:       types.Address(chi.URLParam(r, "address")),
		Relationships: make(map[string]string),
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

func NewBlobIndexFilter(values url.Values) (BlobIndexFilter, error) {
	filter := BlobIndexFilter{
		Page: 1,
	}

	if err := urlval.Decode(values, &filter); err != nil {
		return filter, errors.Wrap(err, "failed to populate")
	}

	return filter, filter.Validate()
}

func (r BlobIndexFilter) Validate() error {
	return ValidateStruct(&r,
		Field(&r.Page, Min(uint64(1))),
		Field(&r.Type, Min(uint64(1))),
	)
}

func BlobIndex(w http.ResponseWriter, r *http.Request) {
	request, err := NewBlobIndexRequest(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	q := BlobQ(r).
		ByOwner(request.Address).
		ByRelationships(request.Relationships)

	filter := request.Filter
	if filter.Type != nil {
		blobType := types.BlobType(cast.ToInt32(filter.Type))
		q = q.ByType(blobType)
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
