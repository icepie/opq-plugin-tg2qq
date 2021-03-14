package proxy

import (
	"crypto/tls"
	"net/http"
	"net/url"

	"golang.org/x/net/proxy"
)

func HttpClient(addr string, auth ...*proxy.Auth) (client *http.Client, err error) {
	proxy, err := url.Parse(addr)
	if err != nil {
		return
	}

	transport := &http.Transport{
		Proxy:           http.ProxyURL(proxy),
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client = &http.Client{Transport: transport}

	return
}
