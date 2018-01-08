package handlers

import (
	"net/http"
	"strconv"

	"encoding/json"

	"github.com/go-chi/chi"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/swarmfund/api/internal/api/resources"
	"gitlab.com/swarmfund/api/internal/types"
)

type (
	BlobIndexRequest struct {
		Address types.Address
		Filters map[string]string
		Type    *types.BlobType
	}
)

func NewBlobIndexRequest(r *http.Request) (BlobIndexRequest, error) {
	request := BlobIndexRequest{
		Address: types.Address(chi.URLParam(r, "address")),
		Filters: make(map[string]string),
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
		request.Filters[k] = values.Get(k)
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
		ByRelationships(request.Filters)

	if request.Type != nil {
		q = q.ByType(*request.Type)
	}

	records, err := q.Select()
	if err != nil {
		Log(r).WithError(err).Error("failed to get blobs")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	var response struct {
		Data []resources.Blob `json:"data"`
	}

	response.Data = make([]resources.Blob, 0, len(records))

	for _, record := range records {
		response.Data = append(response.Data, resources.NewBlob(&record))
	}

	json.NewEncoder(w).Encode(response)

}
