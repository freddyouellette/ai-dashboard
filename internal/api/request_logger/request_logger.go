package request_logger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/freddyouellette/ai-dashboard/plugins/plugin_models"
)

type Logger interface {
	Error(msg string, fields map[string]interface{})
	Warning(msg string, fields map[string]interface{})
	Info(msg string, fields map[string]interface{})
}

type RequestUtils interface {
	GetContextInt(r *http.Request, key any, def int) int
}

type RequestLogger struct {
	logger       Logger
	options      Options
	requestUtils RequestUtils
}

type Options struct {
	LogHeaders      bool
	LogRequestBody  bool
	LogResponseBody bool
	PrettyJson      bool
}

func NewRequestLogger(
	logger Logger,
	options Options,
	requestUtils RequestUtils,
) *RequestLogger {
	return &RequestLogger{
		logger:       logger,
		options:      options,
		requestUtils: requestUtils,
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
		info := map[string]interface{}{}
		info["user_id"] = l.requestUtils.GetContextInt(r, plugin_models.UserIdContextKey{}, 0)
		info["start"] = time.Now()
		info["method"] = r.Method
		info["url"] = r.URL.Path

		observer := &responseWriterObserver{ResponseWriter: w}
		next.ServeHTTP(observer, r)

		if l.options.LogHeaders {
			info["headers"] = r.Header
		}

		if l.options.LogRequestBody {
			body, err := io.ReadAll(r.Body)
			if err != nil {
				info["request_body"] = "<<<Error reading body>>>\n"
			} else if len(body) != 0 {
				info["request_body"] = sprintJson(body, l.options.PrettyJson)
			}
		}

		end := time.Now()
		info["end"] = end
		info["duration"] = time.Since(end).String()
		if l.options.LogResponseBody {
			info["response_body"] = sprintJson(observer.body, l.options.PrettyJson)
		}

		l.logger.Info("API HTTP Request", info)
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
