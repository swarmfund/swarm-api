package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/swarmfund/api/internal/types"
)

func TestMediaType(t *testing.T) {
	cases := []struct {
		name      string
		config    map[string][]string
		DocType   types.DocumentType
		MediaType string
		expected  bool
	}{
		{
			name:      "empty",
			config:    nil,
			DocType:   types.DocumentTypeAlpha,
			MediaType: "application/pdf",
			expected:  false,
		},
		{
			name: "allowed for all documents",
			config: map[string][]string{
				"general": {"image/jpeg"},
			},
			DocType:   types.DocumentTypeAlpha,
			MediaType: "image/jpeg",
			expected:  true,
		},
		{
			name: "not allowed in general",
			config: map[string][]string{
				"general": {"image/jpeg"},
			},
			DocType:   types.DocumentTypeAlpha,
			MediaType: "image/png",
			expected:  false,
		},
		{
			name: "allowed only for specific",
			config: map[string][]string{
				"alpha": {"application/pdf"},
			},
			DocType:   types.DocumentTypeAlpha,
			MediaType: "application/pdf",
			expected:  true,
		},
		{
			name: "not allowed only for specific",
			config: map[string][]string{
				"alpha": {"application/pdf"},
			},
			DocType:   types.DocumentTypeAlpha,
			MediaType: "application/xml",
			expected:  false,
		},
		{
			name: "allowed in general and in specific",
			config: map[string][]string{
				"general": {"image/jpeg", "image/tiff", "image/png", "image/gif"},
				"alpha":   {"image/jpeg"},
			},
			DocType:   types.DocumentTypeAlpha,
			MediaType: "image/jpeg",
			expected:  true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := NewMediaTypes(c.config)
			assert.NoError(t, err)
			assert.Equal(t, c.expected, got.IsAllowed(c.DocType, c.MediaType))
		})
	}
}
