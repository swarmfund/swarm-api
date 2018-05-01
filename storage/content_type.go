package storage

import (
	"github.com/pkg/errors"
	"gitlab.com/swarmfund/api/internal/types"
)

var (
	allowedMediaTypes = map[types.DocumentType]map[string]struct{}{
		0: map[string]struct{}{
			"image/jpeg": struct{}{},
			"image/tiff": struct{}{},
			"image/png":  struct{}{},
			"image/gif":  struct{}{},
		},
		types.DocumentTypeFundDocument: map[string]struct{}{
			"application/pdf": struct{}{},
		},
		types.DocumentTypeNavReport: map[string]struct{}{
			"application/pdf": struct{}{},
		},
		types.DocumentTypeAlpha: map[string]struct{}{
			"application/pdf": struct{}{},
		},
		types.DocumentTypeBravo: map[string]struct{}{
			"application/pdf": struct{}{},
		},
		types.DocumentTypeCharlie: map[string]struct{}{
			"application/pdf": struct{}{},
		},
		types.DocumentTypeDelta: map[string]struct{}{
			"application/pdf": struct{}{},
		},
		types.DocumentTypeTokenTerms: map[string]struct{}{
			"application/pdf": struct{}{},
		},
		types.DocumentTypeTokenMetrics: map[string]struct{}{
			"application/pdf": struct{}{},
		},
	}

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
