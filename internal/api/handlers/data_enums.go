package handlers

import (
	"net/http"
	"gitlab.com/swarmfund/api/assets"
	"encoding/json"
)

func GetDataEnums(w http.ResponseWriter, r *http.Request){
	w.Header().Set("content-type", "application/json")
	data := assets.Enums.Data()
	json.NewEncoder(w).Encode(data)
}
