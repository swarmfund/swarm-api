package track

import (
	"net/http"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/tokend/go/signcontrol"
)

type EventType int32

const (
	EventTypeGetWallet EventType = 1 << iota
)

type Event struct {
	Address string
	Signer  string
	Details EventDetails
}

type EventDetails struct {
	Type    EventType
	Request *RequestDetails
}

type RequestDetails struct {
	IP        string
	UserAgent string
	Path      string
}

type Q interface {
	Track(Event) error
}

type Tracker struct {
	entry *logan.Entry
	q     Q
}

func NewTracker(entry *logan.Entry, q Q) *Tracker {
	return &Tracker{
		entry, q,
	}
}

func (t *Tracker) track(event Event) {
	defer func() {
		if rvr := recover(); rvr != nil {
			t.entry.WithRecover(rvr).Error("tracker panicked")
		}
	}()
	if err := t.q.Track(event); err != nil {
		t.entry.WithError(err).WithFields(logan.F{
			"event": event,
		}).Error("failed to save event")
	}
}
func (t *Tracker) GetWallet(request *http.Request) {
	sig, _ := signcontrol.IsSigned(request)
	t.track(Event{
		Address: sig.Address,
		Signer:  sig.Signer,
		Details: EventDetails{
			Type: EventTypeGetWallet,
			Request: &RequestDetails{
				IP:        request.RemoteAddr,
				UserAgent: request.Header.Get("user-agent"),
			},
		},
	})
}
