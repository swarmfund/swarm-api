package handlers

import (
	"encoding/json"
	"gitlab.com/swarmfund/api/assets"
	"net/http"
)

type GetEnumsResponse struct {
	Data map[string]interface{} `json:"data"`
}

func GetEnums(w http.ResponseWriter, r *http.Request) {

	response := GetEnumsResponse{
		Data: assets.Enums.Data(),
	}

	defer json.NewEncoder(w).Encode(&response)
}
