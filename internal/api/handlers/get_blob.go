package handlers

import (
	"net/http"

	"encoding/json"

	"github.com/go-chi/chi"
	. "github.com/go-ozzo/ozzo-validation"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/swarmfund/api/internal/api/resources"
	"gitlab.com/swarmfund/api/internal/types"
)

type (
	GetBlobRequest struct {
		Address types.Address `json:"-"`
		BlobID  string        `json:"-"`
	}
)

func NewGetBlobRequest(r *http.Request) (GetBlobRequest, error) {
	request := GetBlobRequest{
		Address: types.Address(chi.URLParam(r, "address")),
		BlobID:  chi.URLParam(r, "blob"),
	}
	return request, request.Validate()
}

func (r GetBlobRequest) Validate() error {
	return ValidateStruct(&r,
		Field(&r.Address, Required),
		Field(&r.BlobID, Required),
	)
}
func GetBlob(w http.ResponseWriter, r *http.Request) {
	request, err := NewGetBlobRequest(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	// TODO check allowed

	blob, err := BlobQ(r).Get(request.BlobID)
	if err != nil {
		Log(r).WithError(err).Error("failed to get blob")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	if blob == nil {
		ape.RenderErr(w, problems.NotFound())
		return
	}

	response := CreateBlobResponse{
		Data: resources.NewBlob(blob),
	}

	json.NewEncoder(w).Encode(&response)
}