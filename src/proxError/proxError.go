package proxError

import (
	"errors"
	"log"
)

func PanicIfErr(err error, logger *log.Logger) {

	if err != nil {
		logger.Panicln(err)
	}
}

var (
	Unknown = errors.New("UK")
)
