package api

import (
	"encoding/json"
	"fmt"
	"strconv"
)

const (
	_ DocumentType = iota
	DocumentTypeCertification
	DocumentTypeAddressProof
	DocumentTypeBankProof
	DocumentTypeSelfie
	DocumentTypeIDProof // 5
	DocumentTypeRegistrationCertificate
	DocumentTypeIncorporationCertificate
	DocumentTypeMOA
	DocumentTypeAOA
	DocumentTypePOA // 10
	DocumentTypePhoneProof
	DocumentTypeAudit
	DocumentTypeSignedForm
	_
	DocumentTypeProofOfIncome // 15
	//
	DocumentTypeAssetLogo
)

var (
	kycDocs = map[DocumentType]bool{
		DocumentTypeCertification:            true,
		DocumentTypeAddressProof:             true,
		DocumentTypeBankProof:                true,
		DocumentTypeSelfie:                   true,
		DocumentTypeIDProof:                  true,
		DocumentTypeRegistrationCertificate:  true,
		DocumentTypeIncorporationCertificate: true,
		DocumentTypeMOA:                      true,
		DocumentTypeAOA:                      true,
		DocumentTypePOA:                      true,
		DocumentTypePhoneProof:               true,
		DocumentTypeAudit:                    true,
		DocumentTypeSignedForm:               true,
	}

	proofOfIncomeDocs = map[DocumentType]bool{
		DocumentTypeProofOfIncome: true,
	}
)

type DocumentType int

func (t *DocumentType) UnmarshalJSON(bytes []byte) error {
	// one might think unmarshal interface{} and then type switch is a better way,
	// that's until you discover it will produce float64
	var i int64
	if err := json.Unmarshal(bytes, &i); err == nil {
		*t = DocumentType(i)
		return nil
	}

	var s string
	if err := json.Unmarshal(bytes, &s); err == nil {
		i, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return err
		}
		*t = DocumentType(i)
		return nil
	}
	return fmt.Errorf("can't unmarshal DocumentType")
}

// IsKYC return true if document is used for KYC flow
func (t DocumentType) IsKYC() bool {
	_, ok := kycDocs[t]
	return ok
}

// IsProofOfIncome return true if document is used for POI review flow
func (t DocumentType) IsProofOfIncome() bool {
	_, ok := proofOfIncomeDocs[t]
	return ok
}
