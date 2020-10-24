package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"example.com/morpheus020/gophercise/shortner/urlshort"
	"github.com/boltdb/bolt"
)

func main() {
	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := urlshort.MapHandler(pathsToUrls, mux)

	json := `
[
	{
		"path":"/g",
		"url":"https://google.com"
	},
	{
		"path":"/f",
		"url":"http://facebook.com"
	}
]
`
	jsonHandler, err1 := urlshort.JSONHandler([]byte(json), mapHandler)
	if err1 != nil {
		panic(err1)
	}
	// Build the YAMLHandler using the mapHandler as the
	// fallback
	yaml := `
- path: /urlshort
  url: https://github.com/gophercises/urlshort
- path: /urlshort-final
  url: https://github.com/gophercises/urlshort/tree/solution	
`
	yamlHandler, err := urlshort.YAMLHandler([]byte(yaml), jsonHandler)
	if err != nil {
		panic(err)
	}

	// Open the database.
	db, err := bolt.Open("temp.db", 0666, nil)
	if err != nil {
		panic(err)
	}
	defer os.Remove(db.Path())

	// Insert data into a bucket.
	if err := db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucket([]byte("shorturl"))
		if err != nil {
			return err
		}
		if err := b.Put([]byte("/t"), []byte("https://twitter.com")); err != nil {
			return err
		}
		if err := b.Put([]byte("/l"), []byte("https://linkedin.com")); err != nil {
			return err
		}
		return nil
	}); err != nil {
		log.Fatal(err)
	}
	dbHandler := urlshort.DBHandler(db, yamlHandler)
	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", dbHandler)

	// Close database to release the file lock.
	if err := db.Close(); err != nil {
		log.Fatal(err)
	}
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
