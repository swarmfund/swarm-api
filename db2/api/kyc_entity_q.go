package api

import (
	"database/sql"

	"github.com/lann/squirrel"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"gitlab.com/swarmfund/api/internal/kyc"
)

type KYCQI interface {
	Create(uid int64, eid string, entity kyc.Entity) error
	Update(eid string, entity kyc.Entity) error
	Delete(eid int64) error
	Select(uid int64) ([]KYCEntityRecord, error)
	Get(eid string) (*KYCEntityRecord, error)
}

const (
	kycEntitiesTable                = "kyc_entities"
	kycEntitiesIndividualConstraint = "kyc_entities_individual_constraint"
)

var (
	ErrKYCEntitiesConstraintViolated = errors.New("kyc_entities constraint violated")
)

type KYCEntityRecord struct {
	ID     string     `db:"id"`
	Entity kyc.Entity `db:"data"`
}

type KYCQ struct {
	parent *Q
}

func (q *UsersQ) KYC() KYCQI {
	return &KYCQ{
		parent: q.parent,
	}
}

func (q *KYCQ) Select(uid int64) (result []KYCEntityRecord, err error) {
	stmt := squirrel.Select("id, data").
		From(kycEntitiesTable).
		Where("user_id = ?", uid)
	err = q.parent.Select(&result, stmt)
	return result, err
}

func (q *KYCQ) Get(eid string) (*KYCEntityRecord, error) {
	stmt := squirrel.Select("id", "data").
		From(kycEntitiesTable).
		Where("id = ?", eid)
	var result KYCEntityRecord
	err := q.parent.Get(&result, stmt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &result, err
}

func (q *KYCQ) Create(uid int64, eid string, entity kyc.Entity) error {
	stmt := squirrel.Insert(kycEntitiesTable).SetMap(map[string]interface{}{
		"id":      eid,
		"user_id": uid,
		"data":    entity,
		"type":    entity.Type,
	})
	_, err := q.parent.Exec(stmt)
	if err != nil {
		cause := errors.Cause(err)
		pqerr, ok := cause.(*pq.Error)
		if ok {
			if pqerr.Constraint == kycEntitiesIndividualConstraint {
				return ErrKYCEntitiesConstraintViolated
			}
		}
	}
	return err
}

func (q *KYCQ) Update(eid string, entity kyc.Entity) error {
	stmt := squirrel.Update(kycEntitiesTable).
		Set("data", entity).
		Where("id = ?", eid)
	_, err := q.parent.Exec(stmt)
	return err
}

func (q *KYCQ) Delete(eid int64) error {
	stmt := squirrel.Delete("kyc_entities").Where("id = ?", eid)
	_, err := q.parent.Exec(stmt)
	return err
}
