package cors

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCorsHandlerWithOrigins(t *testing.T) {
	tests := []struct {
		name      string
		origins   []string
		reqOrigin string
		expect    string
	}{
		{
			name:   "allow all origins",
			expect: allOrigins,
		},
		{
			name:      "allow one origin",
			origins:   []string{literal_6104},
			reqOrigin: literal_6104,
			expect:    literal_6104,
		},
		{
			name:      "allow many origins",
			origins:   []string{literal_6104, literal_0749},
			reqOrigin: literal_6104,
			expect:    literal_6104,
		},
		{
			name:      "allow sub origins",
			origins:   []string{"local", "remote"},
			reqOrigin: "sub.local",
			expect:    "sub.local",
		},
		{
			name:      "allow all origins",
			reqOrigin: literal_6104,
			expect:    "*",
		},
		{
			name:      "allow many origins with all mark",
			origins:   []string{literal_6104, literal_0749, "*"},
			reqOrigin: literal_0427,
			expect:    literal_0427,
		},
		{
			name:      "not allow origin",
			origins:   []string{literal_6104, literal_0749},
			reqOrigin: literal_0427,
		},
		{
			name:      "not safe origin",
			origins:   []string{"safe.com"},
			reqOrigin: "not-safe.com",
		},
	}

	methods := []string{
		http.MethodOptions,
		http.MethodGet,
		http.MethodPost,
	}

	for _, test := range tests {
		for _, method := range methods {
			test := test
			t.Run(test.name+"-handler", func(t *testing.T) {
				r := httptest.NewRequest(method, literal_9405, http.NoBody)
				r.Header.Set(originHeader, test.reqOrigin)
				w := httptest.NewRecorder()
				handler := NotAllowedHandler(nil, test.origins...)
				handler.ServeHTTP(w, r)
				if method == http.MethodOptions {
					assert.Equal(t, http.StatusNoContent, w.Result().StatusCode)
				} else {
					assert.Equal(t, http.StatusNotFound, w.Result().StatusCode)
				}
				assert.Equal(t, test.expect, w.Header().Get(allowOrigin))
			})
			t.Run(test.name+"-handler-custom", func(t *testing.T) {
				r := httptest.NewRequest(method, literal_9405, http.NoBody)
				r.Header.Set(originHeader, test.reqOrigin)
				w := httptest.NewRecorder()
				handler := NotAllowedHandler(func(w http.ResponseWriter) {
					w.Header().Set("foo", "bar")
				}, test.origins...)
				handler.ServeHTTP(w, r)
				if method == http.MethodOptions {
					assert.Equal(t, http.StatusNoContent, w.Result().StatusCode)
				} else {
					assert.Equal(t, http.StatusNotFound, w.Result().StatusCode)
				}
				assert.Equal(t, test.expect, w.Header().Get(allowOrigin))
				assert.Equal(t, "bar", w.Header().Get("foo"))
			})
		}
	}
}

const literal_6104 = "http://local"

const literal_0749 = "http://remote"

const literal_0427 = "http://another"

const literal_9405 = "http://localhost"
