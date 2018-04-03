package notificator

import (
	"bytes"
	"fmt"
	"html/template"
	"net/url"

	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/notificator"
	"gitlab.com/swarmfund/api/internal/clienturl"
	"gitlab.com/swarmfund/horizon-connector/v2"
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

func (c *Connector) Init(connector *horizon.Connector) error {
	emailConfirmation, err := getTemplate("email_confirm", connector)
	if err != nil {
		return errors.Wrap(err, "failed to get template")
	}
	c.conf.EmailConfirmation = emailConfirmation

	kycApprove, err := getTemplate("kyc_approve", connector)
	if err != nil {
		return errors.Wrap(err, "failed to get template")
	}

	c.conf.KYCApprove = kycApprove

	kycReject, err := getTemplate("kyc_reject", connector)
	if err != nil {
		return errors.Wrap(err, "failed to get template")
	}

	c.conf.KYCReject = kycReject

	return nil
}

func getTemplate(name string, connector *horizon.Connector) (*template.Template, error) {
	q := connector.Templates()

	body, err := q.Get(name)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("unable to download %s", name))
	}

	return template.New(name).Parse(string(body))
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
		Header: "Swarm Verification Request",
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
		Header: "Swarm Verification Request",
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
