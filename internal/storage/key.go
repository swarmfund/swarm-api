package storage

import (
	"crypto/rand"

	"bytes"
	"encoding/base32"
	"encoding/binary"

	"strings"

	"github.com/pkg/errors"
	"gitlab.com/swarmfund/api/internal/types"
)

type VersionByte byte

const (
	VersionByteDocument VersionByte = 104
)

var (
	keyMask = []byte{
		0, 37, 199, 88, 239, 73, 154, 8,
		47, 246, 43, 30, 74, 11, 52, 73,
		62, 28, 45, 145, 46, 197, 162, 206,
		3, 26, 43, 8, 70, 205, 126, 12,
	}
)

type Key struct {
	VersionByte VersionByte
	UserID      int64
	Type        types.DocumentType
	Random      [19]byte
}

func NewKey(uid int64, documentType types.DocumentType) *Key {
	key := Key{
		VersionByte: VersionByteDocument,
		UserID:      uid,
		Type:        documentType,
	}
	_, err := rand.Read(key.Random[:])
	if err != nil {
		panic(errors.Wrap(err, "failed to read entropy"))
	}
	return &key
}

func (k *Key) MarshalText() ([]byte, error) {
	buf := new(bytes.Buffer)

	fields := []interface{}{
		k.VersionByte, k.UserID, k.Type, k.Random,
	}
	for _, field := range fields {
		if err := binary.Write(buf, binary.BigEndian, field); err != nil {
			return nil, errors.Wrap(err, "failed to write field")
		}
	}

	raw := [32]byte{}
	for i := 0; i < len(raw); i++ {
		raw[i] = buf.Bytes()[i] ^ keyMask[i]
	}

	encoded := base32.HexEncoding.WithPadding(base32.NoPadding).EncodeToString(raw[:])
	return []byte(strings.ToLower(encoded)), nil
}

func (k *Key) UnmarshalText(src []byte) error {
	str := strings.ToUpper(string(src))
	decoded, err := base32.HexEncoding.WithPadding(base32.NoPadding).DecodeString(str)
	if err != nil {
		return errors.Wrap(err, "failed to decode")
	}

	raw := [32]byte{}
	for i := 0; i < len(raw); i++ {
		raw[i] = decoded[i] ^ keyMask[i]
	}

	reader := bytes.NewReader(raw[:])
	fields := []interface{}{
		&k.VersionByte, &k.UserID, &k.Type, &k.Random,
	}
	for _, field := range fields {
		if err := binary.Read(reader, binary.BigEndian, field); err != nil {
			return errors.Wrap(err, "failed to read field")
		}
	}

	return nil
}
