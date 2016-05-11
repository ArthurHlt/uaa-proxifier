package main

import (
	"net/url"
	"net/http/httputil"
	"net/http"
	"crypto/tls"
)

type Prox struct {
	target            *url.URL
	proxy             *httputil.ReverseProxy
	defaultTransport  http.RoundTripper
	transportUserInfo http.RoundTripper
}

func NewCustomProxy(target string, skipInsecure bool) *Prox {
	url, err := url.Parse(target)
	fatalIf("Parsing uaa url", err)
	defTr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: skipInsecure},
		Proxy: http.ProxyFromEnvironment,
	}
	alterTr := &TransportUserInfo{
		RoundTripper: defTr,
	}
	proxy := httputil.NewSingleHostReverseProxy(url)
	proxy.Transport = defTr
	return &Prox{target: url, proxy: proxy, defaultTransport: defTr, transportUserInfo: alterTr}
}

func (p *Prox) handle(w http.ResponseWriter, r *http.Request) {
	r.Host = p.target.Host
	r.Header.Del("X-Forwarded-Proto")
	if r.URL.Path != "/userinfo" {
		p.proxy.Transport = p.defaultTransport
		p.proxy.ServeHTTP(w, r)
		return
	}
	p.proxy.Transport = p.transportUserInfo
	p.proxy.ServeHTTP(w, r)
}


