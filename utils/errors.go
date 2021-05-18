package utils

import (
	"go.uber.org/zap"
)

func CheckErrors(err error, code, msg, solution string) {

	if err != nil {
		logger := zap.NewExample().Sugar()
		defer logger.Sync()
		logger.Infow(msg,
			"code", code,
			"solution", solution)
	} else {
		return
	}
}
