package errors

import core "dappco.re/go"

func New(text string) error {
	return core.NewError(text)
}

func Is(err, target error) bool {
	return core.Is(err, target)
}

func As(err error, target any) bool {
	return core.As(err, target)
}

func Join(errs ...error) error {
	return core.ErrorJoin(errs...)
}
