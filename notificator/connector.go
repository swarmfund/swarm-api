package notificator

import (
	"html/template"
	"net/url"

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
	loader      TemplateLoader
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
	c.loader = loader

	return nil
}

type TemplateLoader interface {
	Get(id string) ([]byte, error)
}
