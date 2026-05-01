package constant

import (
	"database/sql/driver"
	"fmt"
)

type PointAction string

const (
	PointActionWatchHistory PointAction = "watch_history"
	PointActionReview       PointAction = "review"
)

func (s *PointAction) Scan(value interface{}) error {
	switch val := value.(type) {
	case []uint8:
		*s = PointAction(val)
	case string:
		*s = PointAction(val)
	default:
		return fmt.Errorf("unsupported type for PointAction: %T", value)
	}
	return nil
}

func (s PointAction) Value() (driver.Value, error) {
	return string(s), nil
}
