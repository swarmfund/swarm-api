package handlers

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	res "gitlab.com/swarmfund/api/resource"
	"gitlab.com/swarmfund/api/db2/api"
)

type (
	DetailsRequest struct {
		Addresses []string `json:"addresses"`
	}

	DetailsResp struct {
		Request DetailsRequest
		Records  []api.User
		Resource res.ShortenUsersDetails
	}
)

func PostDetails(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		Log(r).WithError(err).Error("Incorrect body")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	var request DetailsRequest
	err = json.Unmarshal(body, &request)
	if err != nil{
		Log(r).WithError(err).Error("Can't unmarshal body")
		ape.RenderErr(w, problems.BadRequest(err)...)
	}

	users := make([]api.User, 0)

	for _, addr := range request.Addresses{
		user, err := UsersQ(r).ByAddress(addr)
		if err != nil{
			Log(r).WithError(err).Error("failed to get users")
			ape.RenderErr(w, problems.InternalError())
			return
		}
		users = append(users, *user)
	}

	var response DetailsResp
	response.Records = users

	response.Resource.Populate(response.Records)
	json.NewEncoder(w).Encode(&response.Resource)

}