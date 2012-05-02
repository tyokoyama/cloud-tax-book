package controller

import (
	"appengine"
	"html/template"
	"net/http"
	"model"
)

var htmltemplate = template.Must(template.ParseFiles("view/index.html"))
func indexHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	cashes, err := model.QueryCash(c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var viewcashes = make([]model.ViewCash, len(cashes))
	for pos, cash := range cashes {
		viewcashes[pos].Create(cash)
	}

	if err := htmltemplate.Execute(w, viewcashes); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
