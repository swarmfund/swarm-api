package notificator

import (
	"bytes"
	"fmt"
	"html/template"
	"net/url"

	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/notificator"
	"gitlab.com/swarmfund/api/internal/clienturl"
)

const (
	NotificationTypeVerificationEmail = 3
)

type Config struct {
	Disabled     bool
	Endpoint     string
	Secret       string
	Public       string
	ClientRouter string

	EmailConfirmation *template.Template
	KYCApprove        *template.Template
	KYCReject         *template.Template
}

type Connector struct {
	notificator *notificator.Connector
	conf        Config
}

func NewConnector(conf Config) *Connector {
	// TODO move this to config
	endpoint, err := url.Parse(conf.Endpoint)
	if err != nil {
		panic(err)
	}

	return &Connector{
		notificator: notificator.NewConnector(
			notificator.Pair{Secret: conf.Secret, Public: conf.Public},
			*endpoint,
		),
		conf: conf,
	}
}

func (c *Connector) SendVerificationLink(email string, payload clienturl.Payload) error {
	encoded, err := payload.Encode()
	if err != nil {
		return errors.Wrap(err, "failed to encode payload")
	}
	letter := &Letter{
		Header: "Swarm Email Verification",
		Link:   fmt.Sprintf("%s/%s", c.conf.ClientRouter, encoded),
	}

	var buff bytes.Buffer
	err = c.conf.EmailConfirmation.Execute(&buff, letter)
	if err != nil {
		return errors.Wrap(err, "Error while populating template for notify approval")
	}
	msg := &notificator.EmailRequestPayload{
		Destination: email,
		Subject:     letter.Header,
		Message:     buff.String(),
	}

	return c.send(NotificationTypeVerificationEmail, email, msg)
}

func (c *Connector) NotifyApproval(email string) error {
	letter := &Letter{
		Header: "Swarm Fund KYC Review",
		Link:   c.conf.ClientRouter,
	}

	var buff bytes.Buffer
	if err := c.conf.KYCApprove.Execute(&buff, letter); err != nil {
		return errors.Wrap(err, "failed to render template")
	}
	msg := &notificator.EmailRequestPayload{
		Destination: email,
		Subject:     letter.Header,
		Message:     buff.String(),
	}

	return c.send(NotificationTypeVerificationEmail, email, msg)
}

func (c *Connector) NotifyRejection(email string) error {
	letter := &Letter{
		Header: "Swarm Fund KYC Review",
		Link:   c.conf.ClientRouter,
	}

	var buff bytes.Buffer
	if err := c.conf.KYCReject.Execute(&buff, letter); err != nil {
		return errors.Wrap(err, "failed to render template")
	}
	msg := &notificator.EmailRequestPayload{
		Destination: email,
		Subject:     letter.Header,
		Message:     buff.String(),
	}

	return c.send(NotificationTypeVerificationEmail, email, msg)
}

func (c *Connector) send(requestType int, token string, payload notificator.Payload) error {
	if c.conf.Disabled {
		// TODO log warning
		return nil
	}

	response, err := c.notificator.Send(requestType, token, payload)
	if err != nil {
		return err
	}

	if !response.IsSuccess() {
		return errors.New("notification request not accepted")
	}
	return nil
}
