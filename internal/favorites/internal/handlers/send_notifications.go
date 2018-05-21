package handlers

import (
	"net/http"

	"encoding/json"
	"io/ioutil"

	"github.com/go-chi/chi"
	"github.com/spf13/cast"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/swarmfund/api/internal/api/handlers"
	"gitlab.com/tokend/go/doorman"
)

type NotificationLetter struct {
	Message string `json:"message"`
}

//SendNotifications check signature of sale owner
//if ok, takes from DB all emails of users which add to favorites sale with a key in params
//then send to notificator-server
func SendNotifications(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "key")
	saleID, err := cast.ToUint64E(key)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	sale, err := handlers.Horizon(r).Sales().SaleByID(saleID)
	if err != nil {
		handlers.Log(r).WithError(err).Error("failed to get sale")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	if err := handlers.Doorman(r, doorman.SignatureOf(sale.Owner)); err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	emails, err := FavoritesQ(r).GetEmails(key)
	if err != nil {
		handlers.Log(r).WithError(err).Error("failed to load users emails")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	//if no one add sale to favorite or no such sale
	if len(emails) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	defer r.Body.Close()
	//read message from request
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	letter := &NotificationLetter{}
	if err := json.Unmarshal(body, letter); err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	err = handlers.Notificator(r).SendSaleNotifications(emails, letter.Message)
	if err != nil {
		handlers.Log(r).WithError(err).Error("failed to send notifications")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
