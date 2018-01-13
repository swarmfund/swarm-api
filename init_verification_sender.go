package api

import (
	"time"

	"gitlab.com/swarmfund/api/internal/clienturl"
)

func initVerificationSender(app *App) {
	go func() {
		tokensQ := app.EmailTokensQ()
		ticker := time.NewTicker(5 * time.Second)
		for ; ; <-ticker.C {
			tokens, err := tokensQ.GetUnsent()
			if err != nil {
				panic(err)
			}
			for _, token := range tokens {
				payload := clienturl.EmailVerification(token.WalletID, token.Token)
				err = app.Config().Notificator().SendVerificationLink(token.Email, payload)
				if err != nil {
					panic(err)
				}
				if err := tokensQ.MarkSent(token.ID); err != nil {
					panic(err)
				}
			}
		}
	}()
}

func init() {
	appInit.Add("verification-sender", initVerificationSender, "api-db")
}
