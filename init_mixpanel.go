package api

import (
	"time"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/swarmfund/api/internal/hose"
	"gitlab.com/tokend/go/xdr"
)

func init() {
	appInit.Add("mixpanel", func(app *App) {
		// only events occurred after cursor will be submitted
		cursor := time.Now()
		mixpanel := app.Config().Mixpanel()
		horizon := app.Config().Horizon()
		log := app.Config().Log().WithField("service", "mixpanel")

		app.txBus.Subscribe(func(event hose.TransactionEvent) {
			if event.Transaction == nil {
				return
			}

			if event.Transaction.CreatedAt.Before(cursor) {
				// event is before cursor, skipping
				return
			}

			envelope := event.Transaction.Envelope()
			for _, op := range envelope.Tx.Operations {
				switch op.Body.Type {
				case xdr.OperationTypeCreateIssuanceRequest:
					body := op.Body.CreateIssuanceRequestOp
					fields := logan.F{
						"type":     "issuance-request",
						"request":  body.Reference,
						"receiver": body.Request.Receiver.AsString(),
					}
					address, err := horizon.Accounts().ByBalance(body.Request.Receiver.AsString())
					if err != nil {
						log.WithError(err).WithFields(fields).Error("failed to get address")
						continue
					}
					if address == nil {
						log.WithFields(fields).Error("address not found")
						continue
					}
					if err := mixpanel.IssuanceRequest(*address, &event.Transaction.CreatedAt, body); err != nil {
						log.WithError(err).WithFields(fields).Error("failed to submit event")
						continue
					}
				case xdr.OperationTypeCreateWithdrawalRequest:
					body := op.Body.CreateWithdrawalRequestOp
					fields := logan.F{
						"type":      "withdrawal-request",
						"requester": body.Request.Balance.AsString(),
					}
					address, err := horizon.Accounts().ByBalance(body.Request.Balance.AsString())
					if err != nil {
						log.WithError(err).WithFields(fields).Error("failed to get address")
						continue
					}
					if address == nil {
						log.WithFields(fields).Error("address not found")
						continue
					}
					if err := mixpanel.WithdrawalRequest(*address, &event.Transaction.CreatedAt, body); err != nil {
						log.WithError(err).WithFields(fields).Error("failed to submit event")
						continue
					}
				}
			}
		})
	}, "tx-watcher")
}
