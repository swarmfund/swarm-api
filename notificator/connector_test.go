package notificator

import (
	"testing"

	"gitlab.com/swarmfund/api/assets"
	"gitlab.com/swarmfund/api/config"
	"gitlab.com/swarmfund/api/internal/clienturl"
)

var (
	conf = config.Notificator{
		Endpoint:          "http://18.195.18.3:9009",
		Secret:            "",
		Public:            "",
		ClientRouter:      "http://client/hredir",
		EmailConfirmation: assets.Templates.Lookup("email_confirm"),
	}

	conn = NewConnector(conf)
)

func TestSendVerificationLink(t *testing.T) {
	payload := clienturl.Payload{
		Status: 200,
		Type:   clienturl.RedirectTypeEmailVerification,
	}
	err := conn.SendVerificationLink("gilirivi@minex-coin.com", payload)
	if err != nil {
		t.Error(err)
	}

}
