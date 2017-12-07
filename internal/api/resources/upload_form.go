package resources

import (
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
)

type UploadForm struct {
	Type       string               `json:"type"`
	Attributes UploadFormAttributes `json:"attributes"`
}

type UploadFormAttributes struct {
	Bucket         string `json:"bucket"`
	Key            string `json:"key"`
	Policy         string `json:"policy"`
	URL            string `json:"url"`
	XAMZAlgorithm  string `json:"x-amz-algorithm" mapstructure:"x-amz-algorithm"`
	XAMZCredential string `json:"x-amz-credential" mapstructure:"x-amz-credential"`
	XAMZDate       string `json:"x-amz-date" mapstructure:"x-amz-date"`
	XAMZSignature  string `json:"x-amz-signature" mapstructure:"x-amz-signature"`
}

func NewUploadForm(form map[string]string) UploadForm {
	r := UploadForm{
		Type: "upload_policy",
	}
	if err := mapstructure.Decode(&form, &r.Attributes); err != nil {
		panic(errors.Wrap(err, "failed to decode upload policy"))
	}
	return r
}
