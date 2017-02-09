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
	proxy.ServeHTTP(w, req)
}
