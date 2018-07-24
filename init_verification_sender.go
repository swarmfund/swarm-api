package api

import (
	"time"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/api/internal/clienturl"
)

func initVerificationSender(app *App) {
	go func() {
		log := logan.New()
		ticker := time.NewTicker(5 * time.Second)
		for ; ; <-ticker.C {
			err := sendVerifications(app, log)
			if err != nil {
				log.WithError(err).Error("Failed to send verifications")
			}

			err = sendWelcomeEmails(app, log)
			if err != nil {
				log.WithError(err).Error("Failed to send welcome emails")
			}
		}
	}()
}

func sendVerifications(app *App, log *logan.Entry) error {
	defer func() {
		if rvr := recover(); rvr != nil {
			log.WithRecover(rvr).Error("sendVerifications panicked")
		}
	}()

	tokensQ := app.EmailTokensQ()
	tokens, err := tokensQ.GetUnsent()
	if err != nil {
		return errors.Wrap(err, "failed to get unsent verifications")
	}

	for _, token := range tokens {
		payload := clienturl.EmailVerification(token.WalletID, token.Token)
		err = app.Config().Notificator().SendVerificationLink(token.Email, payload)
		if err != nil {
			log.WithError(err).WithField("email", token.Email).Warn("failed to send verification link")
			continue
		}

		if err := tokensQ.MarkSent(token.ID); err != nil {
			return errors.Wrap(err, "failed to mark notification as sent")
		}
	}

	return nil
}

func sendWelcomeEmails(app *App, log *logan.Entry) error {
	defer func() {
		if rvr := recover(); rvr != nil {
			log.WithRecover(rvr).Error("sendWelcomeEmails panicked")
		}
	}()

	tokensQ := app.EmailTokensQ()
	tokens, err := tokensQ.GetUnsentWelcomeEmail()
	if err != nil {
		return errors.Wrap(err, "failed to get tokens with unsent welcome emails")
	}

	for _, token := range tokens {
		err = app.Config().Notificator().SendWelcomeEmail(token.Email)
		if err != nil {
			log.WithError(err).WithField("email", token.Email).Warn("failed to send welcome email")
			continue
		}

		if err := tokensQ.MarkSentWelcomeEmail(token.ID); err != nil {
			return errors.Wrap(err, "failed to mark welcome email as sent")
		}
	}

	return nil
}

func init() {
	appInit.Add("verification-sender", initVerificationSender, "api-db")
}
