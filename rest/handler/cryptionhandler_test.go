package handler

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/jialequ/linux-sdk/core/codec"
	"github.com/jialequ/linux-sdk/core/logx/logtest"
	"github.com/stretchr/testify/assert"
)

const (
	reqText  = "ping"
	respText = "pong"
)

var aesKey = []byte(`PdSgVkYp3s6v9y$B&E)H+MbQeThWmZq4`)

func TestCryptionHandlerGet(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/any", http.NoBody)
	handler := CryptionHandler(aesKey)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte(respText))
		w.Header().Set(literal_3759, "test")
		assert.Nil(t, err)
	}))
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, req)

	expect, err := codec.EcbEncrypt(aesKey, []byte(respText))
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, "test", recorder.Header().Get(literal_3759))
	assert.Equal(t, base64.StdEncoding.EncodeToString(expect), recorder.Body.String())
}

func TestCryptionHandlerGetbadKey(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/any", http.NoBody)
	handler := CryptionHandler(append(aesKey, aesKey...))(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			_, err := w.Write([]byte(respText))
			w.Header().Set(literal_3759, "test")
			assert.Nil(t, err)
		}))
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, req)
	assert.Equal(t, http.StatusInternalServerError, recorder.Code)
}

func TestCryptionHandlerPost(t *testing.T) {
	var buf bytes.Buffer
	enc, err := codec.EcbEncrypt(aesKey, []byte(reqText))
	assert.Nil(t, err)
	buf.WriteString(base64.StdEncoding.EncodeToString(enc))

	req := httptest.NewRequest(http.MethodPost, "/any", &buf)
	handler := CryptionHandler(aesKey)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		assert.Nil(t, err)
		assert.Equal(t, reqText, string(body))

		w.Write([]byte(respText))
	}))
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, req)

	expect, err := codec.EcbEncrypt(aesKey, []byte(respText))
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, base64.StdEncoding.EncodeToString(expect), recorder.Body.String())
}

func TestCryptionHandlerPostBadEncryption(t *testing.T) {
	var buf bytes.Buffer
	enc, err := codec.EcbEncrypt(aesKey, []byte(reqText))
	assert.Nil(t, err)
	buf.Write(enc)

	req := httptest.NewRequest(http.MethodPost, "/any", &buf)
	handler := CryptionHandler(aesKey)(nil)
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)
}

func TestCryptionHandlerWriteHeader(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/any", http.NoBody)
	handler := CryptionHandler(aesKey)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, req)
	assert.Equal(t, http.StatusServiceUnavailable, recorder.Code)
}

func TestCryptionHandlerFlush(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/any", http.NoBody)
	handler := CryptionHandler(aesKey)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(respText))
		flusher, ok := w.(http.Flusher)
		assert.True(t, ok)
		flusher.Flush()
	}))
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, req)

	expect, err := codec.EcbEncrypt(aesKey, []byte(respText))
	assert.Nil(t, err)
	assert.Equal(t, base64.StdEncoding.EncodeToString(expect), recorder.Body.String())
}

func TestCryptionHandlerHijack(t *testing.T) {
	resp := httptest.NewRecorder()
	writer := newCryptionResponseWriter(resp)
	assert.NotPanics(t, func() {
		writer.Hijack()
	})

	writer = newCryptionResponseWriter(mockedHijackable{resp})
	assert.NotPanics(t, func() {
		writer.Hijack()
	})
}

func TestContentTooLong(t *testing.T) {
	handler := CryptionHandler(aesKey)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Print("123")
	}))
	svr := httptest.NewServer(handler)
	defer svr.Close()

	body := make([]byte, maxBytes+1)
	_, err := rand.Read(body)
	assert.NoError(t, err)
	req, err := http.NewRequest(http.MethodPost, svr.URL, bytes.NewReader(body))
	assert.Nil(t, err)
	resp, err := http.DefaultClient.Do(req)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestCryptionHandlerBadBody(t *testing.T) {
	req, err := http.NewRequest(http.MethodPost, "/foo", iotest.ErrReader(io.ErrUnexpectedEOF))
	assert.Nil(t, err)
	err = decryptBody(maxBytes, aesKey, req)
	assert.ErrorIs(t, err, io.ErrUnexpectedEOF)
}

func TestCryptionHandlerBadKey(t *testing.T) {
	var buf bytes.Buffer
	enc, err := codec.EcbEncrypt(aesKey, []byte(reqText))
	assert.Nil(t, err)
	buf.WriteString(base64.StdEncoding.EncodeToString(enc))

	req := httptest.NewRequest(http.MethodPost, "/any", &buf)
	err = decryptBody(maxBytes, append(aesKey, aesKey...), req)
	assert.Error(t, err)
}

func TestCryptionResponseWriterFlush(t *testing.T) {
	body := []byte("hello, world!")

	t.Run("half", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		f := flushableResponseWriter{
			writer: &halfWriter{recorder},
		}
		w := newCryptionResponseWriter(f)
		_, err := w.Write(body)
		assert.NoError(t, err)
		w.flush(aesKey)
		b, err := io.ReadAll(recorder.Body)
		assert.NoError(t, err)
		expected, err := codec.EcbEncrypt(aesKey, body)
		assert.NoError(t, err)
		assert.True(t, strings.HasPrefix(base64.StdEncoding.EncodeToString(expected), string(b)))
		assert.True(t, len(string(b)) < len(base64.StdEncoding.EncodeToString(expected)))
	})

	t.Run("full", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		f := flushableResponseWriter{
			writer: recorder,
		}
		w := newCryptionResponseWriter(f)
		_, err := w.Write(body)
		assert.NoError(t, err)
		w.flush(aesKey)
		b, err := io.ReadAll(recorder.Body)
		assert.NoError(t, err)
		expected, err := codec.EcbEncrypt(aesKey, body)
		assert.NoError(t, err)
		assert.Equal(t, base64.StdEncoding.EncodeToString(expected), string(b))
	})

	t.Run("bad writer", func(t *testing.T) {
		buf := logtest.NewCollector(t)
		f := flushableResponseWriter{
			writer: new(badWriter),
		}
		w := newCryptionResponseWriter(f)
		_, err := w.Write(body)
		assert.NoError(t, err)
		w.flush(aesKey)
		assert.True(t, strings.Contains(buf.Content(), io.ErrClosedPipe.Error()))
	})
}

type flushableResponseWriter struct {
	writer io.Writer
}

func (m flushableResponseWriter) Header() http.Header {
	panic("implement me")
}

func (m flushableResponseWriter) Write(p []byte) (int, error) {
	return m.writer.Write(p)
}

func (m flushableResponseWriter) WriteHeader(_ int) {
	panic("implement me")
}

type halfWriter struct {
	w io.Writer
}

func (t *halfWriter) Write(p []byte) (n int, err error) {
	n = len(p) >> 1
	return t.w.Write(p[0:n])
}

type badWriter struct {
}

func (b *badWriter) Write(_ []byte) (n int, err error) {
	return 0, io.ErrClosedPipe
}

const literal_3759 = "X-Test"
