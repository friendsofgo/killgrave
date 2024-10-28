package killgrave

import (
	"net/http"

	log "github.com/sirupsen/logrus"
)

const (
	methodField string = "Method"
	urlField           = "URL"
)

func LogFieldsFromRequest(req *http.Request) log.Fields {
	return log.Fields{
		methodField: req.Method,
		urlField:    req.URL,
	}
}
