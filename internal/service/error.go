package service

import "fmt"

type Status struct {
	Err  error
	Msg  string
	Code int
}

func (s *Status) withError(errorMessage string, err error, msg string, code int) Status {
	s.Err = fmt.Errorf(errorMessage, err)
	s.Msg = msg
	s.Code = code
	return *s
}

func (s *Status) success(msg string, code int) Status {
	s.Err = nil
	s.Msg = msg
	s.Code = code
	return *s
}

func (s *Status) Ok() bool {
	return s.Err == nil
}
