package app

import (
	"fmt"
	"github.com/ylapshin/urlshortener/internal/urlstg"
	"io"
	"net/http"
)

const respPath = "/"

type App struct {
	data *urlstg.URLStg
}

func fullURL(req *http.Request, id string) string {
	proto := "https://"
	if req.TLS == nil {
		proto = "http://"
	}

	return proto + req.Host + respPath + id
}

func (app *App) rootMethodGet(res http.ResponseWriter, req *http.Request) {
	id := req.URL.Path[1:]
	if len(id) != urlstg.KeyLen {
		http.Error(res, fmt.Sprintf("%s - wrong URL id", id), http.StatusBadRequest)
		return
	}

	url, err := app.data.Get(id)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	res.Header().Set("Location", url)
	res.WriteHeader(http.StatusTemporaryRedirect)
}

func (app *App) rootMethodPost(res http.ResponseWriter, req *http.Request) {
	contentType := req.Header.Get("Content-Type")
	if contentType != "text/plain" {
		http.Error(res, fmt.Sprintf("%s - invalid content type", contentType), http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	id, _ := app.data.Reg(string(body))

	res.Header().Set("Content-Type", "text/plain")
	res.WriteHeader(http.StatusCreated)
	_, err = res.Write([]byte(fullURL(req, id)))
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
	}
}

func (app *App) RootHandler(res http.ResponseWriter, req *http.Request) {

	if req.Method == http.MethodGet {
		app.rootMethodGet(res, req)
	}

	if req.Method == http.MethodPost {
		app.rootMethodPost(res, req)
	}
}

func New() *App {
	return &App{data: urlstg.New()}
}
