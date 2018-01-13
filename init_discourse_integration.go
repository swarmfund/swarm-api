package api

import (
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/swarmfund/api/db2/api"
	"gitlab.com/swarmfund/api/internal/discourse"
	"gitlab.com/swarmfund/api/internal/hose"
	"gitlab.com/swarmfund/go/xdr"
)

func init() {
	appInit.Add("discourse-integration", func(app *App) {
		log := app.Config().Log().WithField("service", "user-create-listener")

		connector := app.Config().Discourse()
		if connector.Disabled {
			return
		}

		// create user listener
		app.userBus.Subscribe(func(event hose.UserEvent) {
			if event.Type != hose.UserEventTypeCreated {
				return
			}
			err := connector.CreateUser(discourse.CreateUser{
				Active:   true,
				Approved: true,
				Email:    event.User.Email,
			})
			entry := log.WithField("user", event.User.Address)
			if err != nil {
				entry.WithError(err).Error("failed to create discourse user")
				return
			}
			log.Debug("discourse user created")
		})

		// create investment group
		app.txBus.Subscribe(func(event hose.TransactionEvent) {
			if event.Transaction == nil {
				return
			}
			for _, change := range event.Transaction.LedgerChanges() {
				if change.Type != xdr.LedgerEntryChangeTypeCreated {
					continue
				}
				if change.Created.Data.Type != xdr.LedgerEntryTypeAsset {
					continue
				}
				data := change.Created.Data.Asset
				log := log.WithFields(logan.F{
					"asset": data.Code,
					"tx":    event.Transaction.PagingToken,
				})
				// TODO check category exists before create
				err := connector.CreateCategory(discourse.CreateCategory{
					Name: string(data.Code),
				})
				if err != nil {
					log.WithError(err).Error("failed to create category")
					return
				}
				log.Debug("discourse category created")
			}
		})

		// investment listener
		//app.txBus.Subscribe(func(event hose.TransactionEvent) {
		//	if event.Transaction == nil {
		//		return
		//	}
		//	for _, change := range event.Transaction.LedgerChanges() {
		//		if change.Type != xdr.LedgerEntryChangeTypeUpdated {
		//			return
		//		}
		//		if change.Updated.Data.Type == xdr.LedgerEntryTypeBalance {
		//			return
		//		}
		//		data := change.Updated.Data.Balance
		//		if data.Amount > 0 {
		//			data.Asset
		//			data.AccountId
		//		}
		//	}
		//})

		// ensure all existing users processed
		func() {
			var users []api.User
			if err := app.APIQ().Users().Select(&users); err != nil {
				log.WithError(err).Error("failed to create users")
				return
			}

			for _, user := range users {
				err := connector.CreateUser(discourse.CreateUser{
					Active:   true,
					Approved: true,
					Email:    user.Email,
				})
				entry := log.WithField("user", user.Address)
				if err != nil {
					entry.WithError(err).Error("failed to create discourse user")
					return
				}
				log.Debug("discourse user created")
			}
		}()
	})
}
