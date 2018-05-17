package storage

import (
	"fmt"

	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/api/internal/types"
)

type MediaTypes struct {
	allowed map[types.DocumentType][]string
}

func NewMediaTypes(mediaTypes map[string][]string) (MediaTypes, error) {
	allowedMediaTypes := map[types.DocumentType][]string{}
	for docName, mediaType := range mediaTypes {
		var docType types.DocumentType

		if err := docType.UnmarshalJSON([]byte(fmt.Sprintf(`"%s"`, docName))); err != nil {
			return MediaTypes{}, errors.Wrap(err, "failed to get document type")
		}

		allowedMediaTypes[docType] = mediaType
	}

	return MediaTypes{allowed: allowedMediaTypes}, nil
}

func (m *MediaTypes) IsAllowed(docType types.DocumentType, mediaType string) bool {
	general, ok := m.allowed[types.DocumentTypeGeneral]
	if ok {
		for _, val := range general {
			if val == mediaType {
				return true
			}
		}
	}

	other, ok := m.allowed[docType]
	if ok {
		for _, val := range other {
			if val == mediaType {
				return true
			}
		}
	}
	return false
}
