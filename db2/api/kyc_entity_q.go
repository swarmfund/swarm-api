package api

import (
	"github.com/lann/squirrel"
)

type KYCQI interface {
	Create(entity KYCEntity) (int64, error)
	Update(eid int64, data []byte) error
	Delete(eid int64) error
}

type KYCQ struct {
	parent *Q
}

func (q *UsersQ) KYC() KYCQI {
	return &KYCQ{
		parent: q.parent,
	}
}

func (q *KYCQ) Create(entity KYCEntity) (int64, error) {
	var result int64
	stmt := squirrel.Insert("kyc_entities").SetMap(map[string]interface{}{
		"user_id": entity.UserID,
		"data":    entity.Data,
		"type":    entity.Type,
	}).Suffix("returning id")
	err := q.parent.Get(&result, stmt)
	return result, err
}

func (q *KYCQ) Update(eid int64, data []byte) error {
	stmt := squirrel.Update("kyc_entities").
		Set("data", data).
		Where("id = ?", eid)
	_, err := q.parent.Exec(stmt)
	return err
}

func (q *KYCQ) Delete(eid int64) error {
	stmt := squirrel.Delete("kyc_entities").Where("id = ?", eid)
	_, err := q.parent.Exec(stmt)
	return err
}
