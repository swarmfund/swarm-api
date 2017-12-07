package handlers

import (
	"encoding/json"
	"net/http"

	"fmt"

	"strconv"

	"net/url"

	"github.com/go-ozzo/ozzo-validation"
	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/swarmfund/api/db2/api"
	"gitlab.com/swarmfund/api/internal/api/movetoape"
	"gitlab.com/swarmfund/api/internal/api/resources"
	"gitlab.com/swarmfund/api/internal/types"
	"gitlab.com/swarmfund/go/doorman"
)

type (
	UsersIndexResponse struct {
		Data  UsersIndexResponseData `json:"data"`
		Links Links                  `json:"links"`
	}
	UsersIndexResponseData []resources.User
	Links                  struct {
		Self string `json:"self"`
		Prev string `json:"prev,omitempty"`
		Next string `json:"next,omitempty"`
	}
	UserIndexFilters struct {
		Page  uint64
		State *uint64
	}
)

func (f *UserIndexFilters) Query() url.Values {
	query := url.Values{}
	query.Add("page", fmt.Sprintf("%d", f.Page))
	if f.State != nil {
		query.Add("state", fmt.Sprintf("%d", *f.State))
	}
	return query
}

func UsersIndex(w http.ResponseWriter, r *http.Request) {
	var err error
	query := r.URL.Query()
	filters := UserIndexFilters{
		Page: 1,
	}
	pageRaw := query.Get("page")
	if pageRaw != "" {
		filters.Page, err = strconv.ParseUint(pageRaw, 0, 64)
		if err != nil {
			ape.RenderErr(w, problems.BadRequest(validation.Errors{
				"page": errors.New("positive integer expected"),
			})...)
			return
		}
		if filters.Page == 0 {
			ape.RenderErr(w, problems.BadRequest(validation.Errors{
				"page": errors.New("should be positive"),
			})...)
			return
		}
	}
	stateRaw := query.Get("state")
	if stateRaw != "" {
		state, err := strconv.ParseUint(stateRaw, 0, 64)
		if err != nil {
			ape.RenderErr(w, problems.BadRequest(validation.Errors{
				"state": errors.New("positive integer expected"),
			})...)
			return
		}
		filters.State = &state
	}

	// TODO unhardcode
	if err := Doorman(r, doorman.SignerOf("GD7AHJHCDSQI6LVMEJEE2FTNCA2LJQZ4R64GUI3PWANSVEO4GEOWB636")); err != nil {
		movetoape.RenderDoormanErr(w, err)
		return
	}

	var records []api.User

	q := UsersQ(r).Page(filters.Page)

	if filters.State != nil {
		q = q.ByState(types.UserState(*filters.State))
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

	query = filters.Query()
	response.Links.Self = fmt.Sprintf("%s?%s", r.URL.Path, query.Encode())

	filters.Page += 1
	query = filters.Query()
	response.Links.Next = fmt.Sprintf("%s?%s", r.URL.Path, query.Encode())

	if filters.Page > 2 {
		filters.Page -= 2
		query = filters.Query()
		response.Links.Prev = fmt.Sprintf("%s?%s", r.URL.Path, query.Encode())
	}

	json.NewEncoder(w).Encode(&response)
}

type FilterLinks struct {
	Self string `json:"self"`
	Next string `json:"next"`
	Prev string `json:"prev,omitempty"`
}

func NewFilterLinks(r *http.Request, filter interface{}) {

}
