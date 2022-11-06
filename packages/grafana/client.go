package grafana

import (
	"net/http"
	"time"
)

type transport struct {
	token string
}

func (t *transport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Add("Authorization", "Bearer "+t.token)
	return http.DefaultTransport.RoundTrip(req)
}

type client struct {
	url string
	cli http.Client
}

func NewClient(url string, token string) *client {
	return &client{
		url: url,
		cli: http.Client{
			Transport: &transport{token: token},
			Timeout:   30 * time.Second,
		},
	}
}
