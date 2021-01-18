package gopher_query

import (
	"net"
	"time"
)

type Client struct {
	dialer *net.Dialer
}

func NewClient() Client {
	return Client{dialer: &net.Dialer{Timeout: 0}}
}

func NewClientWithTimeOut(timeout time.Duration) Client {
	return Client{dialer: &net.Dialer{Timeout: timeout}}
}