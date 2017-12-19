package api

import (
	"time"

	"gitlab.com/distributed_lab/figure"
	"gitlab.com/swarmfund/api/db2/api"
	"gitlab.com/swarmfund/api/internal/data"
	"gitlab.com/swarmfund/api/log"
)

func deleteExpiredWallets(log *log.Entry, emailTokensQ data.EmailTokensQ, walletQI api.WalletQI, expireDuration time.Duration) {
	var emTokens []data.EmailToken
	var err error
	for {
		time.Sleep(expireDuration / 10)

		emTokens, err = emailTokensQ.GetUnconfirmed()
		if err != nil {
			log.WithError(err).Error("unable to get unconfirmed email tokens")
			continue
		}

		if len(emTokens) == 0 {
			continue
		}

		walletIDs := []string{}
		for _, emt := range emTokens {
			if emt.LastSentAt.UTC().Before(time.Now().UTC().Add(-1 * expireDuration)) {
				walletIDs = append(walletIDs, emt.WalletID)
			}
		}

		if len(walletIDs) == 0 {
			continue
		}

		err = walletQI.DeleteWallets(walletIDs)
		if err != nil {
			log.WithError(err).Error("Unable to delete expired wallets")
			continue
		}

		log.WithField("quantity", len(walletIDs)).Debug("Successfully deleted wallets")
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

	go deleteExpiredWallets(entry, app.EmailTokensQ(), app.APIQ().Wallet(), config.Expiration)
}

func init() {
	appInit.Add("wallet-cleaner", initWalletCleaner, "api-db")
}
