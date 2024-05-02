package mdlware

import (
	"log"
	"net/http"
	"net/http/httputil"
	"time"

	"homework/internal/pkg/kafkalogger"
)

type logger interface {
	LogMessage(message kafkalogger.Message) error
}

// Logger is a middleware that log request to specified logger
func Logger(s logger, handler http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		rawRequest, err := httputil.DumpRequest(req, true)
		if err != nil {
			log.Println("httputil.DumpRequest:", err)
		}

		err = s.LogMessage(kafkalogger.Message{
			RequestTime: time.Now(),
			HTTPMethod:  req.Method,
			RawRequest:  string(rawRequest),
		})
		if err != nil {
			log.Println("LogMessage:", err)
		}

		handler.ServeHTTP(w, req)
	}
}
