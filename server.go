package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

func main() {
	fmt.Println("Starting web server")
	f, err := os.OpenFile("Resume.pdf", os.O_RDONLY, os.ModeDevice)
	if err != nil {
		panic(err)
	}
	resume, err := io.ReadAll(f)
	if err != nil {
		panic(err)
	}
	fileHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write(resume)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	})
	cachedHandler := CacheHandler(fileHandler, 3600)
	err = http.ListenAndServe(":5000", cachedHandler)
	if err != nil {
		panic(err)
	}
}

func CacheHandler(next http.Handler, seconds int) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "public, max-age=3600")
		w.Header().Set("Expires", time.Now().Add(time.Duration(seconds)*time.Second).Format(http.TimeFormat))
		w.Header().Set("Last-Modified", time.Now().UTC().Format(http.TimeFormat))

		next.ServeHTTP(w, r)
	})
}
