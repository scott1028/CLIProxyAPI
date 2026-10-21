package handlers

import (
	"bytes"
	"compress/gzip"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func gzipRequestBody(t *testing.T, payload []byte) []byte {
	t.Helper()
	var compressed bytes.Buffer
	writer := gzip.NewWriter(&compressed)
	if _, errWrite := writer.Write(payload); errWrite != nil {
		t.Fatalf("gzip write: %v", errWrite)
	}
	if errClose := writer.Close(); errClose != nil {
		t.Fatalf("gzip close: %v", errClose)
	}
	return compressed.Bytes()
}

func newGzipRequestContext(t *testing.T, payload []byte) *gin.Context {
	t.Helper()
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/chat/completions", bytes.NewReader(gzipRequestBody(t, payload)))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Request.Header.Set("Content-Encoding", "gzip")
	return c
}

func TestReadRequestBodyDecodesGzip(t *testing.T) {
	payload := []byte(`{"model":"test-model","messages":[{"role":"user","content":"hello"}]}`)
	c := newGzipRequestContext(t, payload)

	got, err := ReadRequestBody(c)
	if err != nil {
		t.Fatalf("ReadRequestBody() error = %v", err)
	}
	if !bytes.Equal(got, payload) {
		t.Fatalf("ReadRequestBody() = %s, want %s", got, payload)
	}
}

func TestReadRequestBodyRejectsOversizedGzip(t *testing.T) {
	payload := []byte(`{"input":"` + strings.Repeat("a", 256) + `"}`)
	c := newGzipRequestContext(t, payload)
	c.Set(MaxRequestBodyBytesContextKey, int64(64))

	_, err := ReadRequestBody(c)
	var maxBytesErr *http.MaxBytesError
	if !errors.As(err, &maxBytesErr) {
		t.Fatalf("ReadRequestBody() error = %v, want *http.MaxBytesError", err)
	}
	if maxBytesErr.Limit != 64 {
		t.Fatalf("MaxBytesError.Limit = %d, want 64", maxBytesErr.Limit)
	}
}
