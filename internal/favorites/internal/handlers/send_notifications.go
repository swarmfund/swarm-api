package handlers

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/spf13/cast"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/swarmfund/api/internal/api/handlers"
	"gitlab.com/tokend/go/doorman"
)

type NotificationLetter struct {
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

	//if no one add sale to favorite
	if len(emails) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

}
