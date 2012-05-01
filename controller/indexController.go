package controller

import (
	"appengine"
	"html/template"
	"net/http"
	"model"
)

func init() {
	http.HandleFunc("/", indexHandler)
}

var htmltemplate = template.Must(template.ParseFiles("view/index.html"))
func indexHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	cashes, err := model.QueryCash(c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := htmltemplate.Execute(w, cashes); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
