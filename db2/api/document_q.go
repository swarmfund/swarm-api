package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/lib/pq"
)

var (
	ErrBadDocumentVersion = errors.New("bad document version")
)

func quoteEscape(s string) string {
	return strings.Replace(s, `'`, `''`, -1)
}

// plsq json magic ahead,
// if your intention is to add more just move documents to separate table and implement proper locking
// takes documents[document.type], default is {}
// creates {document.key: document} object
// merges two in a way when same key overrides value
// then builds {document.type: {object_create_above}}
// strips attributes with null values
// and finally merges it to user documents object
func updateUserDocument(docType DocumentType, key, document string) string {
	return fmt.Sprintf(`
		update %s
		set documents = documents || jsonb_strip_nulls(jsonb_build_object('%d', coalesce(documents -> '%d', '{}'::jsonb) || jsonb_build_object('%s', '%s'::jsonb))),
		    updated_at = timestamp 'now',
		    documents_version = documents_version + 1
		where id = $1
		  and documents_version = $2`, pq.QuoteIdentifier(tableUser), docType, docType, quoteEscape(key), quoteEscape(document))
}

type DocumentsQI interface {
	Set(userID int64, document *Document) error
	Delete(userID int64, docType DocumentType, key string) error
}

type DocumentsQ struct {
	version int64
	parent  *Q
}

func (q *UsersQ) Documents(version int64) DocumentsQI {
	return &DocumentsQ{
		version: version,
		parent:  q.parent,
	}
}

func (q *DocumentsQ) Set(userID int64, document *Document) error {
	bytes, err := json.Marshal(document)
	if err != nil {
		return err
	}
	result, err := q.parent.DB.Exec(updateUserDocument(document.Type, document.Key, string(bytes)), userID, q.version)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err == nil && rows == 0 {
		err = ErrBadDocumentVersion
	}
	return err
}

func (q *DocumentsQ) Delete(userID int64, docType DocumentType, key string) error {
	result, err := q.parent.DB.Exec(updateUserDocument(docType, key, "null"), userID, q.version)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err == nil && rows == 0 {
		err = ErrBadDocumentVersion
	}
	return err
}
