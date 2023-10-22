package request_logger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type RequestLogger struct {
	logger  *log.Logger
	options Options
}

type Options struct {
	LogHeaders      bool
	LogRequestBody  bool
	LogResponseBody bool
	PrettyJson      bool
}

func NewRequestLogger(
	logger *log.Logger,
	options Options,
) *RequestLogger {
	return &RequestLogger{
		logger:  logger,
		options: options,
	}
}

type responseWriterObserver struct {
	http.ResponseWriter
	body       []byte
	statusCode int
}

func (wo *responseWriterObserver) Write(b []byte) (int, error) {
	wo.body = append(wo.body, b...)
	return wo.ResponseWriter.Write(b)
}

func (wo *responseWriterObserver) WriteHeader(statusCode int) {
	wo.statusCode = statusCode
	wo.ResponseWriter.WriteHeader(statusCode)
}

func (l *RequestLogger) CreateRequestLoggerHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		observer := &responseWriterObserver{ResponseWriter: w}
		next.ServeHTTP(observer, r)

		message := ""

		message += fmt.Sprintf("API REQ %s %s\n", r.Method, r.URL.Path)

		if l.options.LogHeaders {
			for name, values := range r.Header {
				for _, value := range values {
					message += fmt.Sprintf("%s: %s\n", name, value)
				}
			}
		}

		if l.options.LogRequestBody {
			body, err := io.ReadAll(r.Body)
			if err != nil {
				message += "<<<Error reading body>>>\n"
			} else if len(body) != 0 {
				message += sprintJson(body, l.options.PrettyJson)
			}
		}

		message += fmt.Sprintf("API RES %d %s\n", observer.statusCode, time.Since(start).String())
		if l.options.LogResponseBody {
			message += sprintJson(observer.body, l.options.PrettyJson)
		}

		l.logger.Print(message)
	})
}

func sprintJson(b []byte, pretty bool) string {
	if len(b) != 0 {
		if pretty {
			var prettyJSON bytes.Buffer
			err := json.Indent(&prettyJSON, b, "", "\t")
			if err == nil {
				return fmt.Sprintf("%s\n", prettyJSON.String())
			} else {
				return fmt.Sprintf("%s\n", string(b))
			}
		} else {
			return fmt.Sprintf("%s\n", string(b))
		}
	}
	return ""
}
