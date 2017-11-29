package api

var (
	DocumentVersionLatest = "latest"
)

type Document struct {
	Key         string                 `json:"key"`
	Type        DocumentType           `json:"type"`
	EntityID    int64                  `json:"entity_id"`
	ContentType string                 `json:"content_type"`
	Checksum    string                 `json:"checksum"`
	CreatedAt   int64                  `json:"created_at"`
	Version     string                 `json:"version"`
	Meta        map[string]interface{} `json:"meta,omitempty"`
}
