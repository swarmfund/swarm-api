package horizon

type Signers struct {
	Signers []Signer `json:"signers"`
}

type SignerType struct {
	Name  string `json:"name"`
	Value int32  `json:"value"`
}

// Signer represents one of an account's signers.
type Signer struct {
	PublicKey      string       `json:"public_key"`
	Weight         int32        `json:"weight"`
	SignerTypeI    int32        `json:"signer_type_i"`
	SignerTypes    []SignerType `json:"signer_types"`
	SignerIdentity int32        `json:"signer_identity"`
	SignerName     string       `json:"signer_name"`
}
