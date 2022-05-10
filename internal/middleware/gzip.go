package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

const compression = "gzip"

type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w gzipWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func Gzip(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if strings.Contains(request.Header.Get("Content-Encoding"), compression) {
			reader, err := gzip.NewReader(request.Body)

			if err != nil {
				http.Error(writer, err.Error(), http.StatusInternalServerError)
				return
			}

			request.Body = reader
			defer request.Body.Close()
		}

		if strings.Contains(request.Header.Get("Accept-Encoding"), compression) {
			gz, err := gzip.NewWriterLevel(writer, gzip.BestSpeed)

			if err != nil {
				http.Error(writer, err.Error(), http.StatusInternalServerError)
				return
			}

			defer gz.Close()
			writer.Header().Set("Content-Encoding", compression)
			writer = gzipWriter{ResponseWriter: writer, Writer: gz}
		}

		next.ServeHTTP(writer, request)
	})
}
