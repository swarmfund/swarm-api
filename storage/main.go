package storage

import (
	"net/url"
	"strings"
	"time"

	"github.com/minio/minio-go"
	"github.com/minio/minio-go/pkg/policy"
	"github.com/pkg/errors"
	"gitlab.com/swarmfund/api/log"
)

const (
	bucketName   = "api"
	publicPrefix = "dpu"
)

type Connector struct {
	Minio            *minio.Client
	Log              *log.Entry
	MinContentLength int64
	MaxContentLength int64
}

func (c *Connector) ensureInitialized() error {
	if err := c.makeBucket(bucketName); err != nil {
		return errors.Wrap(err, "failed to create bucket")
	}

	err := c.Minio.SetBucketPolicy(bucketName, publicPrefix, policy.BucketPolicyReadOnly)
	if err != nil {
		return errors.Wrap(err, "failed to ensure bucket policy")
	}

	return nil
}

func (c *Connector) UploadFormData(key string) (map[string]string, error) {
	policy := minio.NewPostPolicy()

	policy.SetBucket(bucketName)
	policy.SetKey(strings.ToLower(key))
	// TODO investigate expire
	policy.SetExpires(time.Now().UTC().Add(72 * time.Hour))
	policy.SetContentLengthRange(c.MinContentLength, c.MaxContentLength)

	url, formData, err := c.Minio.PresignedPostPolicy(policy)
	if err != nil {
		if minioErr, ok := err.(minio.ErrorResponse); ok {
			if minioErr.Code == "NoSuchBucket" {
				// bucket does not exists, yet
				err = c.ensureInitialized()
				if err != nil {
					return nil, err
				}
				return c.UploadFormData(key)
			} else {
				return nil, err
			}
		} else {
			return nil, err
		}
	}
	formData["url"] = url.String()

	return formData, nil
}

func (c *Connector) makeBucket(bucket string) error {
	bucket = strings.ToLower(bucket)
	err := c.Minio.MakeBucket(bucket, "")
	if err != nil {
		if minioErr, ok := err.(minio.ErrorResponse); ok {
			if minioErr.Code == "BucketAlreadyOwnedByYou" {
				// seems like race on bucket create
				err = nil
			}
		}
		// check error again, it could get reset above
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Connector) DocumentURL(key string) (*url.URL, error) {
	return c.Minio.PresignedGetObject(strings.ToLower(bucketName), strings.ToLower(key), 3600*time.Second, nil)
}

func (c *Connector) Delete(bucket, key string) error {
	return c.Minio.RemoveObject(strings.ToLower(bucket), strings.ToLower(key))
}

func (c *Connector) Get(bucket, key string) (*minio.Object, error) {
	return c.Minio.GetObject(strings.ToLower(bucket), strings.ToLower(key))
}

func (c *Connector) Exists(bucket, key string) (bool, error) {
	object, err := c.Minio.GetObject(strings.ToLower(bucket), strings.ToLower(key))
	if err != nil {
		return false, err
	}
	defer object.Close()
	_, err = object.Stat()
	if err != nil {
		if mErr, ok := err.(minio.ErrorResponse); ok {
			if mErr.Code == "NoSuchKey" || mErr.Code == "NoSuchBucket" {
				return false, nil
			}
		}
		return false, err
	}
	return true, err
}

func (q *Connector) Bucket(bucket string) ([]minio.ObjectInfo, error) {
	doneCh := make(chan struct{})
	defer close(doneCh)
	result := []minio.ObjectInfo{}
	for message := range q.Minio.ListObjects(strings.ToLower(bucket), "", false, doneCh) {
		if message.Err != nil {
			return nil, message.Err
		}
		result = append(result, message)
	}

	return result, nil
}

func (c *Connector) DeleteBucket(bucket string) error {
	bucket = strings.ToLower(bucket)
	objects, err := c.Bucket(bucket)
	if err != nil {
		if mErr, ok := err.(minio.ErrorResponse); ok {
			if mErr.Code == "NoSuchBucket" {
				return nil
			}
		}
		return err
	}
	rmCh := make(chan string)
	errCh := c.Minio.RemoveObjects(bucket, rmCh)
	for _, object := range objects {
		rmCh <- object.Key
	}
	close(rmCh)
	mErr := <-errCh
	if mErr.Err != nil {
		return mErr.Err
	}
	return c.Minio.RemoveBucket(bucket)
}
