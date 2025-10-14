package util

import "errors"

func JoinErrors(errs ...error) error {
	return errors.Join(errs...)
}
