package kyctracker

import (
	"time"

	"fmt"

	"gitlab.com/swarmfund/api/db2"
	"gitlab.com/swarmfund/api/db2/api"
	. "gitlab.com/swarmfund/api/errors"
	"gitlab.com/swarmfund/api/log"
	"gitlab.com/swarmfund/api/notificator"
	"github.com/pkg/errors"
)

type Tracker struct {
	log         *log.Entry
	q           api.QInterface
	notificator *notificator.Connector
}

func New(entry *log.Entry, q api.QInterface, connector *notificator.Connector) *Tracker {
	return &Tracker{
		log:         entry,
		q:           q,
		notificator: connector,
	}
}

func (tracker *Tracker) Run() {
	ticker := time.NewTicker(10 * time.Second)
	closeCh := make(chan struct{})
	for ; ; <-ticker.C {
		tracker.log.Info("checking kyc states")
		pending := false
		for user := range tracker.iterateUsers(closeCh) {
			ok, err := tracker.processUser(&user)
			if ok {
				pending = true
			}
			if err != nil {
				tracker.log.WithError(err).Error("failed to process user")
			}
		}
		if pending {
			if err := tracker.sendNotifications(); err != nil {
				tracker.log.WithError(err).Error("failed to send notifications")
				continue
			}
			tracker.log.Info("kyc notification has been sent")
		}
	}
}

func (tracker *Tracker) sendNotifications() error {
	recipients, err := tracker.q.Notifications().GetRecipients(api.NotificationTypeKYC)
	if err != nil {
		return errors.Wrap(err, "failed to get recipients")
	}
	for _, recipient := range recipients {
		if err := tracker.notificator.NotifyKYCReviewPending(recipient); err != nil {
			return errors.Wrap(err, "failed to send notification")
		}
	}
	return nil
}

func (tracker *Tracker) processUser(user *api.User) (bool, error) {
	if user.State != api.UserWaitingForApproval {
		// we are currently interested in users waiting for approval
		return false, nil
	}

	changed, err := tracker.q.KYCTracker().EnsureRecord(user.ID, user.State)
	if err != nil {
		return false, errors.Wrap(err, "failed to ensure record")
	}

	return changed, nil
}

func (tracker *Tracker) iterateUsers(closeCh <-chan struct{}) <-chan api.User {
	results := make(chan api.User)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				tracker.log.WithError(FromPanic(r)).Error("iterate panicked")
			}
			close(results)
		}()
		cursor := db2.PageQuery{
			Order: "desc",
			Limit: 100,
		}
		for {
			users := []api.User{}
			// TODO optimize with updated_at filter
			if err := tracker.q.Users().Page(cursor).Select(&users); err != nil {
				tracker.log.WithError(err).Error("failed to get users")
				return
			}

			if len(users) == 0 {
				// end of iteration
				return
			}
			for _, user := range users {
				select {
				case <-closeCh:
					// received close signal
					return
				default:
					results <- user
					cursor.Cursor = fmt.Sprintf("%d", user.ID)
				}
			}
		}
	}()
	return results
}
