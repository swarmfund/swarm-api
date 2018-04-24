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
	tableUserLimit = 100
)

var (
	tableUserAliased = fmt.Sprintf("%s u", tableUser)
	selectUser       = sq.Select(
		"u.id",
		"u.email",
		"u.address",
		"u.kyc_sequence",
		"u.reject_reason",
		"r.address as recovery_address",
		"a.state as airdrop_state",
		"(select json_agg(kyc) from kyc_entities kyc where kyc.user_id=u.id) as kyc_entities",
		"b.value as kyc_blob_value",
		"b.id as kyc_blob_id",
		// state and type might be nil if ingestion is still in progress,
		// make sure resources render some meaningful stub values that will not break clients
		"coalesce(us.state, 0) as user_state",
		"coalesce(us.type, 0) as user_type",
	).
		LeftJoin("user_states us on us.address=u.address").
		Join("recoveries r on r.wallet=u.email").
		// joining left since it's optional due to late migration
		LeftJoin("airdrops a on a.owner=u.address").
		LeftJoin("blobs b ON us.kyc_blob = b.id and b.type = ?", types.BlobTypeKYCForm).
		From(tableUserAliased)

	insertUser = sq.Insert(tableUser)
	updateUser = func(address string) sq.UpdateBuilder {
		return sq.Update(tableUserAliased).
			Where("u.address = ?", address).
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
	UserStateQ

	New() UsersQI
	Transaction(func(UsersQI) error) error

	UpdateAirdropState(address types.Address, state types.AirdropState) error

	// DEPRECATED
	Participants(participants map[int64][]Participant) error

	// Select methods
	ByState(state types.UserState) UsersQI
	ByType(tpe types.UserType) UsersQI
	EmailMatches(string) UsersQI
	AddressMatches(string) UsersQI
	ByFirstName(firstName string) UsersQI
	ByLastName(lastName string) UsersQI
	ByCountry(country string) UsersQI
	Select(dest interface{}) error
	Page(page uint64) UsersQI

	ByAddress(address string) (*User, error)
	// DEPRECATED
	ByAddresses(addresses []string) ([]User, error)
	// DEPRECATED
	ByEmail(email string) (*User, error)

	// Create expected to set inserted record ID on u
	Create(u *User) error
	Update(u *User) error

	KYC() KYCQI
}

func (q *Q) Users() UsersQI {
	return &UsersQ{
		parent: &Q{
			q.Clone(),
		},
		sql: selectUser,
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

func (q *UsersQ) UpdateAirdropState(address types.Address, state types.AirdropState) error {
	stmt := sq.Insert("airdrops").SetMap(map[string]interface{}{
		"owner": address,
		"state": state,
	}).Suffix("ON CONFLICT (owner) DO UPDATE SET state = EXCLUDED.state;")

	_, err := q.parent.Exec(stmt)
	return err
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

func (q *UsersQ) ByState(state types.UserState) UsersQI {
	q.sql = q.sql.Where("us.state & ? != 0", state)
	return q
}

func (q *UsersQ) ByType(tpe types.UserType) UsersQI {
	q.sql = q.sql.Where("us.type & ? != 0", tpe)
	return q
}

func (q *UsersQ) EmailMatches(str string) UsersQI {
	q.sql = q.sql.Where("u.email ilike ?", fmt.Sprint("%", str, "%"))
	return q
}

func (q *UsersQ) AddressMatches(str string) UsersQI {
	q.sql = q.sql.Where("u.address ilike ?", fmt.Sprint("%", str, "%"))
	return q
}

func (q *UsersQ) Update(user *User) error {
	stmt := updateUser(string(user.Address)).
		SetMap(map[string]interface{}{
			"type":          user.UserType,
			"state":         user.State,
			"kyc_sequence":  user.KYCSequence,
			"reject_reason": user.RejectReason,
		})
	_, err := q.parent.Exec(stmt)
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

func (q *UsersQ) SetState(update UserStateUpdate) error {
	clauses := map[string]interface{}{
		"address":    update.Address,
		"updated_at": update.Timestamp,
	}
	if update.State != nil {
		clauses["state"] = *update.State
	}
	if update.Type != nil {
		clauses["type"] = *update.Type
	}
	if update.KYCBlob != nil {
		clauses["kyc_blob"] = *update.KYCBlob
	}

	stmt := sq.Insert("user_states as us").SetMap(clauses).Suffix(`
		ON CONFLICT (address) DO UPDATE
			SET state = coalesce(excluded.state, us.state),
				type = coalesce(excluded.type, us.type),
				updated_at = excluded.updated_at,
				kyc_blob = coalesce(excluded.kyc_blob, us.kyc_blob)
			WHERE us.updated_at <= excluded.updated_at
	`)
	_, err := q.parent.Exec(stmt)
	if err != nil {
		cause := errors.Cause(err)
		pqerr, ok := cause.(*pq.Error)
		if ok {
			if pqerr.Constraint == "user_states_users_fkey" {
				return ErrUsersConflict
			}
		}
	}
	return err
}

func (q *UsersQ) Create(u *User) error {
	sql := insertUser.SetMap(map[string]interface{}{
		"address": u.Address,
		"email":   u.Email,
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

func (q *UsersQ) ByFirstName(firstName string) UsersQI {
	q.sql = q.sql.Where("? in (b.value::jsonb#>>'{first_name}', b.value::jsonb#>>'{v2,first_name}')", firstName)
	return q
}

func (q *UsersQ) ByLastName(lastName string) UsersQI {
	q.sql = q.sql.Where("? in (b.value::jsonb#>>'{last_name}', b.value::jsonb#>>'{v2,last_name}')", lastName)
	return q
}

func (q *UsersQ) ByCountry(country string) UsersQI {
	q.sql = q.sql.Where("? in (b.value::jsonb#>>'{address, country}', b.value::jsonb#>>'{v2,address,country}')", country)
	return q
}
