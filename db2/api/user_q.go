package api

import (
	"database/sql"
	"fmt"
	"time"

	sq "github.com/lann/squirrel"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"gitlab.com/swarmfund/api/internal/types"
)

const (
	tableUser      = "users"
	tableUserLimit = 20
)

var (
	tableUserAliased = fmt.Sprintf("%s u", tableUser)
	selectUser       = sq.Select(
		"u.*",
		"(select json_agg(kyc) from kyc_entities kyc where kyc.user_id=u.id) as kyc_entities").
		From(tableUserAliased)

	insertUser = sq.Insert(tableUser)
	updateUser = func(address string) sq.UpdateBuilder {
		return sq.Update(tableUserAliased).
			Where("address = ?", address).
			Set("updated_at", time.Now())
	}
	ErrUsersConflict = errors.New("address constraint violated")
)

//UsersQ is a helper struct to aid in configuring queries that loads users
type UsersQ struct {
	Err    error
	parent *Q
	sql    sq.SelectBuilder
}

type UsersQI interface {
	New() UsersQI
	Transaction(func(UsersQI) error) error

	// Participants
	Participants(participants map[int64][]Participant) error

	// Select methods
	ByState(state types.UserState) UsersQI
	RecoveryPending() UsersQI
	LimitReviewRequests() UsersQI

	ByAddress(address string) (*User, error)
	ByAddresses(addresses []string) ([]User, error)
	ByEmail(email string) (*User, error)
	ByID(uid int64) (*User, error)
	Select(dest interface{}) error

	// Create expected to set inserted record ID on u
	Create(u *User) error
	Update(u *User) error
	Delete(accountID string) error

	// Change state methods
	Approve(user *User) error
	Reject(user *User) error
	ChangeState(address types.Address, state types.UserState) error
	LimitReviewState(address string, state UserLimitReviewState) error

	// Update methods
	SetRecoveryState(address string, state UserRecoveryState) error
	SetEmail(address, email string) error

	WithAddress(addresses ...string) UsersQI
	Page(page uint64) UsersQI
	First() (*User, error)

	Documents(version int64) DocumentsQI
	KYC() KYCQI
}

func (q *Q) Users() UsersQI {
	return &UsersQ{
		parent: q,
		sql:    selectUser,
	}
}

func (q *UsersQ) New() UsersQI {
	return q.parent.Users()
}

func (q *UsersQ) Transaction(fn func(UsersQI) error) error {
	return q.parent.Transaction(func() error {
		return fn(q)
	})
}

func (q *UsersQ) WithIntegration(exchange string) UsersQI {
	if q.Err != nil {
		return q
	}
	q.sql = q.sql.
		Column("i.meta as integration_meta").
		LeftJoin("exchange_integrations i on i.address=u.address").
		Where("i.exchange = ?", exchange)
	return q
}

func (q *UsersQ) WithAddress(addresses ...string) UsersQI {
	if q.Err != nil {
		return q
	}
	q.sql = q.sql.
		Where(sq.Eq{"u.address": addresses})
	return q
}

func (q *UsersQ) ByState(state types.UserState) UsersQI {
	q.sql = q.sql.Where("state = ?", state)
	return q
}

func (q *UsersQ) RecoveryPending() UsersQI {
	q.sql = q.sql.Where("recovery_state = ?", UserRecoveryStatePending)
	return q
}

func (q *UsersQ) LimitReviewRequests() UsersQI {
	if q.Err != nil {
		return q
	}
	q.sql = q.sql.Where("limit_review_state = ?", UserLimitReviewPending)
	return q
}

func (q *UsersQ) Update(user *User) error {
	stmt := updateUser(string(user.Address)).
		Set("type", user.UserType)
	_, err := q.parent.Exec(stmt)
	return err
}

func (q *UsersQ) SetRecoveryState(address string, state UserRecoveryState) error {
	sql := updateUser(address).
		Set("recovery_state", state)
	_, err := q.parent.Exec(sql)
	return err
}

func (q *UsersQ) SetEmail(address, email string) error {
	sql := updateUser(address).
		Set("email", email)
	_, err := q.parent.Exec(sql)
	return err
}

//Approve updates row in `users`, set 'state' == approved
func (q *UsersQ) Approve(user *User) error {
	sql := updateUser(string(user.Address)).
		Where("documents_version = ?", user.DocumentsVersion).
		Set("documents_version", user.DocumentsVersion+1).
		Set("state", types.UserStateApproved).
		Set("documents", user.Documents)
	result, err := q.parent.Exec(sql)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		err = ErrBadDocumentVersion
	}
	return err
}

func (q *UsersQ) ChangeState(address types.Address, state types.UserState) error {
	sql := updateUser(string(address)).
		Set("state", state)

	_, err := q.parent.Exec(sql)
	return err
}

func (q *UsersQ) LimitReviewState(address string, state UserLimitReviewState) error {
	sql := updateUser(address).
		Set("limit_review_state", state)

	_, err := q.parent.Exec(sql)
	return err
}

// ByAddress loads a row from `users`, by address
func (q *UsersQ) ByAddress(address string) (*User, error) {
	dest := new(User)
	stmt := selectUser.
		Where("u.address = ?", address).
		Limit(1)
	err := q.parent.Get(dest, stmt)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	return dest, err
}

// ByEmail loads a row from `users`, by email
func (q *UsersQ) ByEmail(email string) (*User, error) {
	dest := new(User)
	sqlq := selectUser.Limit(1).Where("u.email = ?", email)

	err := q.parent.Get(dest, sqlq)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	return dest, err
}

func (q *UsersQ) ByID(uid int64) (*User, error) {
	var user User
	stmt := selectUser.Limit(1).Where("u.id = ?", uid)
	err := q.parent.Get(&user, stmt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &user, err
}

func (q *UsersQ) Create(u *User) error {
	sql := insertUser.SetMap(map[string]interface{}{
		"address": u.Address,
		"email":   u.Email,
		"type":    u.UserType,
		"state":   u.State,
	}).Suffix("returning id")

	err := q.parent.Get(&(u.ID), sql)
	if err != nil {
		cause := errors.Cause(err)
		pqerr, ok := cause.(*pq.Error)
		if ok {
			if pqerr.Constraint == "unique_address" {
				return ErrUsersConflict
			}
		}
	}
	return err
}

func (q *UsersQ) Reject(user *User) error {
	sql := fmt.Sprintf(`update %s
		set documents = $1,
			documents_version = $2,
		 	state = $3,
		    updated_at = timestamp 'now'
		where address = $4
		  and documents_version = $5`, tableUser)

	result, err := q.parent.DB.Exec(sql, user.Documents, user.DocumentsVersion+1, types.UserStateRejected, user.Address, user.DocumentsVersion)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		err = ErrBadDocumentVersion
	}

	return err
}

func (q *UsersQ) ByAddresses(addresses []string) (result []User, err error) {
	if len(addresses) == 0 {
		return result, nil
	}

	sql := selectUser.Where(sq.Eq{"u.address": addresses})

	err = q.parent.Select(&result, sql)
	return result, err
}

func (q *UsersQ) Page(page uint64) UsersQI {
	if q.Err != nil {
		return q
	}

	q.sql = q.sql.Offset(tableUserLimit * (page - 1)).Limit(tableUserLimit)
	return q
}

// Select loads the results of the query specified by `q` into `dest`.
func (q *UsersQ) Select(dest interface{}) error {
	if q.Err != nil {
		return q.Err
	}

	q.Err = q.parent.Select(dest, q.sql)
	return q.Err
}

func (q *UsersQ) First() (*User, error) {
	if q.Err != nil {
		return nil, q.Err
	}
	var result User
	err := q.parent.Get(&result, q.sql.Limit(1))
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &result, err
}

func (q *UsersQ) Delete(accountID string) error {
	stmt := sq.Delete(tableUser).Where("address = ?", accountID)
	_, err := q.parent.Exec(stmt)
	return err
}

func (q *UsersQ) Participants(ops map[int64][]Participant) error {
	var accountIDs []types.Address
	for _, op := range ops {
		for _, participant := range op {
			accountIDs = append(accountIDs, participant.AccountID)
		}
	}

	if len(accountIDs) == 0 {
		return nil
	}

	stmt := selectUser.Where(sq.Eq{"u.address": accountIDs})

	var users []User
	err := q.parent.Select(&users, stmt)
	if err != nil {
		return err
	}

	usersMap := map[types.Address]User{}
	for _, user := range users {
		usersMap[user.Address] = user
	}

	for _, op := range ops {
		for pi, _ := range op {
			participant := op[pi]
			if user, ok := usersMap[participant.AccountID]; ok {
				participant.Email = &user.Email
				op[pi] = participant
			}
		}
	}

	return nil
}
