package blacklist

import (
	"strings"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

var (
	ErrUnallowedValue = errors.New("unallowed value")
	ErrUnallowedType  = errors.New("bad type of value")
)

type Approver struct {
	blacklist []string
}

func NewApprover(unallowed ...string) *Approver {
	unallowed = append(unallowed, "") //default value of type string is also unallowed for domain
	return &Approver{
		blacklist: unallowed,
	}
}

func getDomain(email string) (domain string) {
	splitted := strings.Split(email, "@")
	if len(splitted) > 1 {
		domain = splitted[1]
	}

	return domain
}

func (a Approver) Validate(value interface{}) error {
	str, ok := value.(string)
	if !ok {
		return errors.Wrap(ErrUnallowedType, "value should be of type string")
	}

	domain := getDomain(str)
	for _, unallowed := range a.blacklist {
		if domain == unallowed {
			return errors.Wrap(ErrUnallowedValue, "failed to approve email, domain in blacklist ", logan.F{"domain": unallowed})
		}
	}

	return nil
}
