package clienturl

import "net/http"

func DiscourseSSO(signature, returnURL, nonce string) Payload {
	return Payload{
		Status: http.StatusOK,
		Type:   RedirectTypeDiscourseSSO,
		Meta: map[string]interface{}{
			"signature": signature,
			"return":    returnURL,
			"nonce":     nonce,
		},
	}
}
