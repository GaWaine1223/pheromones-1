package pheromone

const (
	ErrLocalSocketTimeout = 1001
)

type Error int

func (err Error) Error() string {
	return errMap[err]
}

var errMap = map[Error]string{
	ErrLocalSocketTimeout : "read timeout",
}