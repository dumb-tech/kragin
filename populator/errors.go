package populator

import "fmt"

const (
	errCodeWrongConfigStoreDataType = iota
	errCodeSpecifiedConfigIsMissing
	errCodeWrongInputDataType
	errCodeWrongInvalidStructType
	errCodeFieldIsPrivate
	errCodeRequiredValueIsMissing
)

type Error struct {
	code int
	text string
}

func (p Error) Error() string { return p.text }
func (p Error) Is(err error) bool {
	e, ok := err.(Error)
	return ok && (e.code == p.code)
}

func newPopulatorError(code int, text string) Error {
	return Error{
		code: code,
		text: text,
	}
}

var ErrWrongInputConfigStoreType = func(actual, wants any) error {
	return newPopulatorError(
		errCodeWrongConfigStoreDataType,
		fmt.Sprintf("wrong input config store data type. expected %T, got %T", wants, actual),
	)
}

var ErrSpecifiedConfigIsMissing = func(name string) error {
	return newPopulatorError(
		errCodeSpecifiedConfigIsMissing,
		fmt.Sprintf("config with name %s is missing in config store", name),
	)
}

var ErrWrongInputDataType = func(actual, wants any) error {
	return newPopulatorError(
		errCodeWrongInputDataType,
		fmt.Sprintf("wrong input data type. expected %T, got %T", wants, actual),
	)
}

var ErrInvalidTargetType = func(field string) Error {
	return newPopulatorError(
		errCodeWrongInvalidStructType,
		fmt.Sprintf("invalid target type for field %s", field))
}
var ErrFieldIsPrivate = func(field string) Error {
	return newPopulatorError(
		errCodeFieldIsPrivate,
		fmt.Sprintf("field %s is private", field))
}

var ErrRequiredValueIsMissing = func(field string) error {
	return newPopulatorError(
		errCodeRequiredValueIsMissing,
		fmt.Sprintf("field %s is required", field),
	)
}
