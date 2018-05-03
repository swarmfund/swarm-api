package storage

type MediaTypes struct {
	types []string
}

func NewMediaTypes(types []string) MediaTypes {
	return MediaTypes{types: types}
}

func (m *MediaTypes) IsAllowed(mediaType string) bool {
	for _, t := range m.types {
		if t == mediaType {
			return true
		}
	}
	return false
}
