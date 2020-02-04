package logUtil

import (
	"log"
)

func FormatError(source string, er error) error {
	log.Printf("[%s]: %s", source, er)
	return er
}
