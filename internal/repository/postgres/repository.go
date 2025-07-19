package postgres

import "errors"

var (
	ErrWrongUserAgent = errors.New("different user agent")
)
