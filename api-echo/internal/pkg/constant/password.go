package constant

import (
	"database/sql/driver"
	"fmt"
)

type (
	PasswordHasher string
)

// --- Scanner --- //
// db の値を models の型に合わせて変換する

func (s *PasswordHasher) Scan(value interface{}) error {
	switch val := value.(type) {
	case []uint8:
		*s = PasswordHasher(val)
	case string:
		*s = PasswordHasher(val)
	default:
		return fmt.Errorf("unsupported type for PasswordHasher: %T", value)
	}
	return nil
}

// --- Valuer --- //
// models の型を db の型に合わせて変換する

func (s PasswordHasher) Value() (driver.Value, error) {
	return string(s), nil
}
