package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/hectorgabucio/genity-hector/internal/data"
)

const POST_DATA_PATH = "/post-data/"
const GET_DATA_PATH = "/get-data/"

type Title struct {
	Title string
}

type app struct {
	dataRepository data.DataRepository
}

func (app *app) postData(w http.ResponseWriter, req *http.Request) {

	if req.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	title, err := processTitleParam(req.URL.Path, POST_DATA_PATH)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	data := &data.Data{Title: title}

	if err := app.dataRepository.Add(data); err != nil {
		log.Println(err)
		http.Error(w, "Internal error adding new data", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	respondJson(Title{Title: data.Title}, w)
}

func (app *app) getData(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	title, err := processTitleParam(req.URL.Path, GET_DATA_PATH)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	data := &data.Data{Title: title}

	retrieved, err := app.dataRepository.Get(data)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal error getting data", http.StatusInternalServerError)
		return
	}
	if retrieved == nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	respondJson(retrieved, w)
}

func respondJson(body interface{}, w http.ResponseWriter) {
	response, err := json.Marshal(body)
	if err != nil {
		http.Error(w, "Error parsing to JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, string(response))
}

func processTitleParam(path string, prefix string) (string, error) {
	title := strings.TrimPrefix(path, prefix)
	if title == "" {
		return "", errors.New("bad title in request")
	}
	return title, nil
}

func main() {

	app := app{data.NewDataRepository()}
	defer app.dataRepository.CloseConn()

	http.HandleFunc(POST_DATA_PATH, app.postData)
	http.HandleFunc(GET_DATA_PATH, app.getData)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalln(err)
	}
}
