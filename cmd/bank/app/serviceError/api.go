package serviceError

import "fmt"

type ApiError struct {
	ID   int64
	Code int64
}

func (a ApiError) Error() string {
	return fmt.Sprintf("api serviceError: id - %d, code - %d", a.ID, a.Code)
}

var ErrOriginal = &ApiError{
	ID:   1,
	Code: 201,
}

type ServiceError struct {
	Err    error
	Params interface{}
}

func (s ServiceError) Error() string {
	return fmt.Sprintf("call with %v params got %v", s.Params, s.Err)
}
func (s ServiceError) Unwrap() error {
	return s.Err
}
