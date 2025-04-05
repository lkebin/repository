package example

import (
	"database/sql/driver"
	"errors"
	"strings"

	"github.com/google/uuid"
)

type ID string

func (id ID) IsZero() bool {
	return id == ""
}

func (id ID) IsValid() bool {
	_, err := uuid.Parse(string(id))
	if err != nil {
		return false
	}
	return true
}

func (id *ID) Scan(value any) error {
	return convertToString(id, value)
}

func (id ID) Value() (driver.Value, error) {
	return convertToByte(id)
}

func (id ID) String() string {
	return string(id)
}

func (id ID) Short() string {
	return strings.ReplaceAll(string(id), "-", "")
}

func ParseIDFromShort(v string) (ID, error) {
	id, err := uuid.Parse(v)
	if err != nil {
		return "", err
	}

	return ID(id.String()), nil
}

type NullID struct {
	ID    ID
	Valid bool
}

func (n *NullID) Scan(value any) error {
	if value == nil {
		n.ID, n.Valid = "", false
		return nil
	}
	n.Valid = true
	return convertToString(&n.ID, value)
}

func (n NullID) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return convertToByte(n.ID)
}

func convertToString(id *ID, value any) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("invalid type for ID")
	}

	if len(b) != 16 {
		return errors.New("invalid data length for ID")
	}

	u := binToUuid(b)
	*id = ID(u.String())

	return nil
}

func convertToByte(id ID) ([]byte, error) {
	u, err := uuid.Parse(string(id))
	if err != nil {
		return nil, err
	}

	uu := uuidToBin(u)

	return []byte(uu), nil
}

func uuidToBin(u uuid.UUID) []byte {
	var b = make([]byte, 16)

	copy(b[0:], u[6:8])
	copy(b[2:], u[4:6])
	copy(b[4:], u[:4])
	copy(b[8:], u[8:])

	return b
}

func binToUuid(b []byte) uuid.UUID {
	var u uuid.UUID

	copy(u[0:], b[4:8])
	copy(u[4:], b[2:4])
	copy(u[6:], b[:2])
	copy(u[8:], b[8:])

	return u
}
