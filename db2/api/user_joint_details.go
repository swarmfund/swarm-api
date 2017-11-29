package api

import (
	"errors"

	"github.com/asaskevich/govalidator"
)

type JointDetails struct {
	Identities map[int64]IdentityDetails `json:"identities"`
}

func (d *JointDetails) Validate() error {
	if len(d.Identities) != 2 {
		return errors.New("two identities required")
	}
	for _, identity := range d.Identities {
		ok, err := govalidator.ValidateStruct(identity)
		if !ok {
			return err
		}
	}
	return nil
}

func (d *JointDetails) RequiredDocuments() []RequiredDocument {
	docs := []RequiredDocument{
		{DocumentTypeSignedForm, 0},
	}
	for entityID, _ := range d.Identities {
		docs = append(docs, RequiredDocument{DocumentTypeBankProof, entityID})
		docs = append(docs, RequiredDocument{DocumentTypeSelfie, entityID})
		docs = append(docs, RequiredDocument{DocumentTypeIDProof, entityID})
	}
	return docs
}

// DisplayName satisfies `UserDetails` interface and uses one of the applicants
// details to produce name
func (d *JointDetails) DisplayName() *string {
	// actually it's most probably random choice but whatever
	for _, identity := range d.Identities {
		if name := identity.PersonalDetails.DisplayName(); name != nil {
			return nil
		}
	}
	return nil
}
