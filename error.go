package backoff

// Error is the backoff error.
type Error struct {
	Inner error  // Inner is the inner error.
	Type  string // Type is the error type.
	Alg   string // Alg is the backoff algorithm.
	Msg   string // Msg is the error message.
}

func (e *Error) Unwrap() error {
	return e.Inner
}

func (e *Error) Error() string {
	s := "go-backoff/backoff: " + e.Type + ":"
	if e.Msg != "" {
		s += " " + e.Msg
	}
	if e.Inner != nil {
		s = s + " [" + e.Inner.Error() + "]"
	}
	return s
}

func (e *Error) Is(target error) bool {
	ee, ok := target.(*Error)
	if ok {
		return e.Type == ee.Type
	}
	return false
}

func rangeError(msg string) *Error {
	return &Error{
		Type: "range",
		Msg:  msg,
	}
}
