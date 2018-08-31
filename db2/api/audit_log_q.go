package api

import (
	sq "github.com/lann/squirrel"
	"github.com/pkg/errors"
	"gitlab.com/swarmfund/api/internal/types"
)

const tableAuditLogs = "audit_logs"

var auditLogsSelect = sq.Select("*").From(tableAuditLogs)
var auditLogsInsert = sq.Insert(tableAuditLogs)

type AuditLogQI interface {
	New() AuditLogQI

	Create(userAction *AuditLogAction) error
	ByAddress(address types.Address) ([]AuditLogAction, error)
	LastSeen(address types.Address) (*AuditLogAction, error)
}

type AuditLogQ struct {
	parent *Q
	sql    sq.SelectBuilder
}

func (q *Q) AuditLog() AuditLogQI {
	return &AuditLogQ{
		parent: &Q{
			q.Clone(),
		},
		sql: auditLogsSelect,
	}
}

func (q *AuditLogQ) New() AuditLogQI {
	return q.parent.AuditLog()
}

func (q AuditLogQ) Create(action *AuditLogAction) (err error) {
	sql := auditLogsInsert.SetMap(map[string]interface{}{
		"user_address": action.UserAddress,
		"action":       action.ActionType,
		"performed_at": action.PerformedAt,
		"details":      action.Details,
	})

	_, err = q.parent.Exec(sql)
	return err
}

func (q AuditLogQ) ByAddress(address types.Address) ([]AuditLogAction, error) {
	sql := q.sql.Where("user_address = ?", address)

	var result []AuditLogAction

	err := q.parent.Select(&result, sql)
	return result, err
}
func (q AuditLogQ) LastSeen(address types.Address) (*AuditLogAction, error) {
	sql := q.sql.Where("user_address = ?", address).
		Where("details::jsonb#>>'{geoinfo, ip}' is not null").
		OrderBy("performed_at desc").
		Limit(1)

	var result []AuditLogAction

	err := q.parent.Select(&result, sql)

	if err != nil {
		return nil, err
	}

	if result == nil {
		return nil, errors.New("no logs found")
	}
	return &result[0], nil
}
