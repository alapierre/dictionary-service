package util

type ZeroRowsAffectedError struct {
	Err error
}

func (e ZeroRowsAffectedError) Error() string {
	return "database update error: " + e.Err.Error()
}

func NewZeroRowsAffectedError(err error) ZeroRowsAffectedError {
	return ZeroRowsAffectedError{Err: err}
}
