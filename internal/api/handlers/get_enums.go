package handlers

import (
	"encoding/json"
	"gitlab.com/swarmfund/api/assets"
	"net/http"
)


func GetEnums(w http.ResponseWriter, r *http.Request) {
	response := assets.Enums.Data()

	json.NewEncoder(w).Encode(&response)
}
