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

type populatorError struct {
	code int
	text string
}

func (p populatorError) Error() string { return p.text }
func (p populatorError) Is(err error) bool {
	e, ok := err.(populatorError)
	return ok && (e.code == p.code)
}

func newPopulatorError(code int, text string) populatorError {
	return populatorError{
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

var ErrInvalidTargetType = func(field string) error {
	return newPopulatorError(
		errCodeWrongInvalidStructType,
		fmt.Sprintf("invalid target type for field %s", field))
}
var ErrFieldIsPrivate = func(field string) error {
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
