package notificator

import (
	"encoding/base64"
	"fmt"
	"net/url"

	"gitlab.com/tokend/go/hash"
)

const TimeLayout = "2006-01-02 15:04:05"

type TransferTypeLetter int

const (
	PAYMENT   = 1 + iota
	INVOICE   = 1 + iota
	OFFER     = 1 + iota
	FORFEIT   = 1 + iota
	DEPOSIT   = 1 + iota
	DEMURRAGE = 1 + iota
)

type operationNoticeLetter struct {
	Id        string
	Header    string
	Email     string
	Addressee string
	Amount    string
	Asset     string
	Date      string
	Link      string
	Message   string
	Type      string
}

type CoinsEmissionNoticeLetter struct {
	operationNoticeLetter
	Status       string
	RejectReason string
}

type DemurrageNoticeLetter struct {
	operationNoticeLetter
	PeriodFrom string
	PeriodTo   string
}

type ForfeitNoticeLetter struct {
	operationNoticeLetter
	Status string
}

type InvoiceNoticeLetter struct {
	operationNoticeLetter
	Counterparty     string
	CounterpartyType string
	Status           string
	RejectReason     string
}
type OfferNoticeLetter struct {
	operationNoticeLetter
	Price       string
	OrderPrice  string
	QuoteAmount string
	Fee         string
	Direction   string
}

type PaymentNoticeLetter struct {
	operationNoticeLetter

	Action           string
	Counterparty     string
	CounterpartyType string
	Fee              string
	FullAmount       string
	Reference        string
}

type TransferNoticeI interface {
	AddLoginLink(clientURL string) error
	GetEmail() string
	GetHeader() string
	GetToken() string
}

func (l *operationNoticeLetter) AddLoginLink(clientURL string) error {
	link, err := url.Parse(fmt.Sprintf("%s/login", clientURL))
	if err != nil {
		return err
	}
	query := link.Query()
	query.Set("username", l.Email)
	link.RawQuery = query.Encode()
	l.Link = link.String()
	return nil
}

func (l *operationNoticeLetter) GetHeader() string {
	return l.Header
}

func (l *operationNoticeLetter) GetEmail() string {
	return l.Email
}

func (l *operationNoticeLetter) GetToken() string {
	hash := hash.Hash([]byte(l.Id))
	return base64.URLEncoding.EncodeToString(hash[:])
}
