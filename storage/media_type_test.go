package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/swarmfund/api/internal/types"
)

func TestMediaType(t *testing.T) {
	mediaTypes, err := NewMediaTypes(
		map[string][]string{
			"general": {"image/jpeg", "image/tiff", "image/png", "image/gif"},
			"alpha":   {"application/pdf"},
		})
	assert.NoError(t, err)

	cases := []struct {
		name      string
		DocType   types.DocumentType
		MediaType string
		expected  bool
	}{
		{
			name:      "allowed for all documents",
			DocType:   types.DocumentTypeDelta,
			MediaType: "image/jpeg",
			expected:  true,
		},
		{
			name:      "not allowed",
			DocType:   types.DocumentTypeAlpha,
			MediaType: "audio/mp4",
			expected:  false,
		},
		{
			name:      "allowed only for document type alpha",
			DocType:   types.DocumentTypeAlpha,
			MediaType: "application/pdf",
			expected:  true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := mediaTypes.IsAllowed(c.DocType, c.MediaType)
			assert.Equal(t, c.expected, got)
		})
	}
}
