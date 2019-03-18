package httpclientmock

import (
	"context"
	"net/http"
	"time"

	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/env"
	"github.com/getfider/fider/app/services/httpclient"
)

func init() {
	//Increase transport timeouts when running Tests
	if env.IsTest() {
		transport := http.DefaultTransport.(*http.Transport)
		transport.TLSHandshakeTimeout = 30 * time.Second
	}
}

type Service struct{}

func (s Service) Category() string {
	return "httpclient"
}

var RequestsHistory = make([]*http.Request, 0)

func (s Service) Enabled() bool {
	return env.IsTest()
}

func (s Service) Init() {
	RequestsHistory = make([]*http.Request, 0)
	bus.AddHandler(s, requestHandler)
}

func requestHandler(ctx context.Context, cmd *httpclient.Request) error {
	req, err := http.NewRequest(cmd.Method, cmd.URL, cmd.Body)
	if err != nil {
		return err
	}

	for k, v := range cmd.Headers {
		req.Header.Set(k, v)
	}
	if cmd.BasicAuth != nil {
		req.SetBasicAuth(cmd.BasicAuth.User, cmd.BasicAuth.Password)
	}

	RequestsHistory = append(RequestsHistory, req)

	cmd.ResponseStatusCode = http.StatusOK
	cmd.ResponseBody = []byte("")
	return nil
}