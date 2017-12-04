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
	"gitlab.com/swarmfund/api/internal/api/resources"
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
		Page uint64
	}
)

func (f *UserIndexFilters) Query() url.Values {
	query := url.Values{}
	query.Add("page", fmt.Sprintf("%d", f.Page))
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
				"page": errors.New("integer expected"),
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

	// TODO check allowed

	var records []api.User
	if err := UsersQ(r).Page(filters.Page).Select(&records); err != nil {
		Log(r).WithError(err).Error("failed to get users")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	response := UsersIndexResponse{
		Data: make(UsersIndexResponseData, 0, len(records)),
	}
	for _, record := range records {
		response.Data = append(response.Data, resources.User{
			Type: string(record.UserType),
			ID:   record.Address,
			Attributes: resources.UserAttributes{
				Email: record.Email,
			},
		})
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
