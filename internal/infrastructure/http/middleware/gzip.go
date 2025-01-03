package middleware

import (
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"slices"
	"strings"
)

const (
	ContentEncodingHeader = "Content-Encoding"
	ContentEncodingGzip   = "gzip"
	AcceptEncodingHeader  = "Accept-Encoding"
)

var (
	supportedGzipContentTypes = []string{
		"application/json",
		"text/html",
	}
)

type compressorWriter struct {
	http.ResponseWriter
	Writer       io.Writer
	wroteHeader  bool
	compressible bool
}

func (cw *compressorWriter) Write(p []byte) (int, error) {
	if !cw.wroteHeader {
		// Ensure headers are written before any body data
		cw.WriteHeader(http.StatusOK)
	}

	if cw.compressible && cw.Writer != nil {
		return cw.Writer.Write(p)
	}

	// Write uncompressed response if not compressible
	return cw.ResponseWriter.Write(p)
}

func (cw *compressorWriter) WriteHeader(code int) {
	if cw.wroteHeader {
		// Avoid writing headers multiple times
		cw.ResponseWriter.WriteHeader(code)
		return
	}

	cw.wroteHeader = true
	defer cw.ResponseWriter.WriteHeader(code)

	if cw.Header().Get("Content-Encoding") != "" {
		return
	}

	if !cw.isCompressible(code) {
		cw.compressible = false
		return
	}

	cw.compressible = true
	cw.Header().Set("Content-Encoding", ContentEncodingGzip)
	cw.Header().Add("Vary", "Accept-Encoding")
	cw.Header().Del("Content-Length")
}

func (cw *compressorWriter) isCompressible(code int) bool {
	if (code >= 100 && code < 200) || code == http.StatusNoContent {
		return false
	}
	
	contentType := cw.Header().Get("Content-Type")
	for _, ct := range supportedGzipContentTypes {
		if strings.HasPrefix(contentType, ct) {
			return true
		}
	}

	return false
}

func WithCompress(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !isGizipSupported(r) {
			fmt.Println("doesn't supports gzip")
			next.ServeHTTP(w, r)
			return
		}

		g := gzip.NewWriter(w)
		defer g.Close()

		gz := &compressorWriter{
			ResponseWriter: w,
			Writer:         g,
		}

		next.ServeHTTP(gz, r)
	})
}

func isGizipSupported(r *http.Request) bool {
	if r.Header.Get(AcceptEncodingHeader) == "" {
		return false
	}

	acceptEncoding := r.Header.Get(AcceptEncodingHeader)
	supportedCompressionFormats := strings.Split(acceptEncoding, ",")

	return slices.Contains(supportedCompressionFormats, ContentEncodingGzip)
}
