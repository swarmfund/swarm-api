package api

// RejectReasons interface describing minimal behaviour required from user reject
// reasons implementations
type RejectReasons interface {
	// Empty returns true if reject reasons are clean and user could be submitted
	// for verification again
	Empty() bool
}

// ValidateRejectReasons validate and sanitize reject reasons JSON according
// to user type
// TODO move to reject reasons interface
func (user *User) ValidateRejectReasons(raw []byte) ([]byte, error) {
	switch user.UserType {
	//case UserTypeIndividual:
	//	rr := IndividualRejectReasons{}
	//	if err := json.Unmarshal(raw, &rr); err != nil {
	//		return nil, errors.Wrap(err, "failed to unmarshal")
	//	}
	//	raw, err := json.Marshal(&rr)
	//	if err != nil {
	//		return nil, errors.Wrap(err, "failed to marshal")
	//	}
	//	return raw, nil
	//case UserTypeJoint:
	//	rr := JointRejectReasons{}
	//	if err := json.Unmarshal(raw, &rr); err != nil {
	//		return nil, errors.Wrap(err, "failed to unmarshal")
	//	}
	//	raw, err := json.Marshal(&rr)
	//	if err != nil {
	//		return nil, errors.Wrap(err, "failed to marshal")
	//	}
	//	return raw, nil
	//case UserTypeBusiness:
	//	rr := BusinessRejectReasons{}
	//	if err := json.Unmarshal(raw, &rr); err != nil {
	//		return nil, errors.Wrap(err, "failed to unmarshal")
	//	}
	//	raw, err := json.Marshal(&rr)
	//	if err != nil {
	//		return nil, errors.Wrap(err, "failed to marshal")
	//	}
	//	return raw, nil
	default:
		panic("user type is unknown")
	}
}
