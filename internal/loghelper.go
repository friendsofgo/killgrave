package killgrave

import (
	"net/http"

	log "github.com/sirupsen/logrus"
)

func LogFieldsFromRequest(req *http.Request) log.Fields {
	return log.Fields{
		"method": req.Method,
		"url":    req.URL,
	}
}
