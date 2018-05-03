package storage

import (
	"fmt"

	"github.com/pkg/errors"
	"gitlab.com/swarmfund/api/internal/types"
)

var (
	allowedMediaTypes = map[types.DocumentType]map[string]struct{}{}

	IsAllowedContentType = func(docType types.DocumentType) func(value interface{}) error {
		return func(value interface{}) error {
			str, ok := value.(string)
			if !ok {
				return errors.New("string expected")
			}
			if IsContentTypeAllowed(docType, str) {
				return nil
			}
			return errors.New("not allowed")
		}
	}
)

func IsContentTypeAllowed(docType types.DocumentType, mediaType string) bool {
	_, ok := allowedMediaTypes[0][mediaType]
	if ok {
		return true
	}
	specific, ok := allowedMediaTypes[docType]
	if ok {
		_, ok := specific[mediaType]
		if ok {
			return true
		}
	}
	return false
}

func SetMediaTypes(mediaTypes map[string][]string) error {
	for docName, extensions := range mediaTypes {
		var docType types.DocumentType
		//init current doc type
		if err := docType.UnmarshalJSON([]byte(fmt.Sprintf(`"%s"`, docName))); err != nil {
			return errors.Wrap(err, "failed to get document type name")
		}

		if extensions != nil {
			for _, t := range extensions {
				allowedMediaTypes[types.DocumentTypeAlpha] = map[string]struct{}{
					t: {},
				}
			}
		}
	}

	return nil
}
