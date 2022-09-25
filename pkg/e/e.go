package e

import "fmt"

// Обертка ошибок с кастомной ошибкой в качестве пояснения.
func Wrap(msg string, err error) error {
	return fmt.Errorf("%s %w", msg, err)
}
