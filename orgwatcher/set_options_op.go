package orgwatcher

type SetOptionsOp struct {
	ID            string `json:"id"`
	SourceAccount string `json:"source_account"`
	Type          string `json:"type"`
	SignerKey     string `json:"signer_key"`
	SignerWeight  int32  `json:"signer_weight"`
}
