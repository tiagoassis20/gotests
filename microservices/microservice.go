package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

type Metrics struct {
	Requests     int
	BadRequest   int
	ServerErrors int
	Redirect     int
	NoAnswer     int
}

func (m Metrics) String() string {
	return fmt.Sprintf("requests:%d, BadRequest:%d, ServerErrors:%d, Redirect:%d, NoAnswer:%d", m.Requests, m.BadRequest, m.ServerErrors, m.Redirect, m.NoAnswer)

}
func (m *Metrics) MetricsMidleware(logger *log.Logger) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			m.Requests++
			next.ServeHTTP(w, r)
			status := w.Header().Get("Status")
			if status == "" {
				m.NoAnswer++
			} else if status >= "500" {
				m.ServerErrors++
			} else if status >= "400" {
				m.BadRequest++
			} else if status >= "300" {
				m.Redirect++
			}

			logger.Println(m)

		})
	}
}

func main() {
	router := mux.NewRouter()
	logger := log.New(os.Stdout, "logger: ", log.LstdFlags)
	metrics := &Metrics{}

	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger.Println(r.Method, r.RequestURI, r.Proto)
			next.ServeHTTP(w, r)
			logger.Println(r.Method, r.RequestURI, r.Proto, w.Header()["Status"])
		})
	}, metrics.MetricsMidleware(logger))

	router.HandleFunc("/puppies", func(w http.ResponseWriter, _ *http.Request) {
		fmt.Fprintln(w, "Yes, puppies!")
		w.Header().Add("status", fmt.Sprintf("%d", http.StatusOK))
	})

	router.HandleFunc("/puppies/{id}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		fmt.Fprintf(w, "This method will act on the puppy with id %s\n", vars["id"])
		w.Header().Add("status", fmt.Sprintf("%d", http.StatusOK))
	})

	http.ListenAndServe(":5001", router)
}
