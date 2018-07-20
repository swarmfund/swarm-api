package notificator

import (
	"fmt"

	"bytes"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/distributed_lab/notificator"
	"gitlab.com/swarmfund/api/internal/clienturl"
)

func (c *Connector) SendVerificationLink(email string, payload clienturl.Payload) error {
	encoded, err := payload.Encode()
	if err != nil {
		return errors.Wrap(err, "failed to encode payload")
	}
	letter := &Letter{
		Header: "Email Verification",
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
		Header: "Verification Request",
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
		Header: "Verification Request",
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

func (c *Connector) SendWelcomeEmail(email string) error {
	letter := &Letter{
		Header: "Welcome!",
		Link:   c.conf.ClientRouter,
	}

	var buff bytes.Buffer
	err := c.conf.WelcomeEmail.Execute(&buff, letter)
	if err != nil {
		return errors.Wrap(err, "failed to render template for welcome email")
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
		c.log.WithFields(logan.F{"request_type": requestType, "token": token}).Warn("disabled")
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
