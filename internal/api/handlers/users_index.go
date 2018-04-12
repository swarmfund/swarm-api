package handlers

import (
	"encoding/json"
	"net/http"

	. "github.com/go-ozzo/ozzo-validation"
	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/swarmfund/api/db2/api"
	"gitlab.com/swarmfund/api/internal/api/movetoape"
	"gitlab.com/swarmfund/api/internal/api/resources"
	"gitlab.com/swarmfund/api/internal/api/urlval"
	"gitlab.com/swarmfund/api/internal/types"
	"gitlab.com/swarmfund/go/doorman"
)

type (
	UsersIndexResponse struct {
		Data  UsersIndexResponseData `json:"data"`
		Links urlval.FilterLinks     `json:"links"`
	}
	UsersIndexResponseData []resources.User
	UserIndexFilters       struct {
		Page    uint64  `url:"page"`
		State   *uint64 `url:"state"`
		Type    *uint64 `url:"type"`
		Email   *string `url:"email"`
		Address *string `url:"address"`

		//Relationships
		FirstName *string `url:"first_name"`
		LastName  *string `url:"last_name"`
		Country   *string `url:"country"`
	}
)

func NewUserFilters(r *http.Request) (UserIndexFilters, error) {
	filters := UserIndexFilters{
		Page: 1,
	}
	if err := urlval.Decode(r.URL.Query(), &filters); err != nil {
		return filters, errors.Wrap(err, "failed to populate")
	}
	return filters, filters.Validate()
}

func (r UserIndexFilters) Validate() error {
	return ValidateStruct(&r,
		Field(&r.Page, Min(uint64(1))),
		Field(&r.State, Min(uint64(1))),
		Field(&r.Type, Min(uint64(1))),
	)
}

func UsersIndex(w http.ResponseWriter, r *http.Request) {
	filters, err := NewUserFilters(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	if err := Doorman(r, doorman.SignerOf(CoreInfo(r).GetMasterAccountID())); err != nil {
		movetoape.RenderDoormanErr(w, err)
		return
	}

	var records []api.User

	q := UsersQ(r).Page(filters.Page)

	if filters.State != nil {
		q = q.ByState(types.UserState(*filters.State))
	}

	if filters.Type != nil {
		q = q.ByType(types.UserType(*filters.Type))
	}

	if filters.Email != nil {
		q = q.EmailMatches(*filters.Email)
	}

	if filters.Address != nil {
		q = q.AddressMatches(*filters.Address)
	}

	if filters.FirstName != nil {
		q = q.ByFirstName(*filters.FirstName)
	}

	if filters.LastName != nil {
		q = q.ByLastName(*filters.LastName)
	}

	if filters.Country != nil {
		q = q.ByCountry(*filters.Country)
	}

	if err := q.Select(&records); err != nil {
		Log(r).WithError(err).Error("failed to get users")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	response := UsersIndexResponse{
		Data: make(UsersIndexResponseData, 0, len(records)),
	}
	for _, record := range records {
		response.Data = append(response.Data, resources.NewUser(&record))
	}

	response.Links = urlval.Encode(r, filters)

	json.NewEncoder(w).Encode(&response)
}
