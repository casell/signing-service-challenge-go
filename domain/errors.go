package domain

import "fmt"

type ErrInvalidAlgorithm struct {
	algorithm string
}

func (e ErrInvalidAlgorithm) Error() string {
	return fmt.Sprintf("invalid algorithm %s", e.algorithm)
}
