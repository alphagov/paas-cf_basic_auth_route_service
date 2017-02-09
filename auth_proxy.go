package main

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

const (
	CF_FORWARDED_URL_HEADER = "X-CF-Forwarded-Url"
)

type AuthProxy struct {
	username string
	password string
}

func NewAuthProxy(username, password string) http.Handler {
	return &AuthProxy{
		username: username,
		password: password,
	}
}

func (a *AuthProxy) checkAuth(user, pass string) bool {
	return user == a.username && pass == a.password
}

func (a *AuthProxy) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	user, pass, ok := req.BasicAuth()
	if !ok || !a.checkAuth(user, pass) {
		http.Error(w, "Unauthorized.", http.StatusUnauthorized)
		return
	}

	forwardedURL := req.Header.Get(CF_FORWARDED_URL_HEADER)
	if forwardedURL == "" {
		http.Error(w, "Missing Forwarded URL", http.StatusBadRequest)
	}
	url, err := url.Parse(forwardedURL)
	if err != nil {
		http.Error(w, "Invalid forward URL: "+err.Error(), http.StatusBadRequest)
		return
	}

	proxy := httputil.NewSingleHostReverseProxy(url)

	defaultDirector := proxy.Director
	proxy.Director = func(beReq *http.Request) {
		defaultDirector(beReq)

		// Set the Host header to match the forwarded hostname instead of the one from the incoming request.
		req.Host = url.Host

		// Setting a blank User-Agent causes the http lib not to output one, whereas if there
		// is no header, it will output a default one.
		// See: https://github.com/golang/go/blob/release-branch.go1.7/src/net/http/request.go#L503
		if _, present := req.Header["User-Agent"]; !present {
			req.Header.Set("User-Agent", "")
		}
	}

	proxy.ServeHTTP(w, req)
}
