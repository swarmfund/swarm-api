package notificator

import (
	"fmt"
	"html/template"
	"net/url"

	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/notificator"
)

const (
	NotificationTypeVerificationEmail = 3
)

type Config struct {
	Disabled     bool     `fig:"disabled"`
	Endpoint     *url.URL `fig:"endpoint"`
	Secret       string   `fig:"secret"`
	Public       string   `fig:"public"`
	ClientRouter string   `fig:"client_router"`

	EmailConfirmation *template.Template `fig:"-"`
	KYCApprove        *template.Template `fig:"-"`
	KYCReject         *template.Template `fig:"-"`
	WelcomeEmail      *template.Template `fig:"-"`
}

type Connector struct {
	notificator *notificator.Connector
	log         *logan.Entry
	conf        Config
}

func NewConnector(conf Config) *Connector {
	return &Connector{
		notificator: notificator.NewConnector(
			notificator.Pair{Secret: conf.Secret, Public: conf.Public},
			*conf.Endpoint,
		),
		conf: conf,
	}
}

func (c *Connector) Init(loader TemplateLoader, log *logan.Entry) error {
	c.log = log

	emailConfirmation, err := getTemplate("email_confirm", loader)
	if err != nil {
		return errors.Wrap(err, "failed to get template")
	}
	c.conf.EmailConfirmation = emailConfirmation

	kycApprove, err := getTemplate("kyc_approve", loader)
	if err != nil {
		return errors.Wrap(err, "failed to get template")
	}

	c.conf.KYCApprove = kycApprove

	kycReject, err := getTemplate("kyc_reject", loader)
	if err != nil {
		return errors.Wrap(err, "failed to get template")
	}

	c.conf.KYCReject = kycReject

	welcomeEmail, err := getTemplate("welcome_email", loader)
	if err != nil {
		return errors.Wrap(err, "failed to get template")
	}
	c.conf.WelcomeEmail = welcomeEmail

	return nil
}

type TemplateLoader interface {
	Get(id string) ([]byte, error)
}

func getTemplate(name string, loader TemplateLoader) (*template.Template, error) {
	body, err := loader.Get(name)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("unable to download %s", name))
	}

	return template.New(name).Parse(string(body))
}
