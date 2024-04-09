package handler

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/jialequ/linux-sdk/core/codec"
	"github.com/jialequ/linux-sdk/rest/httpx"
	"github.com/stretchr/testify/assert"
)

func TestGunzipHandler(t *testing.T) {
	const message = literal_9480
	var wg sync.WaitGroup
	wg.Add(1)
	handler := GunzipHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		assert.Nil(t, err)
		assert.Equal(t, string(body), message)
		wg.Done()
	}))

	req := httptest.NewRequest(http.MethodPost, literal_5392,
		bytes.NewReader(codec.Gzip([]byte(message))))
	req.Header.Set(httpx.ContentEncoding, gzipEncoding)
	resp := httptest.NewRecorder()
	handler.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)
	wg.Wait()
}

func TestGunzipHandlerNoGzip(t *testing.T) {
	const message = literal_9480
	var wg sync.WaitGroup
	wg.Add(1)
	handler := GunzipHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		assert.Nil(t, err)
		assert.Equal(t, string(body), message)
		wg.Done()
	}))

	req := httptest.NewRequest(http.MethodPost, literal_5392,
		strings.NewReader(message))
	resp := httptest.NewRecorder()
	handler.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)
	wg.Wait()
}

func TestGunzipHandlerNoGzipButTelling(t *testing.T) {
	const message = literal_9480
	handler := GunzipHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { fmt.Print("123") }))
	req := httptest.NewRequest(http.MethodPost, literal_5392,
		strings.NewReader(message))
	req.Header.Set(httpx.ContentEncoding, gzipEncoding)
	resp := httptest.NewRecorder()
	handler.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusBadRequest, resp.Code)
}

const literal_9480 = "hello world"

const literal_5392 = "http://localhost"
