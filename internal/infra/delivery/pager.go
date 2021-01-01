package delivery

import (
	"foodmap/internal/infra/errors"
	"strconv"
)

// ParseLimitAndSkip used in requests with array response
func ParseLimitAndSkip(limit, skip string) (int64, int64, error) {
	var (
		l, s int64
		err  error
	)
	if limit == "" {
		l = 0
	} else {
		l, err = strconv.ParseInt(limit, 10, 64)
	}
	if err != nil {
		return 0, 0, errors.NewValidationError("invalid", "limit")
	}
	if skip == "" {
		l = 0
	} else {
		s, err = strconv.ParseInt(skip, 10, 64)
	}
	if err != nil {
		return 0, 0, errors.NewValidationError("invalid", "skip")
	}
	return l, s, nil
}
