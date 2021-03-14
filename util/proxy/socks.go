package proxy

import (
	"net"
	"net/http"
	"time"

	"golang.org/x/net/proxy"
)

func Socks5Client(addr string, auth ...*proxy.Auth) (client *http.Client, err error) {
	dialer, err := proxy.SOCKS5("tcp", addr,
		nil,
		&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		},
	)
	if err != nil {
		return
	}

	transport := &http.Transport{
		Proxy:               nil,
		Dial:                dialer.Dial,
		TLSHandshakeTimeout: 10 * time.Second,
	}

	client = &http.Client{Transport: transport}

	return
}
