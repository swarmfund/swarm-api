package storage

var allowedMediaTypes = map[string]string{
	"application/pdf": "pdf",
	"image/jpeg":      "jpeg",
	"image/tiff":      "tiff",
	"image/png":       "png",
	"image/gif":       "gif",
}

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
