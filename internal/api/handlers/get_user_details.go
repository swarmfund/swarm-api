package handlers

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/swarmfund/api/resource"
	"gitlab.com/swarmfund/api/db2/api"
)

type (
	DetailsRequest struct {
		Addresses []string `json:"addresses"`
	}

	Details struct {
		Request DetailsRequest
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
	if err != nil{
		Log(r).WithError(err).Error("Can't unmarshal body")
		ape.RenderErr(w, problems.BadRequest(err)...)
	}

	users := make([]api.User, 0)
	//TODO single query
	for _, addr := range res.Request.Addresses{
		user, err := UsersQ(r).ByAddress(addr)
		if err != nil{
			Log(r).WithError(err).Error("failed to get users")
			ape.RenderErr(w, problems.InternalError())
			return
		}
		users = append(users, *user)
	}


	res.Records = users

	res.Resource.Populate(res.Records)
	json.NewEncoder(w).Encode(&res.Resource)

}