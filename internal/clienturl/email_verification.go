package clienturl

import "net/http"

func EmailVerification(walletID, token string) Payload {
	return Payload{
		Status: http.StatusOK,
		Type:   RedirectTypeEmailVerification,
		Meta: map[string]interface{}{
			"wallet_id": walletID,
			"token":     token,
		},
	}
}
