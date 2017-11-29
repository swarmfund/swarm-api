package response_types

import (
	"strconv"
	"time"

	"gitlab.com/swarmfund/api/db2/api"
	"gitlab.com/swarmfund/api/render/hal"
)

const (
	OpBaseStatePending  int32 = 1
	OpBaseStateSuccess        = 2
	OpBaseStateRejected       = 3
	OpBaseStateFailed         = 4
)

var OpBaseStatuses = map[int32]string{
	OpBaseStatePending:  "Pending",
	OpBaseStateSuccess:  "Success",
	OpBaseStateRejected: "Rejected",
	OpBaseStateFailed:   "Failed",
}

type OperationI interface {
	Populate() (err error)
	ParticipantsMap() map[int64][]api.Participant
	UpdateParticipants(pMap map[int64][]api.Participant)
}

type OperationBase struct {
	Links struct {
		Self        hal.Link `json:"self"`
		Transaction hal.Link `json:"transaction"`
		Succeeds    hal.Link `json:"succeeds"`
		Precedes    hal.Link `json:"precedes"`
	} `json:"_links"`

	IdI             int64
	ID              string        `json:"id"`
	PT              string        `json:"paging_token"`
	SourceAccount   string        `json:"source_account"`
	Type            string        `json:"type"`
	TypeI           int32         `json:"type_i"`
	State           int32         `json:"state"`
	Identifier      string        `json:"identifier"`
	LedgerCloseTime time.Time     `json:"ledger_close_time"`
	Participants    []Participant `json:"participants, omitempty"`

	apiParticipants []api.Participant
}

func (b *OperationBase) Populate() (err error) {
	id, err := strconv.ParseInt(b.ID, 10, 64)
	if err != nil {
		return err
	}
	b.IdI = id
	acc := make([]api.Participant, len(b.Participants))
	for i, p := range b.Participants {
		acc[i] = p.ToApiParticipant()
	}

	b.apiParticipants = acc
	return
}

func (b *OperationBase) UpdateParticipants(pMap map[int64][]api.Participant) {
	acc := pMap[b.IdI]

	if len(pMap[b.IdI]) > len(b.Participants) {
		tail := make([]Participant, len(pMap[b.IdI])-len(b.Participants))
		participants := append(b.Participants, tail...)
		b.Participants = participants
	}

	for i, p := range pMap[b.IdI] {
		b.Participants[i].FromApiParticipant(&p)
		acc[i] = b.Participants[i].ToApiParticipant()
	}

	b.apiParticipants = acc
	return
}

func (b *OperationBase) ParticipantsMap() map[int64][]api.Participant {
	return map[int64][]api.Participant{
		b.IdI: b.apiParticipants,
	}
}
