package api

import "github.com/pkg/errors"

type KYCTrackerQI interface {
	EnsureRecord(uid int64, state UserState) (bool, error)
}

func (q *Q) KYCTracker() KYCTrackerQI {
	return &KYCTrackerQ{
		parent: q,
	}
}

type KYCTrackerQ struct {
	parent *Q
}

func (q *KYCTrackerQ) EnsureRecord(uid int64, state UserState) (bool, error) {
	stmt := `
		insert into kyc_tracker (user_id, last_state) values ($1, $2)
		on conflict (user_id)
			do update set last_state=$2
			   where kyc_tracker.last_state != $2`
	result, err := q.parent.DB.Exec(stmt, uid, state)
	if err != nil {
		return false, err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return false, errors.Wrap(err, "failed to get rows affected")
	}
	return affected != 0, nil
}
