package api

import (
	"time"

	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/api/db2/api"
	"gitlab.com/swarmfund/api/internal/data"
	"gitlab.com/swarmfund/api/log"
)

const (
	walletCleanerService = "wallet_cleaner"
)

func deleteExpiredWallets(log *log.Entry, emailTokensQ data.EmailTokensQ, walletQI api.WalletQI, expireDuration time.Duration) {
	do := func() (err error) {
		defer func() {
			if rvr := recover(); rvr != nil {
				err = errors.FromPanic(err)
			}
		}()
		tokens, err := emailTokensQ.GetUnconfirmed()
		if err != nil {
			return errors.Wrap(err, "failed to get email tokens")
		}

		if len(tokens) == 0 {
			return nil
		}

		var expiredWallets []string
		for _, token := range tokens {
			if token.LastSentAt == nil {
				// email has not been sent yet
				continue
			}
			if token.LastSentAt.Add(expireDuration).Before(time.Now()) {
				expiredWallets = append(expiredWallets, token.WalletID)
			}
		}

		err = walletQI.DeleteWallets(expiredWallets)
		if err != nil {
			return errors.Wrap(err, "failed to delete wallets")
		}

		log.WithField("count", len(expiredWallets)).Info("wallets deleted")

		return nil
	}

	for ; ; time.Sleep(expireDuration / 10) {
		if err := do(); err != nil {
			log.WithError(err).Error("failed to clean wallets")
		}
	}
}

func initWalletCleaner(app *App) {
	entry := log.WithField("service", walletCleanerService)
	var config struct {
		Enabled    bool
		Expiration time.Duration
	}
	err := figure.
		Out(&config).
		From(app.config.Get(walletCleanerService)).
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
	appInit.Add(walletCleanerService, initWalletCleaner, "api-db")
}
