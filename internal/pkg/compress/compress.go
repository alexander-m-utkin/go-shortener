package compress

import (
	"compress/zlib"
	"io"
	"net/http"
	"strings"
)

type zlibWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w zlibWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func ZlibHandle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "deflate") {
			next.ServeHTTP(w, r)
			return
		}

		gz := zlib.NewWriter(w)
		defer gz.Close()

		w.Header().Set("Content-Encoding", "deflate")
		next.ServeHTTP(zlibWriter{ResponseWriter: w, Writer: gz}, r)
	})
}
