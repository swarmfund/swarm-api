package handlers

import (
	"net/http"

	"encoding/json"

	"github.com/go-chi/chi"
	. "github.com/go-ozzo/ozzo-validation"
	"github.com/spf13/cast"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/api/internal/api/resources"
	"gitlab.com/swarmfund/api/internal/api/urlval"
	"gitlab.com/swarmfund/api/internal/types"
	"gitlab.com/swarmfund/go/doorman"
)

type (
	BlobIndexFilter struct {
		Page          uint64            `url:"page"`
		Address       types.Address     `url:"address"`
		Type          *types.BlobType   `url:"type"`
		Relationships map[string]string `url:"relationships"`
	}
)

func NewBlobIndexFilter(r *http.Request) (BlobIndexFilter, error) {
	filter := BlobIndexFilter{
		Page:          1,
		Address:       types.Address(chi.URLParam(r, "address")),
		Relationships: make(map[string]string),
	}

	values := r.URL.Query()
	if err := urlval.DecodeWithValues(&values, &filter); err != nil {
		return filter, errors.Wrap(err, "failed to populate")
	}

	if values.Get("type") != "" {
		raw, err := cast.ToInt32E(values.Get("type")) //strconv.ParseInt(values.Get("type"), 0, 32)
		if err != nil {
			return filter, err
		}
		tpe := types.BlobType(raw)
		filter.Type = &tpe

		values.Del("type")
	}

	for k := range values {
		filter.Relationships[k] = values.Get(k)
	}

	return filter, filter.Validate()
}

func (r BlobIndexFilter) Validate() error {
	return ValidateStruct(&r,
		Field(&r.Address, Required),
		Field(&r.Page, Min(uint64(1))),
		Field(&r.Type, Min(int32(1))),
	)
}

func BlobIndex(w http.ResponseWriter, r *http.Request) {
	request, err := NewBlobIndexFilter(r)
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
