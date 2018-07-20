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
	general := m.allowed[types.DocumentTypeGeneral]
	specific := m.allowed[docType]
	allowed := append(general, specific...)
	for _, val := range allowed {
		if val == mediaType {
			return true
		}
	}

	return false
}
