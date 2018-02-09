// generated by jsonenums -tprefix=false -transform=snake -type=DocumentType; DO NOT EDIT
package types

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
)

func init() {
	// stubs for imports
	_ = json.Delim('s')
	_ = driver.Int32

}

var ErrDocumentTypeInvalid = errors.New("DocumentType is invalid")

func init() {
	var v DocumentType
	if _, ok := interface{}(v).(fmt.Stringer); ok {
		_DocumentTypeNameToValue = map[string]DocumentType{
			interface{}(DocumentTypeAssetLogo).(fmt.Stringer).String():     DocumentTypeAssetLogo,
			interface{}(DocumentTypeFundLogo).(fmt.Stringer).String():      DocumentTypeFundLogo,
			interface{}(DocumentTypeFundDocument).(fmt.Stringer).String():  DocumentTypeFundDocument,
			interface{}(DocumentTypeNavReport).(fmt.Stringer).String():     DocumentTypeNavReport,
			interface{}(DocumentTypeAlpha).(fmt.Stringer).String():         DocumentTypeAlpha,
			interface{}(DocumentTypeBravo).(fmt.Stringer).String():         DocumentTypeBravo,
			interface{}(DocumentTypeCharlie).(fmt.Stringer).String():       DocumentTypeCharlie,
			interface{}(DocumentTypeDelta).(fmt.Stringer).String():         DocumentTypeDelta,
			interface{}(DocumentTypeTokenTerms).(fmt.Stringer).String():    DocumentTypeTokenTerms,
			interface{}(DocumentTypeTokenMetrics).(fmt.Stringer).String():  DocumentTypeTokenMetrics,
			interface{}(DocumentTypeKYCIdDocument).(fmt.Stringer).String(): DocumentTypeKYCIdDocument,
			interface{}(DocumentTypeKYCPoa).(fmt.Stringer).String():        DocumentTypeKYCPoa,
		}
	}
}

var _DocumentTypeNameToValue = map[string]DocumentType{
	"asset_logo":      DocumentTypeAssetLogo,
	"fund_logo":       DocumentTypeFundLogo,
	"fund_document":   DocumentTypeFundDocument,
	"nav_report":      DocumentTypeNavReport,
	"alpha":           DocumentTypeAlpha,
	"bravo":           DocumentTypeBravo,
	"charlie":         DocumentTypeCharlie,
	"delta":           DocumentTypeDelta,
	"token_terms":     DocumentTypeTokenTerms,
	"token_metrics":   DocumentTypeTokenMetrics,
	"kyc_id_document": DocumentTypeKYCIdDocument,
	"kyc_poa":         DocumentTypeKYCPoa,
}

var _DocumentTypeValueToName = map[DocumentType]string{
	DocumentTypeAssetLogo:     "asset_logo",
	DocumentTypeFundLogo:      "fund_logo",
	DocumentTypeFundDocument:  "fund_document",
	DocumentTypeNavReport:     "nav_report",
	DocumentTypeAlpha:         "alpha",
	DocumentTypeBravo:         "bravo",
	DocumentTypeCharlie:       "charlie",
	DocumentTypeDelta:         "delta",
	DocumentTypeTokenTerms:    "token_terms",
	DocumentTypeTokenMetrics:  "token_metrics",
	DocumentTypeKYCIdDocument: "kyc_id_document",
	DocumentTypeKYCPoa:        "kyc_poa",
}

// String is generated so DocumentType satisfies fmt.Stringer.
func (r DocumentType) String() string {
	s, ok := _DocumentTypeValueToName[r]
	if !ok {
		return fmt.Sprintf("DocumentType(%d)", r)
	}
	return s
}

// Validate verifies that value is predefined for DocumentType.
func (r DocumentType) Validate() error {
	_, ok := _DocumentTypeValueToName[r]
	if !ok {
		return ErrDocumentTypeInvalid
	}
	return nil
}

// MarshalJSON is generated so DocumentType satisfies json.Marshaler.
func (r DocumentType) MarshalJSON() ([]byte, error) {
	if s, ok := interface{}(r).(fmt.Stringer); ok {
		return json.Marshal(s.String())
	}
	s, ok := _DocumentTypeValueToName[r]
	if !ok {
		return nil, fmt.Errorf("invalid DocumentType: %d", r)
	}
	return json.Marshal(s)
}

// UnmarshalJSON is generated so DocumentType satisfies json.Unmarshaler.
func (r *DocumentType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("DocumentType should be a string, got %s", data)
	}
	v, ok := _DocumentTypeNameToValue[s]
	if !ok {
		return fmt.Errorf("invalid DocumentType %q", s)
	}
	*r = v
	return nil
}

func (t *DocumentType) Scan(src interface{}) error {
	i, ok := src.(int64)
	if !ok {
		return fmt.Errorf("can't scan from %T", src)
	}
	*t = DocumentType(i)
	return nil
}

func (t DocumentType) Value() (driver.Value, error) {
	return int64(t), nil
}
