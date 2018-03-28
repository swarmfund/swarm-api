package handlers

import (
	"net/http"

	"encoding/json"

	"github.com/go-chi/chi"
	. "github.com/go-ozzo/ozzo-validation"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/swarmfund/api/internal/api/movetoape"
	"gitlab.com/swarmfund/api/internal/api/resources"
	"gitlab.com/swarmfund/api/internal/types"
	"gitlab.com/swarmfund/go/doorman"
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
	fieldRule := &FieldRules{}
	if r.Address == "" {
		fieldRule = Field(&r.Address)
	} else {
		fieldRule = Field(&r.Address, Required)
	}

	return ValidateStruct(&r,
		fieldRule,
		Field(&r.BlobID, Required),
	)
}
func GetBlob(w http.ResponseWriter, r *http.Request) {
	request, err := NewGetBlobRequest(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

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

	if !types.IsPublicBlob(blob.Type) {
		if request.Address == "" {
			err := Doorman(r, doorman.SignerOf(CoreInfo(r).GetMasterAccountID()))
			if err != nil {
				movetoape.RenderDoormanErr(w, err)
				return
			}
		} else {
			err := Doorman(r, doorman.SignerOf(string(request.Address)), doorman.SignerOf(CoreInfo(r).GetMasterAccountID()))
			if err != nil {
				movetoape.RenderDoormanErr(w, err)
				return
			}
		}
	}

	response := CreateBlobResponse{
		Data: resources.NewBlob(blob),
	}

	json.NewEncoder(w).Encode(&response)
}
