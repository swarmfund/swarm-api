package handlers

import (
	"encoding/json"
	"net/http"

	"gitlab.com/swarmfund/api/assets"
)

func GetDataEnums(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	data := assets.Enums.Data()
	json.NewEncoder(w).Encode(data)
}
