package logged_client

import (
	"fmt"
	"log"
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

type LoggedClient struct {
	*http.Client
	logger  *log.Logger
	options Options
}

func NewLoggedClient(client *http.Client, logger *log.Logger, options Options) *LoggedClient {
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
	message := fmt.Sprintf("EXT REQ: %s\n", string(reqBytes))

	resp, err := c.Client.Do(req)

	if err != nil {
		message += fmt.Sprintf("CLIENT ERR: %s\n", err.Error())
	} else {
		resBytes, err := httputil.DumpResponse(resp, true)
		if err != nil {
			return nil, err
		}
		message += fmt.Sprintf("EXT RESP: %s\n", string(resBytes))
	}

	c.logger.Print(message)
	return resp, err
}
