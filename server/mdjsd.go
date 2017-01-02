package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"strings"
)

func handler(w http.ResponseWriter, req *http.Request) {
	host := getHost(req.URL.Path)

	w.Header().Set("Content-Type", "text/javascript")
	if origin, ok := detectOrigin(req, host); ok {
		w.Header().Set("Access-Control-Allow-Origin", origin)
	}

	writeBody(w, host)
}

func getHost(path string) string {
	path = strings.TrimPrefix(path, "/")
	path = strings.TrimSuffix(path, ".js")
	return path
}

func writeBody(w http.ResponseWriter, host string) {
	l := 0

	l += writeFile(w, "default.js")

	if matches, _ := filepath.Glob("default/*.js"); matches != nil {
		for _, path := range matches {
			l += writeFile(w, path)
		}
	}

	for host != "" {
		l += writeFile(w, host+".js")
		if sp := strings.SplitN(host, ".", 2); len(sp) == 2 {
			host = sp[1]
		} else {
			host = ""
		}
	}

	if l == 0 {
		w.WriteHeader(http.StatusNoContent)
	}
}

func writeFile(w http.ResponseWriter, path string) int {
	if b, err := ioutil.ReadFile(path); err == nil && len(b) > 0 {
		w.Write(b)
		w.Write([]byte{'\n'})
		return len(b) + 1
	}
	return 0
}

func detectOrigin(req *http.Request, host string) (origin string, ok bool) {
	origin = req.Header.Get("Origin")
	ok = origin != "" && host != "" &&
		(strings.HasSuffix(origin, "://"+host) || strings.HasSuffix(origin, "."+host))
	return
}

func main() {
	var f http.HandlerFunc = handler
	srv := &http.Server{
		Addr:    ":3131",
		Handler: f,
	}
	log.Fatal(srv.ListenAndServeTLS("etc/server.crt", "etc/server.key"))
}
