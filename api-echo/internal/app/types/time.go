package types

import (
	"time"

	"filmogophery/internal/pkg/constant"
)

func ConvertTime2Date(target time.Time) constant.Date {
	return constant.Date(target.Format(constant.DateFormat))
}
