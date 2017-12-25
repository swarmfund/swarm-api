package types

type DocumentKey struct {
	VersionByte  byte
	UserID       int64
	DocumentType DocumentType
	Random       [23]byte
}
