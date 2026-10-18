package utils

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

var supportedContentTypes = []string{
	"application/javascript",
	"application/json",
	"text/css",
	"text/html",
	"text/plain",
	"text/xml",
}

type compressWriter struct {
	w  http.ResponseWriter
	zw *gzip.Writer
}

func newCompressWriter(w http.ResponseWriter) *compressWriter {
	return &compressWriter{
		w:  w,
		zw: gzip.NewWriter(w),
	}
}

func (c *compressWriter) Header() http.Header {
	return c.w.Header()
}

func (c *compressWriter) Write(p []byte) (int, error) {
	return c.zw.Write(p)
}

func (c *compressWriter) WriteHeader(statusCode int) {
	if statusCode < 300 {
		c.w.Header().Set("Content-Encoding", "gzip")
	}
	c.w.WriteHeader(statusCode)
}

func (c *compressWriter) Close() error {
	return c.zw.Close()
}

type compressReader struct {
	r  io.ReadCloser
	zr *gzip.Reader
}

func newCompressReader(r io.ReadCloser) (*compressReader, error) {
	zr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}

	return &compressReader{
		r:  r,
		zr: zr,
	}, nil
}

func (c compressReader) Read(p []byte) (n int, err error) {
	return c.zr.Read(p)
}

func (c *compressReader) Close() error {
	if err := c.r.Close(); err != nil {
		return err
	}
	return c.zr.Close()
}

func checkSupportsGzip(r *http.Request) bool {
	acceptEncoding := r.Header.Get("Accept-Encoding")
	return strings.Contains(acceptEncoding, "gzip")
}

func checkSendsGzip(r *http.Request) bool {
	contentEncoding := r.Header.Get("Content-Encoding")
	return strings.Contains(contentEncoding, "gzip")
}

func checkSupportsContentType(r *http.Request) bool {
	contentType := r.Header.Get("Content-Type")

	for _, c := range supportedContentTypes {
		if strings.Contains(contentType, c) {
			return true
		}
	}

	return false
}

func WithGzip(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ow := w

		if checkSupportsGzip(r) && checkSupportsContentType(r) {
			cw := newCompressWriter(w)
			ow = cw
			defer cw.Close()
		}

		if checkSendsGzip(r) {
			cr, err := newCompressReader(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			r.Body = cr
			defer cr.Close()
		}

		h.ServeHTTP(ow, r)
	}
}
