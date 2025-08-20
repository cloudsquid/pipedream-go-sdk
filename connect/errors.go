package connect

import "errors"

var (
	NotFoundErr error = errors.New("requested resource does not exist")
)
