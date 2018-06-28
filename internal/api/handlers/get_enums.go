package handlers

import (
	"encoding/json"
	"net/http"

	"gitlab.com/swarmfund/api/assets"
)

func GetEnums(w http.ResponseWriter, r *http.Request) {
	response := assets.Enums.Data()

	json.NewEncoder(w).Encode(response)
}
