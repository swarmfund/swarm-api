package storage

import (
	"errors"
	"net/url"
	"strings"
	"time"

	"gitlab.com/swarmfund/api/config"
	"gitlab.com/swarmfund/api/log"
	"github.com/minio/minio-go"
)

type Connector struct {
	minio *minio.Client
	log   *log.Entry
	conf  config.Storage
}

func New(conf config.Storage) (*Connector, error) {
	minioClient, err := minio.New(conf.Host, conf.AccessKey, conf.SecretKey, conf.ForceSSL)
	if err != nil {
		return nil, err
	}
	return &Connector{
		minio: minioClient,
		log:   log.WithField("service", "storage"),
		conf:  conf,
	}, nil
}

func (c *Connector) UploadFormData(bucket, key string) (map[string]string, error) {
	policy := minio.NewPostPolicy()

	policy.SetBucket(strings.ToLower(bucket))
	policy.SetKey(strings.ToLower(key))
	policy.SetExpires(time.Now().UTC().Add(c.conf.FormDataExpire))
	policy.SetContentLengthRange(c.conf.MinContentLength, c.conf.MaxContentLength)

	url, formData, err := c.minio.PresignedPostPolicy(policy)
	if err != nil {
		if minioErr, ok := err.(minio.ErrorResponse); ok {
			if minioErr.Code == "NoSuchBucket" {
				// bucket does not exists, yet
				err = c.makeBucket(bucket)
				if err != nil {
					return nil, err
				}
				return c.UploadFormData(bucket, key)
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
	err := c.minio.MakeBucket(bucket, "")
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

	fields := strings.Split(c.conf.ObjectCreateARN, ":")
	if len(fields) != 6 {
		return errors.New("invalid arn")
	}

	mb, err := c.minio.GetBucketNotification(bucket)
	if err != nil {
		return err
	}

	arn := minio.NewArn(fields[1], fields[2], fields[3], fields[4], fields[5])
	nc := minio.NewNotificationConfig(arn)
	nc.AddEvents(minio.ObjectRemovedAll, minio.ObjectCreatedAll)
	switch fields[2] {
	case "sns":
		mb.AddTopic(nc)
	case "sqs":
		mb.AddQueue(nc)
	case "lambda":
		mb.AddLambda(nc)
	default:
		return errors.New("invalid arn service")
	}

	if err := c.minio.SetBucketNotification(bucket, mb); err != nil {
		return err
	}
	return nil
}

func (c *Connector) DocumentURL(bucket, key string) (*url.URL, error) {
	return c.minio.PresignedGetObject(strings.ToLower(bucket), strings.ToLower(key), 3600*time.Second, nil)
}

func (c *Connector) Delete(bucket, key string) error {
	return c.minio.RemoveObject(strings.ToLower(bucket), strings.ToLower(key))
}

func (c *Connector) Get(bucket, key string) (*minio.Object, error) {
	return c.minio.GetObject(strings.ToLower(bucket), strings.ToLower(key))
}

func (c *Connector) Exists(bucket, key string) (bool, error) {
	object, err := c.minio.GetObject(strings.ToLower(bucket), strings.ToLower(key))
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
	for message := range q.minio.ListObjects(strings.ToLower(bucket), "", false, doneCh) {
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
	errCh := c.minio.RemoveObjects(bucket, rmCh)
	for _, object := range objects {
		rmCh <- object.Key
	}
	close(rmCh)
	mErr := <-errCh
	if mErr.Err != nil {
		return mErr.Err
	}
	return c.minio.RemoveBucket(bucket)
}
