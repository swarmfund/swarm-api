package api

import (
	"time"

	"gitlab.com/swarmfund/api/db2/api"
	"gitlab.com/swarmfund/api/internal/data"
	"gitlab.com/swarmfund/api/log"
)

func deleteExpiredWallets(emailTokensQ data.EmailTokensQ, walletQI api.WalletQI, expireDuration time.Duration) {
	var emTokens []data.EmailToken
	var err error
	logger := log.WithField("service", "wallet_cleaner")

	for {
		time.Sleep(expireDuration / 10)

		emTokens, err = emailTokensQ.GetUnconfirmed()
		if err != nil {
			logger.WithError(err).Error("unable to get unconfirmed email tokens")
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
			logger.WithError(err).Error("Unable to delete expired wallets")
			continue
		}

		logger.WithField("quantity", len(walletIDs)).Debug("Successfully deleted wallets")
	}
}

func initWalletCleaner(app *App) {
	cfg := app.config.WalletCleaner()
	if !cfg.Enabled {
		return
	}

	expireDuration, err := time.ParseDuration(cfg.Expiration)
	if err != nil {
		log.WithField("error", err.Error()).Error("failed to init wallet_cleaner")
		return
	}

	go deleteExpiredWallets(app.EmailTokensQ(), app.APIQ().Wallet(), expireDuration)
}

func init() {
	appInit.Add("wallet-cleaner", initWalletCleaner, "api-db")
}
