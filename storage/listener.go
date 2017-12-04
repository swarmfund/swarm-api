package storage

import (
	"encoding/json"

	"strings"
	"time"

	"github.com/minio/minio-go"
	"github.com/pkg/errors"
	"github.com/streadway/amqp"
	"gitlab.com/swarmfund/api/config"
	"gitlab.com/swarmfund/api/db2/api"
	errors2 "gitlab.com/swarmfund/api/errors"
	"gitlab.com/swarmfund/api/log"
)

type Consumer struct {
	conn     *amqp.Connection
	ch       *amqp.Channel
	tag      string
	messages <-chan amqp.Delivery
	log      *log.Entry
	apiQ     api.QInterface
	conf     config.Storage
	storage  *Connector
}

func NewConsumer(conf config.Storage, apiQ api.QInterface, storage *Connector) (*Consumer, error) {
	conn, err := amqp.Dial(conf.ListenerBrokerURL)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return &Consumer{
		conn:    conn,
		ch:      ch,
		log:     log.WithField("service", "document-consumer"),
		apiQ:    apiQ,
		conf:    conf,
		storage: storage,
	}, nil
}

func (c *Consumer) Prepare() error {
	err := c.ch.ExchangeDeclare(
		c.conf.ListenerExchange,
		c.conf.ListenerExchangeType,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	q, err := c.ch.QueueDeclare(
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	err = c.ch.Qos(
		1,
		0,
		false,
	)
	if err != nil {
		return err
	}

	err = c.ch.QueueBind(
		q.Name,
		c.conf.ListenerBindingKey,
		c.conf.ListenerExchange,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	c.messages, err = c.ch.Consume(
		q.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	return err
}

func (c *Consumer) Consume() {
	go func() {
		defer c.TearDown()
		for msg := range c.messages {
			if err := c.processMessage(msg); err != nil {
				c.log.WithError(err).Error("failed to process message")
				// TODO dead-letter queue
				msg.Nack(false, true)
				continue
			}
			msg.Ack(false)
		}
	}()
}

func (c *Consumer) processMessage(msg amqp.Delivery) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors2.FromPanic(r)
		}
	}()
	var message Message
	if err = json.Unmarshal(msg.Body, &message); err != nil {
		return errors.Wrap(err, "failed to unmarshal")
	}
	err = c.processRecords(message.Records)
	if err != nil {
		return errors.Wrap(err, "failed to process records")
	}
	return nil
}

func (c *Consumer) processRecords(records []Record) error {
	for _, record := range records {
		switch record.EventName {
		case minio.ObjectCreatedPut, minio.ObjectCreatedPost, minio.ObjectCreatedCopy, minio.ObjectCreatedCompleteMultipartUpload:
			return c.processObjectCreated(record)
		case minio.ObjectRemovedDelete, minio.ObjectRemovedDeleteMarkerCreated:
			return c.processObjectRemoved(record)
		default:
			c.log.
				WithField("event", record.EventName).
				Warn("got unexpected event type")
			continue
		}
	}
	return nil
}

func (c *Consumer) processObjectCreated(record Record) error {
	accountID := strings.ToUpper(record.S3.Bucket.Name)

	// update user docs prior to triggering specific type processing
	document, err := c.updateUserDocs(accountID, record)
	if err != nil {
		return err
	}

	user, err := c.apiQ.Users().ByAddress(accountID)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("user not found")
	}

	switch {
	case document.Type.IsRecovery():
		err := c.ProcessRecoveryUpload(user, document)
		if err != nil {
			return err
		}
	case document.Type.IsKYC():
		err := c.ProcessKYCUpload(user, document)
		if err != nil {
			return err
		}
	case document.Type.IsProofOfIncome():
		err := c.ProcessProofOfIncomeUpload(user, document)
		if err != nil {
			return err
		}
	default:
		return errors.New("unknown document type uploaded")
	}
	return nil
}

func (c *Consumer) processObjectRemoved(record Record) error {
	accountID := strings.ToUpper(record.S3.Bucket.Name)

	doc, err := FromKey(record.S3.Object.Key)
	if err != nil {
		return errors.Wrap(err, "failed to parse key")
	}

	user, err := c.apiQ.Users().ByAddress(accountID)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("user not found")
	}

	err = c.apiQ.Users().Documents(user.DocumentsVersion).Delete(user.ID, doc.Type, record.S3.Object.Key)
	return err
}

func (c *Consumer) updateUserDocs(accountID string, record Record) (*api.Document, error) {
	user, err := c.apiQ.Users().ByAddress(accountID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		// TODO shouldn't really happen, probably need to delete doc
		return nil, errors.New("user not found")
	}

	doc, err := FromKey(record.S3.Object.Key)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse key")
	}

	document := &api.Document{
		ContentType: record.S3.Object.ContentType,
		Key:         record.S3.Object.Key,
		Type:        doc.Type,
		EntityID:    doc.EntityID,
		Checksum:    record.S3.Object.ETag,
		CreatedAt:   time.Now().Unix(),
		Version:     doc.Version,
	}

	err = c.apiQ.Users().Documents(user.DocumentsVersion).Set(user.ID, document)
	if err != nil {
		return nil, err
	}

	return document, nil
}

func (c *Consumer) TearDown() {
	c.conn.Close()
	c.ch.Close()
}

type Message struct {
	Records []Record `json:"Records"`
}

type Record struct {
	S3 struct {
		Bucket struct {
			Name string `json:"name"`
		} `json:"bucket"`
		Object struct {
			Key         string `json:"key"`
			ETag        string `json:"eTag"`
			ContentType string `json:"contentType"`
		} `json:"object"`
	} `json:"s3"`
	EventName string    `json:"eventName"`
	EventTime time.Time `json:"eventTime"`
}
