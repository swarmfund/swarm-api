package handlers

import (
	"encoding/json"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/swarmfund/api/db2/api"
	"gitlab.com/swarmfund/api/resource"
	"io/ioutil"
	"net/http"
)

type (
	DetailsRequest struct {
		Addresses []string `json:"addresses"`
	}

	Details struct {
		Request  DetailsRequest
		Records  []api.User
		Resource resource.ShortenUsersDetails
	}
)

func GetUsersDetails(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		Log(r).WithError(err).Error("Incorrect body")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}
	var res Details
	err = json.Unmarshal(body, &res.Request)
	if err != nil {
		Log(r).WithError(err).Error("Can't unmarshal body")
		ape.RenderErr(w, problems.BadRequest(err)...)
	}

	users, err := UsersQ(r).ByAddresses(res.Request.Addresses)

	if err != nil {
		Log(r).WithError(err).Error("Can't find users")
		ape.RenderErr(w, problems.InternalError())
	}

	res.Records = users

	res.Resource.Populate(res.Records)
	json.NewEncoder(w).Encode(&res.Resource)

}
