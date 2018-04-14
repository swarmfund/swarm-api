package storage

import (
	"github.com/pkg/errors"
)

var (
	allowedMediaTypes = map[string]struct{}{}

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

func SetMediaTypes(mediaTypes []string) {
	for _, t := range mediaTypes {
		allowedMediaTypes[t] = struct{}{}
	}
}
