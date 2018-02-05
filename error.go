// Copyright 2018 Lothar . All rights reserved.
// https://github.com/GaWaine1223

package pheromone

const (
	ErrLocalSocketTimeout = 1001

	ErrRemoteSocketEmpty  = 2001
	ErrRemoteSocketExist  = 2002
)

type Error int

func (err Error) Error() string {
	return errMap[err]
}

var errMap = map[Error]string{
	ErrLocalSocketTimeout : "read timeout",
}