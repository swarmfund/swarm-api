package handlers

import (
	"net/http"

	"encoding/json"

	"fmt"

	"encoding/base32"

	"github.com/go-chi/chi"
	. "github.com/go-ozzo/ozzo-validation"
	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/swarmfund/api/internal/api/resources"
	"gitlab.com/swarmfund/api/internal/types"
	"gitlab.com/swarmfund/go/hash"
)

type (
	CreateBlobRequest struct {
		Address types.Address         `json:"-"`
		Data    CreateBlobRequestData `json:"data"`
	}
	CreateBlobRequestData struct {
		Type       types.BlobType              `json:"type"`
		Attributes CreateBlobRequestAttributes `json:"attributes"`
	}
	CreateBlobRequestAttributes struct {
		Value string `json:"value"`
	}
	CreateBlobResponse struct {
		Data resources.Blob `json:"data"`
	}
)

func NewCreateBlobRequest(r *http.Request) (CreateBlobRequest, error) {
	request := CreateBlobRequest{
		Address: types.Address(chi.URLParam(r, "address")),
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return request, errors.Wrap(err, "failed to unmarshal")
	}
	return request, request.Validate()
}

func (r CreateBlobRequest) Validate() error {
	return ValidateStruct(&r,
		Field(&r.Address, Required),
		Field(&r.Data, Required),
	)
}

func (r CreateBlobRequestData) Validate() error {
	return ValidateStruct(&r,
		Field(&r.Type, Required),
		Field(&r.Attributes, Required),
	)
}
func (r CreateBlobRequestAttributes) Validate() error {
	return ValidateStruct(&r,
		Field(&r.Value, Required),
	)
}

func (r CreateBlobRequest) Blob() *types.Blob {
	msg := fmt.Sprintf("%s%d%s", r.Address, r.Data.Type, r.Data.Attributes.Value)
	hash := hash.Hash([]byte(msg))
	return &types.Blob{
		ID:    base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(hash[:]),
		Type:  r.Data.Type,
		Value: r.Data.Attributes.Value,
	}
}

func CreateBlob(w http.ResponseWriter, r *http.Request) {
	request, err := NewCreateBlobRequest(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	// TODO check allowed

	blob := request.Blob()

	if err := BlobQ(r).Create(request.Address, blob); err != nil {
		Log(r).WithError(err).Error("failed to save blob")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	response := CreateBlobResponse{
		Data: resources.NewBlob(blob),
	}

	w.WriteHeader(201)
	json.NewEncoder(w).Encode(&response)
}
