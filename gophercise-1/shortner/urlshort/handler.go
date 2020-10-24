package urlshort

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/boltdb/bolt"
	"gopkg.in/yaml.v2"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r != nil && r.URL != nil {
			url, ok := pathsToUrls[r.URL.Path]
			if ok {
				http.Redirect(w, r, url, 302)
			} else {
				fallback.ServeHTTP(w, r)
			}
		} else {
			fallback.ServeHTTP(w, r)
		}

	}
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//     - path: /some-path
//       url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
type shortnerPair struct {
	Path string
	Url  string
}

func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	var mapSlice []shortnerPair
	err := yaml.Unmarshal(yml, &mapSlice)
	if err != nil {
		return nil, err
	}
	pathsToUrls := make(map[string]string)
	fmt.Println(mapSlice)
	for _, line := range mapSlice {
		pathsToUrls[string(line.Path)] = line.Url
	}
	return MapHandler(pathsToUrls, fallback), nil
}

func JSONHandler(jsn []byte, fallback http.Handler) (http.HandlerFunc, error) {
	var mapSlice []shortnerPair
	err := json.Unmarshal(jsn, &mapSlice)
	if err != nil {
		return nil, err
	}
	pathsToUrls := make(map[string]string)
	fmt.Println(mapSlice)
	for _, line := range mapSlice {
		pathsToUrls[string(line.Path)] = line.Url
	}
	return MapHandler(pathsToUrls, fallback), nil
}
func DBHandler(db *bolt.DB, fallback http.Handler) http.HandlerFunc {
	// Access data from within a read-only transactional block.

	return func(w http.ResponseWriter, r *http.Request) {
		if r != nil && r.URL != nil {
			if err := db.View(func(tx *bolt.Tx) error {
				url := tx.Bucket([]byte("shorturl")).Get([]byte(r.URL.Path))

				if url != nil {
					http.Redirect(w, r, string(url), 302)
				} else {
					fallback.ServeHTTP(w, r)
				}
				return nil
			}); err != nil {
				log.Fatal(err)
			}
		} else {
			fallback.ServeHTTP(w, r)
		}

	}
}
