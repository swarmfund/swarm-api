package notificator

import (
	"fmt"
	"html/template"

	"bytes"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/distributed_lab/notificator"
	"gitlab.com/swarmfund/api/internal/clienturl"
)

var (
	ErrTemplateNotFound = errors.New("Template not found")
)

func tryGetTemplate(name string, loader TemplateLoader) (*template.Template, error) {
	body, err := loader.Get(name)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("unable to download %s", name))
	}
	if body == nil {
		return nil, ErrTemplateNotFound
	}
	return template.New(name).Parse(string(body))
}

func (c *Connector) lazyLoadTemplate(name string, template *template.Template) (*template.Template, error) {
	if template == nil {
		loadedTemplate, err := tryGetTemplate(name, c.loader)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get", logan.F{"template": name})
		}
		return loadedTemplate, nil
	}
	return template, nil
}

func (c *Connector) SendVerificationLink(email string, payload clienturl.Payload) (err error) {
	encoded, err := payload.Encode()
	if err != nil {
		return errors.Wrap(err, "failed to encode payload")
	}
	letter := &Letter{
		Header: "Email Verification",
		Link:   fmt.Sprintf("%s/%s", c.conf.ClientRouter, encoded),
	}

	c.conf.EmailConfirmation, err = c.lazyLoadTemplate("email_confirm", c.conf.EmailConfirmation)
	if err != nil {
		return errors.Wrap(err, "failed to load template")
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

func (c *Connector) NotifyApproval(email string) (err error) {
	letter := &Letter{
		Header: "Verification Request",
		Link:   c.conf.ClientRouter,
	}

	c.conf.KYCApprove, err = c.lazyLoadTemplate("kyc_approve", c.conf.KYCApprove)
	if err != nil {
		return errors.Wrap(err, "failed to load template")
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

func (c *Connector) NotifyRejection(email string) (err error) {
	letter := &Letter{
		Header: "Verification Request",
		Link:   c.conf.ClientRouter,
	}

	c.conf.KYCReject, err = c.lazyLoadTemplate("kyc_reject", c.conf.KYCReject)
	if err != nil {
		return errors.Wrap(err, "failed to load template")
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

func (c *Connector) SendWelcomeEmail(email string) (err error) {
	letter := &Letter{
		Header: "Welcome!",
		Link:   c.conf.ClientRouter,
	}
	c.conf.WelcomeEmail, err = c.lazyLoadTemplate("welcome_email", c.conf.WelcomeEmail)
	if err != nil {
		return err
	}

	var buff bytes.Buffer
	err = c.conf.WelcomeEmail.Execute(&buff, letter)
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
