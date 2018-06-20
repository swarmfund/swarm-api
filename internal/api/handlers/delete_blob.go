package handlers

import (
	"net/http"

	"github.com/go-chi/chi"
	. "github.com/go-ozzo/ozzo-validation"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/swarmfund/api/internal/api/movetoape"
	"gitlab.com/tokend/go/doorman"
)

type (
	DeleteBlobRequest struct {
		BlobID string `json:"-"`
	}
)

func NewDeleteBlobRequest(r *http.Request) (DeleteBlobRequest, error) {
	request := DeleteBlobRequest{
		BlobID: chi.URLParam(r, "blob"),
	}
	return request, request.Validate()
}

func (r DeleteBlobRequest) Validate() error {
	return Errors{
		"blob": Validate(&r.BlobID, Required),
	}.Filter()
}

func DeleteBlob(w http.ResponseWriter, r *http.Request) {
	request, err := NewDeleteBlobRequest(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	blob, err := BlobQ(r).Get(request.BlobID)
	if err != nil {
		ape.RenderErr(w, problems.InternalError())
		Log(r).WithError(err).Error("failed to get blob")
		return
	}

	if blob == nil || blob.DeletedAt != nil {
		// blob is not existent or already deleted
		w.WriteHeader(http.StatusNoContent)
		return
	}
	constrains := []doorman.SignerConstraint{doorman.SignerOf(CoreInfo(r).GetMasterAccountID())}
	if blob.OwnerAddress != nil {
		constrains = append(constrains, doorman.SignerOf(string(*blob.OwnerAddress)))
	}
	if err := Doorman(r, constrains...); err != nil {
		movetoape.RenderDoormanErr(w, err)
		return
	}

	if err := BlobQ(r).MarkDeleted(request.BlobID); err != nil {
		ape.RenderErr(w, problems.InternalError())
		Log(r).WithError(err).Error("failed to mark blob deleted")
	}

	w.WriteHeader(http.StatusNoContent)
}
