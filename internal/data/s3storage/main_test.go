package s3storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/swarmfund/api/internal/storage"

	"gitlab.com/swarmfund/api/internal/types"
)

func TestViperConfig_StorageMediaTypesHook(t *testing.T) {
	localStorage := Storage{}
	var err error
	localStorage.mediaTypes, err = storage.NewMediaTypes(
		map[string][]string{
			"general": {"image/jpeg", "image/jpeg", "image/png", "image/gif"},
			"alpha":   {"application/pdf"},
		})
	assert.NoError(t, err)

	assert.True(t, localStorage.IsContentTypeAllowed(types.DocumentTypeAlpha, "application/pdf"))
	assert.True(t, localStorage.IsContentTypeAllowed(types.DocumentTypeGeneral, "image/jpeg"))
	assert.True(t, localStorage.IsContentTypeAllowed(types.DocumentTypeGeneral, "image/jpeg"))
	assert.True(t, localStorage.IsContentTypeAllowed(types.DocumentTypeGeneral, "image/png"))
	assert.True(t, localStorage.IsContentTypeAllowed(types.DocumentTypeGeneral, "image/gif"))
	assert.True(t, localStorage.IsContentTypeAllowed(types.DocumentTypeAlpha, "image/gif"))
	assert.False(t, localStorage.IsContentTypeAllowed(types.DocumentTypeGeneral, "application/pdf"))
}
