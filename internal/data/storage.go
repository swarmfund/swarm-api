package data

import (
	"net/url"

	"gitlab.com/swarmfund/api/internal/types"
)

type Storage interface {
	SignedObjectURL(key string) (*url.URL, error)
	// TODO move out of here
	IsContentTypeAllowed(docType types.DocumentType, mediaType string) bool
	UploadFormData(key string) (map[string]string, error)
}
