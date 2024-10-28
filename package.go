package main

import (
	"embed"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"text/template"
	"time"
)

func logware(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		// Initialize the status to 200 in case WriteHeader is not called
		rec := statusRecorder{w, 200}

		// do
		next.ServeHTTP(&rec, r)

		// clean up sensitive values in query
		query := r.URL.Query()
		if k := "access_token"; query.Has(k) {
			query.Set("access_token", "...")
		}
		// log request/response and
		path := r.URL.Path
		if v := query.Encode(); v != "" {
			path += "?" + v
		}
		debug.Println(r.Method, path, rec.status, time.Since(start))
	}
}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (rec *statusRecorder) WriteHeader(code int) {
	rec.status = code
	rec.ResponseWriter.WriteHeader(code)
}

func init() {
	page = template.Must(
		template.New("").Funcs(funcMap).ParseFS(asset, "htdocs/*"),
	)
}

var page *template.Template
var funcMap = template.FuncMap{
	"doX": func() string { return "x" },
}

//go:embed htdocs
var asset embed.FS

var debug = log.New(ioutil.Discard, "D ", log.LstdFlags|log.Lshortfile)

func init() {
	if yes, _ := strconv.ParseBool(os.Getenv("D")); yes {
		debug.SetOutput(os.Stderr)
	}
}
