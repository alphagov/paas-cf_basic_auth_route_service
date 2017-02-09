package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
)

func TestAuthProxy(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "test suite")
}

var _ = Describe("Basic Auth proxy", func() {

	var (
		username string = "user"
		password string = "secret"
		proxy    http.Handler
		backend  *ghttp.Server
		req      *http.Request
	)

	BeforeEach(func() {
		proxy = NewAuthProxy(username, password)
		backend = ghttp.NewServer()
		backend.AllowUnhandledRequests = true
		backend.UnhandledRequestStatusCode = http.StatusOK
	})

	AfterEach(func() {
		backend.Close()
	})

	Context("with a request from route-services", func() {
		BeforeEach(func() {
			req = httptest.NewRequest("GET", "/", nil)
			req.Header.Set("X-CF-Forwarded-Url", backend.URL())
			req.Header.Set("X-CF-Proxy-Signature", "Stub signature")
			req.Header.Set("X-CF-Proxy-Metadata", "Stub metadata")
		})

		Context("with the correct username and password", func() {
			BeforeEach(func() {
				req.SetBasicAuth(username, password)
			})

			It("should proxy the request to the backend", func() {
				w := httptest.NewRecorder()
				proxy.ServeHTTP(w, req)

				Expect(w.Code).To(Equal(http.StatusOK))

				Expect(backend.ReceivedRequests()).To(HaveLen(1))

				headers := backend.ReceivedRequests()[0].Header
				Expect(headers.Get("X-CF-Proxy-Signature")).To(Equal("Stub signature"))
				Expect(headers.Get("X-CF-Proxy-Metadata")).To(Equal("Stub metadata"))
			})
		})

		Context("with invalid credentials", func() {
			BeforeEach(func() {
				req.SetBasicAuth(username, "not the password")
			})

			It("returns a 401 Unauthorized", func() {
				w := httptest.NewRecorder()
				proxy.ServeHTTP(w, req)

				Expect(w.Code).To(Equal(http.StatusUnauthorized))
			})

			It("does not make a request to the backend", func() {
				w := httptest.NewRecorder()
				proxy.ServeHTTP(w, req)

				Expect(backend.ReceivedRequests()).To(HaveLen(0))
			})
		})
	})
})
