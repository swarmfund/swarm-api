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
	"gitlab.com/swarmfund/api/internal/api/movetoape"
	"gitlab.com/swarmfund/api/internal/api/resources"
	"gitlab.com/swarmfund/api/internal/data"
	"gitlab.com/swarmfund/api/internal/data/postgres"
	"gitlab.com/swarmfund/api/internal/types"
	"gitlab.com/tokend/go/doorman"
	"gitlab.com/tokend/go/hash"
)

type (
	CreateBlobRequest struct {
		Address *types.Address        `json:"-"`
		Data    CreateBlobRequestData `json:"data"`
	}
	CreateBlobRequestData struct {
		Type          types.BlobType              `json:"type"`
		Attributes    CreateBlobRequestAttributes `json:"attributes"`
		Relationships Relationships               `json:"relationships"`
	}
	CreateBlobRequestAttributes struct {
		Value string `json:"value"`
	}
	CreateBlobResponse struct {
		Data resources.Blob `json:"data"`
	}
	Relationships map[string]Object
	Object        struct {
		Data ObjectData `json:"data"`
	}
	ObjectData struct {
		ID string `json:"id"`
	}
)

func (r ObjectData) Validate() error {
	return ValidateStruct(&r,
		Field(&r.ID, Required),
	)
}

func (r Object) Validate() error {
	return ValidateStruct(&r,
		Field(&r.Data, Required),
	)
}

func (r Relationships) Validate() error {
	errs := Errors{}
	for k, v := range r {
		errs[k] = Validate(v.Data, Required)
	}
	return errs.Filter()
}

func NewCreateBlobRequest(r *http.Request) (CreateBlobRequest, error) {
	request := CreateBlobRequest{}
	address := types.Address(chi.URLParam(r, "address"))
	if address != "" {
		request.Address = &address
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return request, errors.Wrap(err, "failed to unmarshal")
	}
	return request, request.Validate()
}

func (r CreateBlobRequest) Validate() error {
	return Errors{
		"address": Validate(&r.Address),
		"/data":   Validate(&r.Data, Required),
	}.Filter()
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
	var address types.Address
	if r.Address != nil {
		address = *r.Address
	}
	msg := fmt.Sprintf("%s%d%s", address, r.Data.Type, r.Data.Attributes.Value)
	hash := hash.Hash([]byte(msg))

	relationships := types.BlobRelationships{}
	for k, v := range r.Data.Relationships {
		relationships[k] = v.Data.ID
	}

	return &types.Blob{
		ID:            base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(hash[:]),
		Type:          r.Data.Type,
		Value:         r.Data.Attributes.Value,
		Relationships: relationships,
	}
}

func CreateBlob(w http.ResponseWriter, r *http.Request) {
	request, err := NewCreateBlobRequest(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	var constraint doorman.SignerConstraint
	if request.Address != nil {
		// if address is set blob has an owner
		constraint = doorman.SignerOf(string(*request.Address))
	} else {
		// nil address means blob is not bound to any user
		constraint = doorman.SignerOf(CoreInfo(r).GetMasterAccountID())
	}
	if err := Doorman(r, constraint); err != nil {
		movetoape.RenderDoormanErr(w, err)
		return
	}

	blob := request.Blob()

	err = BlobQ(r).Transaction(func(blobs data.Blobs) error {
		if blob.Type == types.BlobTypeNavUpdate {
			if request.Address == nil {
				// not sure why NavUpdate case even exists but master should not use it anyways
				return errors.New("master not allowed to create nav update")
			}
			existing, err := blobs.
				ByOwner(*request.Address).
				ByType(blob.Type).
				ByRelationships(blob.Relationships).
				Select()
			if err != nil {
				return errors.Wrap(err, "failed to get existing blobs")
			}

			if err := blobs.Delete(existing...); err != nil {
				return errors.Wrap(err, "failed to delete existing blobs")
			}
		}
		if err := blobs.Create(request.Address, blob); err != nil {
			return errors.Wrap(err, "failed to create blob")
		}
		return nil
	})
	if err != nil {
		// silencing error to make request idempotent
		if errors.Cause(err) != postgres.ErrBlobsConflict {
			Log(r).WithError(err).Error("failed to save blob")
			ape.RenderErr(w, problems.InternalError())
			return
		}
	}

	response := CreateBlobResponse{
		Data: resources.NewBlob(blob),
	}

	w.WriteHeader(201)
	json.NewEncoder(w).Encode(&response)

	go Tracker(r).CreateBlob(r)
}
