package handlers

import (
	"encoding/json"
	"net/http"

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
	BlobDirectRequest struct {
		Hash string `json:"-"`
	}
)

func NewBlobDirectRequest(r *http.Request) (BlobDirectRequest, error) {
	request := BlobDirectRequest{
		Hash: chi.URLParam(r, "blob"),
	}
	return request, request.Validate()
}

func (r BlobDirectRequest) Validate() error {
	return ValidateStruct(&r,
		Field(&r.Hash),
	)
}
func GetDirectBlob(w http.ResponseWriter, r *http.Request) {
	request, err := NewBlobDirectRequest(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	blob, err := BlobQ(r).Get(request.Hash)
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
		err := Doorman(r, doorman.SignerOf(CoreInfo(r).GetMasterAccountID()))
		if err != nil {
			movetoape.RenderDoormanErr(w, err)
			return
		}
	}

	response := CreateBlobResponse{
		Data: resources.NewBlob(blob),
	}

	json.NewEncoder(w).Encode(&response)
}
