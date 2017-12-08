package storage

import (
	"github.com/pkg/errors"
)

var (
	allowedMediaTypes = map[string]string{
		"application/pdf": "pdf",
		"image/jpeg":      "jpeg",
		"image/tiff":      "tiff",
		"image/png":       "png",
		"image/gif":       "gif",
	}

	IsAllowedContentType = func(value interface{}) error {
		str, ok := value.(string)
		if !ok {
			return errors.New("string expected")
		}
		if IsContentTypeAllowed(str) {
			return nil
		}
		return errors.New("not allowed")
	}
)

func IsContentTypeAllowed(mediaType string) bool {
	_, ok := allowedMediaTypes[mediaType]
	return ok
}

// ContentTypeExtension return file extension for corresponding media type.
// Will panic if media type is unknown, use `IsContentTypeAllowed` first
func ContentTypeExtension(mediaType string) string {
	ext, ok := allowedMediaTypes[mediaType]
	if !ok {
		panic("unknown media type")
	}
	return ext
}
