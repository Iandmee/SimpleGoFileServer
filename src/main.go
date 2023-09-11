package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func uploadHandler(w http.ResponseWriter, r *http.Request, path *string, filename *string) {
	if r.ContentLength > MAX_UPLOAD_SIZE {
		http.Error(w, "The uploaded file is too big. Please use an image less than 1MB in size", http.StatusBadRequest)
		return
	}
	body, err := io.ReadAll(io.LimitReader(r.Body, MAX_UPLOAD_SIZE))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := os.MkdirAll(filepath.Dir(*path), DEFAULT_FOLDER_PERMISSIONS); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := os.WriteFile(*path, body, DEFAULT_FILE_PERMISSIONS); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = w.Write([]byte(fmt.Sprintf("File \"%s\" successfully uploaded.", *filename)))
	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func downloadHandler(w http.ResponseWriter, r *http.Request, path *string) {
	f, err := os.Open(*path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Print(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}(f)
	fi, err := f.Stat()
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	modTime := fi.ModTime()
	http.ServeContent(w, r, *path, modTime, f)
}

func deleteHandler(w http.ResponseWriter, path *string, filename *string) {
	if err := os.Remove(*path); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err := w.Write([]byte(fmt.Sprintf("File \"%s\" successfully deleted.", *filename)))
	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func basicHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		filename := filepath.Base(path)
		switch r.Method {
		case http.MethodGet:
			downloadHandler(w, r, &path)
		case http.MethodPost:
			uploadHandler(w, r, &path, &filename)
		case http.MethodDelete:
			deleteHandler(w, &path, &filename)
		}
	}
}

var MAX_UPLOAD_SIZE int64 = 1024 * 1024
var SERVER_PORT = 8080
var DEFAULT_FILE_PERMISSIONS = os.FileMode(0644)
var DEFAULT_FOLDER_PERMISSIONS = os.FileMode(0644)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", basicHandler())
	log.Printf("Server started on *:%v", SERVER_PORT)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", SERVER_PORT), mux))
}
