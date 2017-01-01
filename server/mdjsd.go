package main

import (
	"io/ioutil"
	"log"
	"net/http"
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
	written := false

	if b, err := ioutil.ReadFile("default.js"); err == nil {
		w.Write(b)
		written = true
	}

	for host != "" {
		if b, err := ioutil.ReadFile(host + ".js"); err == nil {
			w.Write(b)
			written = true
		}
		if sp := strings.SplitN(host, ".", 2); len(sp) == 2 {
			host = sp[1]
		} else {
			host = ""
		}
	}

	if !written {
		w.WriteHeader(http.StatusNoContent)
	}
}

func detectOrigin(req *http.Request, host string) (origin string, ok bool) {
	origin = req.Header.Get("Origin")
	ok = origin != "" && host != "" && strings.HasSuffix(origin, host)
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
