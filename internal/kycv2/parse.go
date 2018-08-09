package kyc

import (
	"encoding/json"

	"github.com/pkg/errors"
)

// ParsingData describes the structure of KYC Blob retrieved form Horizon.
type parsingData struct {
	FirstName  string      `json:"first_name"`
	LastName   string      `json:"last_name"`
	Address    Address     `json:"address"`
	ETHAddress string      `json:"eth_address"`
	Documents  DocumentsV1 `json:"documents"`

	Version string `json:"version"`
	V2      Data   `json:"v2"`
}

type DocumentsV1 struct {
	KYCIdDocument     IDDocument     `json:"kyc_id_document"`
	KYCProofOfAddress ProofOfAddrDoc `json:"kyc_poa"`
}

func ParseKYCData(data string) (*Data, error) {
	var parsingKYCData parsingData
	err := json.Unmarshal([]byte(data), &parsingKYCData)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to unmarshal data bytes into Data structure")
	}

	switch parsingKYCData.Version {
	case "v2":
		return &parsingKYCData.V2, nil
	default:
		// v1
		return &Data{
			FirstName: parsingKYCData.FirstName,
			LastName:  parsingKYCData.LastName,
			Address:   parsingKYCData.Address,
			Documents: Documents{
				IDDocument: IDDocument{
					FaceFile: DocFile{
						ID:   parsingKYCData.Documents.KYCIdDocument.FaceFile.ID,
						Name: parsingKYCData.Documents.KYCIdDocument.FaceFile.Name,
					},
					BackFile: nil,
					// v1 only supported passports
					Type: PassportDocType,
				},
				ProofOfAddr: ProofOfAddrDoc{
					FaceFile: DocFile{
						ID: parsingKYCData.Documents.KYCProofOfAddress.FaceFile.ID,
					},
				},
			},
		}, nil
	}
}

// ParsingFirstNameData describes the structure of shortened(FirstName only) KYC Blob retrieved form Horizon.
type parsingFirstNameData struct {
	FirstName string `json:"first_name"`

	Version string        `json:"version"`
	V2      FirstNameData `json:"v2"`
}

func ParseKYCFirstName(data string) (string, error) {
	var parsingData parsingFirstNameData
	err := json.Unmarshal([]byte(data), &parsingData)
	if err != nil {
		return "", errors.Wrap(err, "Failed to unmarshal data bytes into Data structure")
	}

	switch parsingData.Version {
	case "v2":
		return parsingData.V2.FirstName, nil
	default:
		// v1
		return parsingData.FirstName, nil
	}
}

type AutoGenerated struct {
	Address struct {
		Line1      string `json:"line_1"`
		Line2      string `json:"line_2"`
		City       string `json:"city"`
		Country    string `json:"country"`
		State      string `json:"state"`
		PostalCode string `json:"postal_code"`
	} `json:"address"`
	AltAddress struct {
		AltLine1      string `json:"alt_line_1"`
		AltLine2      string `json:"alt_line_2"`
		AltCity       string `json:"alt_city"`
		AltCountry    string `json:"alt_country"`
		AltState      string `json:"alt_state"`
		AltPostalCode string `json:"alt_postal_code"`
	} `json:"alt_address"`
	PrivacyPolicy bool   `json:"privacy_policy"`
	CompanyName   string `json:"company_name"`
	PhoneNumber   string `json:"phone_number"`
	Website       string `json:"website"`
	Owners        []struct {
		FirstName  string `json:"first_name"`
		LastName   string `json:"last_name"`
		Percentage string `json:"percentage"`
	} `json:"owners"`
	Managers []struct {
		FirstName   string `json:"first_name"`
		LastName    string `json:"last_name"`
		Position    string `json:"position"`
		Years       string `json:"years"`
		DirectPhone string `json:"direct_phone"`
		Email       string `json:"email"`
	} `json:"managers"`
	CompanyType             string `json:"company_type"`
	Teaser                  string `json:"teaser"`
	OwnershipStructure      string `json:"ownership_structure"`
	OrganizationalStructure string `json:"organizational_structure"`
	Documents               struct {
		CompanyOriginationCertificate struct {
			Front struct {
				MimeType string `json:"mime_type"`
				Name     string `json:"name"`
				Key      string `json:"key"`
			} `json:"front"`
		} `json:"company_origination_certificate"`
		KycTaxReturns struct {
			Front struct {
				MimeType string `json:"mime_type"`
				Name     string `json:"name"`
				Key      string `json:"key"`
			} `json:"front"`
		} `json:"kyc_tax_returns"`
		KycInvestmentPresentation struct {
			Front struct {
				MimeType string `json:"mime_type"`
				Name     string `json:"name"`
				Key      string `json:"key"`
			} `json:"front"`
		} `json:"kyc_investment_presentation"`
		KycManagementBios struct {
			Front struct {
				MimeType string `json:"mime_type"`
				Name     string `json:"name"`
				Key      string `json:"key"`
			} `json:"front"`
		} `json:"kyc_management_bios"`
		KycManagementCvs struct {
			Front struct {
				MimeType string `json:"mime_type"`
				Name     string `json:"name"`
				Key      string `json:"key"`
			} `json:"front"`
		} `json:"kyc_management_cvs"`
		KycManagementPictures struct {
			Front struct {
				MimeType string `json:"mime_type"`
				Name     string `json:"name"`
				Key      string `json:"key"`
			} `json:"front"`
		} `json:"kyc_management_pictures"`
	} `json:"documents"`
}
