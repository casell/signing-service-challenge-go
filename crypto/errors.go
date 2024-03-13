package crypto

import "fmt"

type ErrUnknownAlgorithm struct {
	algorithm string
}

func (e ErrUnknownAlgorithm) Error() string {
	return fmt.Sprintf("algorithm: unknown signature algorithm %s", e.algorithm)
}
