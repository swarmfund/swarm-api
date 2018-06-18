package track

import (
	"net/http"

	"encoding/json"

	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/swarmfund/api/actions"
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

func (d *EventDetails) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	if bytes, ok := value.([]byte); ok {
		if err := json.Unmarshal(bytes, &d); err != nil {
			return errors.Wrap(err, "failed to scan EventDetails")
		}
		return nil
	}
	return errors.New("failed to scan EventDetails")
}

type RequestDetails struct {
	IP        string
	UserAgent string
	Path      string
	Location  string
}

type Q interface {
	Track(Event) error
	Last(*Event) (*Event, error)
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

func (t *Tracker) CreateBlob(address string, request *http.Request) {
	ip := request.Header.Get("x-real-ip")

	location, err := actions.GetGeoLocation(ip)
	if err != nil {
		t.entry.WithError(err).WithFields(logan.F{
			"ip": ip,
		}).Error("failed to get geolocation")
	}

	sig, _ := signcontrol.IsSigned(request)
	t.track(Event{
		Address: address,
		Signer:  sig.Signer,
		Details: EventDetails{
			Type: EventTypeGetWallet,
			Request: &RequestDetails{
				IP:        ip,
				UserAgent: request.Header.Get("user-agent"),
				Location:  location,
			},
		},
	})
}

func (t *Tracker) GetLast(event Event) (*Event, error) {
	return t.q.Last(&event)
}
