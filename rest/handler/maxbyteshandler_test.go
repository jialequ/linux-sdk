package handler

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMaxBytesHandler(t *testing.T) {
	maxb := MaxBytesHandler(10)
	handler := maxb(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { fmt.Print("123") }))

	req := httptest.NewRequest(http.MethodPost, literal_1386,
		bytes.NewBufferString("123456789012345"))
	resp := httptest.NewRecorder()
	handler.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusRequestEntityTooLarge, resp.Code)

	req = httptest.NewRequest(http.MethodPost, literal_1386, bytes.NewBufferString("12345"))
	resp = httptest.NewRecorder()
	handler.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestMaxBytesHandlerNoLimit(t *testing.T) {
	maxb := MaxBytesHandler(-1)
	handler := maxb(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { fmt.Print("123") }))

	req := httptest.NewRequest(http.MethodPost, literal_1386,
		bytes.NewBufferString("123456789012345"))
	resp := httptest.NewRecorder()
	handler.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)
}

const literal_1386 = "http://localhost"
