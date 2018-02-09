package api

import (
	"time"

	"gitlab.com/distributed_lab/figure"
	"gitlab.com/swarmfund/api/db2/api"
	"gitlab.com/swarmfund/api/internal/data"
	"gitlab.com/swarmfund/api/log"
)

func deleteExpiredWallets(log *log.Entry, emailTokensQ data.EmailTokensQ, walletQI api.WalletQI, expireDuration time.Duration) {
	for ; ; time.Sleep(expireDuration / 10) {
		tokens, err := emailTokensQ.GetUnconfirmed()
		if err != nil {
			log.WithError(err).Error("unable to get unconfirmed email tokens")
			continue
		}

		if len(tokens) == 0 {
			continue
		}

		expiredWallets := []string{}
		for _, token := range tokens {
			if token.LastSentAt.Add(expireDuration).Before(time.Now()) {
				expiredWallets = append(expiredWallets, token.WalletID)
			}
		}

		err = walletQI.DeleteWallets(expiredWallets)
		if err != nil {
			log.WithError(err).Error("unable to delete expired wallets")
			continue
		}

		log.WithField("count", len(expiredWallets)).Info("wallets deleted")
	}
}

func initWalletCleaner(app *App) {
	service := "wallet_cleaner"
	entry := log.WithField("service", service)
	var config struct {
		Enabled    bool
		Expiration time.Duration
	}
	err := figure.
		Out(&config).
		From(app.config.Get(service)).
		Please()
	if err != nil {
		entry.WithError(err).Error("failed to figure out")
		return
	}

	if !config.Enabled {
		entry.Info("disabled")
		return
	}

	//go deleteExpiredWallets(entry, app.EmailTokensQ(), app.APIQ().Wallet(), config.Expiration)
}

func init() {
	appInit.Add("wallet-cleaner", initWalletCleaner, "api-db")
}
