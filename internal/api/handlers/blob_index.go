package handlers

import (
	"net/http"

	"encoding/json"

	"github.com/go-chi/chi"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/swarmfund/api/internal/api/resources"
)

func BlobIndex(w http.ResponseWriter, r *http.Request) {
	filters := map[string]string{}
	for k := range r.URL.Query() {
		filters[k] = r.URL.Query().Get(k)
	}

	address := chi.URLParam(r, "address")

	records, err := BlobQ(r).Filter(address, filters)
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
