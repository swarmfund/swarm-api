package api

import (
	"encoding/json"

	"github.com/pkg/errors"
)

// HeteroUserDetails provides better way to sanitize user details and store
// results when user type is not know at compile time
type HeteroUserDetails struct {
	UserType UserType
	Entities KYCEntities
	raw      []byte
}

func (r *HeteroUserDetails) UnmarshalJSON(data []byte) error {
	type t struct {
		UserType `json:"user_type"`
		Details  json.RawMessage
	}
	var tt t
	if err := json.Unmarshal(data, &tt); err != nil {
		return err
	}
	// support for setting user type from above
	if r.UserType == UserTypeNil {
		r.UserType = tt.UserType
	}
	r.raw = data
	switch r.UserType {
	case UserTypeIndividual:
		details := IndividualDetails{}
		if err := json.Unmarshal(tt.Details, &details); err != nil {
			return errors.Wrap(err, "failed to unmarshal individual details")
		}
		for _, g := range []KYCEntityGetter{
			details.PersonalDetails, details.Address,
			details.BankDetails, details.EmploymentDetails,
		} {
			entity, err := g.KYCEntity()
			if err != nil {
				return errors.Wrap(err, "failed to get entity")
			}
			r.Entities = append(r.Entities, entity)
		}
	case UserTypeJoint:
		details := JointDetails{}
		if err := json.Unmarshal(tt.Details, &details); err != nil {
			return errors.Wrap(err, "failed to unmarshal joint details")
		}
		for _, identity := range details.Identities {
			entity, err := identity.KYCEntity()
			if err != nil {
				return errors.Wrap(err, "failed to get entity")
			}
			r.Entities = append(r.Entities, entity)
		}
	case UserTypeBusiness:
		details := BusinessDetails{}
		if err := json.Unmarshal(tt.Details, &details); err != nil {
			return errors.Wrap(err, "failed to unmarshal business details")
		}
		// corr and registered address have different entity type
		entity, err := details.CorrespondenceAddress.KYCEntity()
		if err != nil {
			return errors.Wrap(err, "failed to get entity")
		}
		entity.Type = KYCEntityTypeCorrespondenceAddress
		r.Entities = append(r.Entities, entity)

		entity, err = details.RegisteredAddress.KYCEntity()
		if err != nil {
			return errors.Wrap(err, "failed to get entity")
		}
		entity.Type = KYCEntityTypeRegisteredAddress
		r.Entities = append(r.Entities, entity)

		for _, g := range []KYCEntityGetter{
			details.CorporationDetails, details.FinancialDetails,
		} {
			entity, err := g.KYCEntity()
			if err != nil {
				return errors.Wrap(err, "failed to get entity")
			}
			r.Entities = append(r.Entities, entity)
		}

		for _, owner := range details.Owners {
			entity, err := owner.KYCEntity()
			if err != nil {
				return errors.Wrap(err, "failed to get entity")
			}
			entity.Type = KYCEntityTypeBusinessOwner
			r.Entities = append(r.Entities, entity)
		}

		for _, signatory := range details.Signatories {
			entity, err := signatory.KYCEntity()
			if err != nil {
				return errors.Wrap(err, "failed to get entity")
			}
			entity.Type = KYCEntityTypeBusinessSignatory
			r.Entities = append(r.Entities, entity)
		}
	default:
		return ErrUnknownUserType
	}
	return nil
}
