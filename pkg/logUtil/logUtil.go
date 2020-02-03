package logUtil

import (
	"errors"
	"fmt"
)

func FormatError(source string, er error) error {
	return errors.New(fmt.Sprintf("[%s]: %s", source, er))
}
