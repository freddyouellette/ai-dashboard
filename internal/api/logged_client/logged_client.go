package logged_client

import (
	"net/http"
	"net/http/httputil"
)

type Options struct {
	LogRequestHeaders  bool
	LogRequestBody     bool
	LogResponseHeaders bool
	LogResponseBody    bool
	PrettyJson         bool
}

type Logger interface {
	Error(msg string, fields map[string]interface{})
	Warning(msg string, fields map[string]interface{})
	Info(msg string, fields map[string]interface{})
}

type LoggedClient struct {
	*http.Client
	logger  Logger
	options Options
}

func NewLoggedClient(client *http.Client, logger Logger, options Options) *LoggedClient {
	return &LoggedClient{
		Client:  client,
		logger:  logger,
		options: options,
	}
}

func (c *LoggedClient) Do(req *http.Request) (*http.Response, error) {
	reqBytes, err := httputil.DumpRequest(req, true)
	if err != nil {
		return nil, err
	}
	resp, err := c.Client.Do(req)

	if err != nil {
		c.logger.Error("Error making HTTP client request", map[string]interface{}{
			"request": string(reqBytes),
			"error":   err.Error(),
		})
	} else {
		resBytes, err := httputil.DumpResponse(resp, true)
		if err != nil {
			return nil, err
		}
		c.logger.Info("HTTP client Request", map[string]interface{}{
			"request":  string(reqBytes),
			"response": string(resBytes),
			"error":    err.Error(),
		})
	}

	return resp, err
}
