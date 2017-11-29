package clienturl

import "testing"

func TestEmailVerification(t *testing.T) {
	walletID := "wallet-id"
	token := "token"
	payload := EmailVerification(walletID, token)
	if payload.Type != RedirectTypeEmailVerification {
		t.Errorf("expected %d got %d", RedirectTypeEmailVerification, payload.Type)
	}
	got := payload.Meta["token"]
	if got != token {
		t.Errorf("expected %s got %s", token, got)
	}
	got = payload.Meta["wallet_id"]
	if got != walletID {
		t.Errorf("expected %s got %s", got, walletID)
	}
}
