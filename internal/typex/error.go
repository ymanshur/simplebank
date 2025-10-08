package typex

import "fmt"

type ErrUnProcessableEnity string

func (e ErrUnProcessableEnity) Error() string {
	return string(e)
}

type ErrDataNotFound string

func (e ErrDataNotFound) Error() string {
	return string(e)
}

type ErrUnAuthorized string

func (e ErrUnAuthorized) Error() string {
	return string(e)
}

type ErrForbidden string

func (e ErrForbidden) Error() string {
	return string(e)
}

func NewErrDataNotFound(data string) ErrDataNotFound {
	return ErrDataNotFound(fmt.Sprintf("%s not found", data))
}
