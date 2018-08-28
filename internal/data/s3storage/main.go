package s3storage

import (
	"net/url"
	"time"

	"strings"

	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/minio/minio-go"
	"github.com/pkg/errors"
	"gitlab.com/swarmfund/api/internal/data"
	"gitlab.com/swarmfund/api/internal/storage"
	"gitlab.com/swarmfund/api/internal/types"
)

const (
	presignExpire = 1 * time.Hour
)

func NewStorage(session *session.Session, bucket string, mediaTypes storage.MediaTypes) (data.Storage, error) {
	creds, err := session.Config.Credentials.Get()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get credentials")
	}

	endpoint := fmt.Sprintf("s3.%s.amazonaws.com", *session.Config.Region)
	if session.Config.Endpoint != nil {
		endpoint = *session.Config.Endpoint
	}
	mc, err := minio.New(
		endpoint,
		creds.AccessKeyID,
		creds.SecretAccessKey,
		true,
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to init minio client")
	}

	return &Storage{
		s3.New(session),
		mc,
		bucket,
		mediaTypes,
	}, nil
}

type Storage struct {
	s3 *s3.S3
	// keeping minio for backwards compatibility
	minio      *minio.Client
	bucket     string
	mediaTypes storage.MediaTypes
}

func (s *Storage) SignedObjectURL(key string) (*url.URL, error) {
	request, _ := s.s3.GetObjectRequest(&s3.GetObjectInput{
		Bucket:              aws.String(s.bucket),
		Key:                 aws.String(key),
		ResponseContentType: aws.String(""),
	})

	signed, err := request.Presign(presignExpire)
	if err != nil {
		return nil, errors.Wrap(err, "failed to presign object get")
	}

	url, err := url.Parse(signed)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse presigned url")
	}

	return url, err
}

// TODO move out of connector
func (s *Storage) IsContentTypeAllowed(docType types.DocumentType, mediaType string) bool {
	return s.mediaTypes.IsAllowed(docType, mediaType)
}

func (s *Storage) UploadFormData(key string) (map[string]string, error) {
	policy := minio.NewPostPolicy()
	policy.String()
	policy.SetBucket(s.bucket)
	policy.SetKey(strings.ToLower(key))
	policy.SetExpires(time.Now().Add(presignExpire))

	url, formData, err := s.minio.PresignedPostPolicy(policy)
	if err != nil {
		panic(err)
	}
	formData["url"] = url.String()

	return formData, nil
}
