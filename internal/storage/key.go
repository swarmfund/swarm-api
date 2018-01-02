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

const (
	documentByte byte = 27
	publicByte   byte = 233
	privateByte  byte = 48
)

var (
	keyMask = []byte{
		0x00, 0x00, 0x12, 0x1f, 0x88, 0x69, 0x5b, 0x40,
		0xa4, 0xe2, 0x74, 0xa7, 0xf2, 0xd5, 0xce, 0x29,
		0x94, 0xd5, 0x62, 0x3a, 0xd8, 0xf4, 0x24, 0xee,
		0xc8, 0x65, 0x8a, 0x39, 0xb2, 0xa6, 0xfc, 0x1d,
		0xdc, 0x19, 0x38, 0xc3, 0xae, 0xeb, 0x92, 0x8a,
	}
)

type Key struct {
	UserID  int64
	Type    types.DocumentType
	entropy [21]byte
}

func NewKey(uid int64, documentType types.DocumentType) *Key {
	key := Key{
		UserID: uid,
		Type:   documentType,
	}
	_, err := rand.Read(key.entropy[:])
	if err != nil {
		panic(errors.Wrap(err, "failed to read entropy"))
	}
	return &key
}

func (k *Key) MarshalText() ([]byte, error) {
	buf := new(bytes.Buffer)

	buf.WriteByte(documentByte)

	if types.IsPublicDocument(k.Type) {
		buf.WriteByte(publicByte)
	} else {
		buf.WriteByte(privateByte)
	}

	fields := []interface{}{
		k.Type, k.UserID,
	}
	for _, field := range fields {
		if err := binary.Write(buf, binary.LittleEndian, field); err != nil {
			return nil, errors.Wrap(err, "failed to write field")
		}
	}

	buf.Write(k.entropy[:])

	raw := make([]byte, len(buf.Bytes()))
	for i := 0; i < len(raw); i++ {
		raw[i] = buf.Bytes()[i] ^ keyMask[i]
	}

	encoded := base32.StdEncoding.EncodeToString(raw)
	return []byte(strings.ToLower(encoded)), nil
}

func (k *Key) UnmarshalText(src []byte) error {
	str := strings.ToUpper(string(src))
	decoded, err := base32.StdEncoding.DecodeString(str)
	if err != nil {
		return errors.Wrap(err, "failed to decode")
	}

	raw := make([]byte, len(decoded))
	for i := 0; i < len(raw); i++ {
		raw[i] = decoded[i] ^ keyMask[i]
	}

	// ignoring document and visibility bytes
	// TODO check leading bytes
	reader := bytes.NewReader(raw[2:])

	fields := []interface{}{
		&k.Type, &k.UserID,
	}
	for _, field := range fields {
		if err := binary.Read(reader, binary.LittleEndian, field); err != nil {
			return errors.Wrap(err, "failed to read field")
		}
	}

	_, err = reader.Read(k.entropy[:])
	if err != nil {
		return errors.Wrap(err, "failed to read entropyste")
	}
	return nil
}

func EncodeKey(key *Key) string {
	encoded, err := key.MarshalText()
	if err != nil {
		panic(errors.Wrap(err, "failed to encode key"))
	}
	return string(encoded)
}
