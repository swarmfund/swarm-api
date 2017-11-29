package clienturl

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
)

type Payload struct {
	// Status action result status follows HTTP status code semantics
	Status int `json:"status,omitempty"`
	// Type specific redirect type which meta corresponds to
	Type RedirectType `json:"type"`
	// Meta holds additional information need for redirect processing
	Meta map[string]interface{} `json:"meta,omitempty"`
}

// Encode with base64url without padding
func (p *Payload) Encode() (string, error) {
	// default status should be omitted
	if p.Status == http.StatusOK {
		p.Status = 0
	}

	bytes, err := json.Marshal(p)
	if err != nil {
		return "", err
	}
	encoded := base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(bytes)
	return encoded, nil
}
