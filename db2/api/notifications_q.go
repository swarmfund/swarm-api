package api

import (
	"encoding/json"
	"fmt"

	"database/sql"

	"github.com/lann/squirrel"
)

type NotificationType int

const (
	_ NotificationType = iota
	NotificationTypeKYC
)

type NotificationsQI interface {
	Get(address string) (*Notifications, error)
	Disable(address string, tpe NotificationType) error
	Enable(address string, tpe NotificationType) error
	GetRecipients(tpe NotificationType) ([]string, error)
	SetEmail(address, email string) error
}

type NotificationsQ struct {
	parent *Q
}

type Notifications struct {
	Records NotificationRecord `json:"records"`
}

type NotificationRecord []Notification

func (n *NotificationRecord) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	switch v := value.(type) {
	case []byte:
		return json.Unmarshal(v, &n)
	default:
		return fmt.Errorf("unsupported Scan from type %T", v)
	}
}

type Notification struct {
	Type  NotificationType `json:"type"`
	Email string           `json:"email"`
}

func (q *Q) Notifications() NotificationsQI {
	return &NotificationsQ{
		parent: q,
	}
}

func (q *NotificationsQ) SetEmail(address, email string) error {
	stmt := squirrel.Update("notifications").Set("email", email).Where("address = ?", address)
	_, err := q.parent.Exec(stmt)
	return err
}

func (q *NotificationsQ) GetRecipients(tpe NotificationType) ([]string, error) {
	var result []string
	stmt := squirrel.
		Select("email").
		From("notifications").
		Where("type = ?", tpe).
		Where("email is not null")
	err := q.parent.Select(&result, stmt)
	return result, err
}

func (q *NotificationsQ) Get(address string) (*Notifications, error) {
	var result Notifications
	stmt := squirrel.Select("json_agg(n) as records").From("notifications n").Where("address = ?", address)
	err := q.parent.Get(&result, stmt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &result, err
}

func (q *NotificationsQ) Disable(address string, tpe NotificationType) error {
	stmt := squirrel.Delete("notifications").Where(squirrel.Eq{
		"address": address,
		"type":    tpe,
	})
	_, err := q.parent.Exec(stmt)
	return err
}

func (q *NotificationsQ) Enable(address string, tpe NotificationType) error {
	stmt := squirrel.Insert("notifications").SetMap(map[string]interface{}{
		"address": address,
		"type":    tpe,
	}).Suffix("on conflict do nothing")
	_, err := q.parent.Exec(stmt)
	return err
}
