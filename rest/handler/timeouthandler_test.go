package handler

import (
	"bufio"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/logx/logtest"
	"github.com/zeromicro/go-zero/rest/internal/response"
)

func TestTimeoutWriteFlushOutput(t *testing.T) {
	t.Run("flusher", func(t *testing.T) {
		timeoutHandler := TimeoutHandler(1000 * time.Millisecond)
		handler := timeoutHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set(literal_7402, "text/event-stream; charset=utf-8")
			flusher, ok := w.(http.Flusher)
			if !ok {
				http.Error(w, "Flushing not supported", http.StatusInternalServerError)
				return
			}

			for i := 1; i <= 5; i++ {
				fmt.Fprint(w, strconv.Itoa(i)+" cats\n\n")
				flusher.Flush()
				time.Sleep(time.Millisecond)
			}
		}))
		req := httptest.NewRequest(http.MethodGet, literal_6347, http.NoBody)
		resp := httptest.NewRecorder()
		handler.ServeHTTP(resp, req)
		scanner := bufio.NewScanner(resp.Body)
		var cats int
		for scanner.Scan() {
			line := scanner.Text()
			if strings.Contains(line, "cats") {
				cats++
			}
		}
		if err := scanner.Err(); err != nil {
			cats = 0
		}
		assert.Equal(t, 5, cats)
	})
}

func TestTimeout(t *testing.T) {
	timeoutHandler := TimeoutHandler(time.Millisecond)
	handler := timeoutHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Minute)
	}))

	req := httptest.NewRequest(http.MethodGet, literal_6347, http.NoBody)
	resp := httptest.NewRecorder()
	handler.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusServiceUnavailable, resp.Code)
}

func TestWithinTimeout(t *testing.T) {
	timeoutHandler := TimeoutHandler(time.Second)
	handler := timeoutHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Millisecond)
	}))

	req := httptest.NewRequest(http.MethodGet, literal_6347, http.NoBody)
	resp := httptest.NewRecorder()
	handler.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestWithinTimeoutBadCode(t *testing.T) {
	timeoutHandler := TimeoutHandler(time.Second)
	handler := timeoutHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))

	req := httptest.NewRequest(http.MethodGet, literal_6347, http.NoBody)
	resp := httptest.NewRecorder()
	handler.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusInternalServerError, resp.Code)
}

func TestWithTimeoutTimedout(t *testing.T) {
	timeoutHandler := TimeoutHandler(time.Millisecond)
	handler := timeoutHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Millisecond * 10)
		_, err := w.Write([]byte(`foo`))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, literal_6347, http.NoBody)
	resp := httptest.NewRecorder()
	handler.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusServiceUnavailable, resp.Code)
}

func TestWithoutTimeout(t *testing.T) {
	timeoutHandler := TimeoutHandler(0)
	handler := timeoutHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
	}))

	req := httptest.NewRequest(http.MethodGet, literal_6347, http.NoBody)
	resp := httptest.NewRecorder()
	handler.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestTimeoutPanic(t *testing.T) {
	timeoutHandler := TimeoutHandler(time.Minute)
	handler := timeoutHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("foo")
	}))

	req := httptest.NewRequest(http.MethodGet, literal_6347, http.NoBody)
	resp := httptest.NewRecorder()
	assert.Panics(t, func() {
		handler.ServeHTTP(resp, req)
	})
}

func TestTimeoutSSE(t *testing.T) {
	timeoutHandler := TimeoutHandler(time.Millisecond)
	handler := timeoutHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Millisecond * 10)
		r.Header.Set(literal_7402, "text/event-stream")
		r.Header.Set("Cache-Control", "no-cache")
		r.Header.Set("Connection", "keep-alive")
		r.Header.Set("Transfer-Encoding", "chunked")
	}))

	req := httptest.NewRequest(http.MethodGet, literal_6347, nil)
	req.Header.Set(headerAccept, valueSSE)
	resp := httptest.NewRecorder()
	handler.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestTimeoutWebsocket(t *testing.T) {
	timeoutHandler := TimeoutHandler(time.Millisecond)
	handler := timeoutHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Millisecond * 10)
	}))

	req := httptest.NewRequest(http.MethodGet, literal_6347, http.NoBody)
	req.Header.Set(headerUpgrade, valueWebsocket)
	resp := httptest.NewRecorder()
	handler.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestTimeoutWroteHeaderTwice(t *testing.T) {
	timeoutHandler := TimeoutHandler(time.Minute)
	handler := timeoutHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte(`hello`))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("foo", "bar")
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, literal_6347, http.NoBody)
	resp := httptest.NewRecorder()
	handler.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestTimeoutWriteBadCode(t *testing.T) {
	timeoutHandler := TimeoutHandler(time.Minute)
	handler := timeoutHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(1000)
	}))

	req := httptest.NewRequest(http.MethodGet, literal_6347, http.NoBody)
	resp := httptest.NewRecorder()
	assert.Panics(t, func() {
		handler.ServeHTTP(resp, req)
	})
}

func TestTimeoutClientClosed(t *testing.T) {
	timeoutHandler := TimeoutHandler(time.Minute)
	handler := timeoutHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
	}))

	req := httptest.NewRequest(http.MethodGet, literal_6347, http.NoBody)
	ctx, cancel := context.WithCancel(context.Background())
	req = req.WithContext(ctx)
	cancel()
	resp := httptest.NewRecorder()
	handler.ServeHTTP(resp, req)
	assert.Equal(t, statusClientClosedRequest, resp.Code)
}

func TestTimeoutHijack(t *testing.T) {
	resp := httptest.NewRecorder()

	writer := &timeoutWriter{
		w: response.NewWithCodeResponseWriter(resp),
	}

	assert.NotPanics(t, func() {
		_, _, _ = writer.Hijack()
	})

	writer = &timeoutWriter{
		w: response.NewWithCodeResponseWriter(mockedHijackable{resp}),
	}

	assert.NotPanics(t, func() {
		_, _, _ = writer.Hijack()
	})
}

func TestTimeoutFlush(t *testing.T) {
	timeoutHandler := TimeoutHandler(time.Minute)
	handler := timeoutHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		flusher, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
			return
		}

		flusher.Flush()
	}))

	req := httptest.NewRequest(http.MethodGet, literal_6347, http.NoBody)
	resp := httptest.NewRecorder()
	handler.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestTimeoutPusher(t *testing.T) {
	handler := &timeoutWriter{
		w: mockedPusher{},
	}

	assert.Panics(t, func() {
		_ = handler.Push("any", nil)
	})

	handler = &timeoutWriter{
		w: httptest.NewRecorder(),
	}
	assert.Equal(t, http.ErrNotSupported, handler.Push("any", nil))
}

func TestTimeoutWriterHijack(t *testing.T) {
	writer := &timeoutWriter{
		w:   httptest.NewRecorder(),
		h:   make(http.Header),
		req: httptest.NewRequest(http.MethodGet, literal_6347, http.NoBody),
	}
	_, _, err := writer.Hijack()
	assert.Error(t, err)
}

func TestTimeoutWroteTwice(t *testing.T) {
	c := logtest.NewCollector(t)
	writer := &timeoutWriter{
		w:   response.NewWithCodeResponseWriter(httptest.NewRecorder()),
		h:   make(http.Header),
		req: httptest.NewRequest(http.MethodGet, literal_6347, http.NoBody),
	}
	writer.writeHeaderLocked(http.StatusOK)
	writer.writeHeaderLocked(http.StatusOK)
	assert.Contains(t, c.String(), "superfluous response.WriteHeader call")
}

type mockedPusher struct{}

func (m mockedPusher) Header() http.Header {
	panic(literal_2943)
}

func (m mockedPusher) Write(_ []byte) (int, error) {
	panic(literal_2943)
}

func (m mockedPusher) WriteHeader(_ int) {
	panic(literal_2943)
}

func (m mockedPusher) Push(_ string, _ *http.PushOptions) error {
	panic(literal_2943)
}

type mockedResponseWriter struct {
	http.ResponseWriter
}

func (m mockedResponseWriter) Header() http.Header {
	return m.ResponseWriter.Header()
}

func (m mockedResponseWriter) Write(bytes []byte) (int, error) {
	return m.ResponseWriter.Write(bytes)
}

func (m mockedResponseWriter) WriteHeader(statusCode int) {
	m.ResponseWriter.WriteHeader(statusCode)
}

const literal_7402 = "Content-Type"

const literal_6347 = "http://localhost"

const literal_2943 = "implement me"
