package storage

import (
	"fmt"

	"regexp"
	"strconv"

	"net/url"

	"gitlab.com/swarmfund/api/db2/api"
	"github.com/pkg/errors"
)

type Document struct {
	AccountID string
	Type      api.DocumentType
	EntityID  int64
	Version   string
	Extension string
}

func (d *Document) Key() string {
	return fmt.Sprintf("%d-%d-%s.%s",
		d.Type, d.EntityID, d.Version, d.Extension)
}

// FromKey parse object key and returns filled `Document` structure.
// Note: `AcccountID` won't be filled
func FromKey(key string) (*Document, error) {
	r := regexp.MustCompile(`^(?P<type>\d+)-(?P<entity>\d+)-(?P<version>[^.]+)\.(?P<ext>\w+)$`)
	key, err := url.QueryUnescape(key)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unescape key")
	}
	match := r.FindStringSubmatch(key)
	if len(match) != 5 {
		return nil, errors.New("unknown key format")
	}

	rawType, err := strconv.ParseInt(match[1], 10, 64)
	if err != nil {
		return nil, err
	}

	entity, err := strconv.ParseInt(match[2], 10, 64)
	if err != nil {
		return nil, err
	}

	document := Document{
		Type:      api.DocumentType(rawType),
		EntityID:  entity,
		Version:   match[3],
		Extension: match[4],
	}
	return &document, nil
}
