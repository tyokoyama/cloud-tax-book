/*
Copyright 2012 Takashi Yokoyama

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package controller

import (
	"appengine"
	"fmt"
	"html/template"
	"net/http"
	"model"
)

var htmltemplate = template.Must(template.ParseFiles("view/index.html"))
func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	c := appengine.NewContext(r)
	keys, cashes, err := model.QueryCash(c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	keybooks, books, err := model.QueryBook(c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	keytypes, types, err := model.QueryType(c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var initial model.Initial
	if err := initial.ReadDatastore(c); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var viewcashes = make([]model.ViewCash, len(cashes))
	for pos, cash := range cashes {
		viewcashes[pos].Create(cash, keys[pos].IntID())
	}

	var viewbooks = make([]model.ViewBook, len(books))
	for pos, book := range books {
		viewbooks[pos].Create(book, keybooks[pos].IntID())
	}

	var viewtypes = make([]model.ViewType, len(types))
	for pos, t := range types {
		viewtypes[pos].Create(t, keytypes[pos].IntID())
	}

	var views struct {
		Cashes []model.ViewCash
		Books []model.ViewBook
		Category []model.ViewType
		StartCash int64
		StartBook int64
		CurrentCashBalance string
		CurrentBookBalance string
	}

	views.Cashes = viewcashes
	views.Books = viewbooks
	views.Category = viewtypes
	views.StartCash = initial.StartCash
	views.StartBook = initial.StartBook
	if len(viewcashes) <= 0 {
		views.CurrentCashBalance = fmt.Sprintf("%d", initial.StartCash)
	} else {
		views.CurrentCashBalance = viewcashes[len(viewcashes) - 1].Balance
	}

	if len(viewbooks) <= 0 {
		views.CurrentBookBalance = fmt.Sprintf("%d", initial.StartBook)
	} else {
		views.CurrentBookBalance = viewbooks[len(viewbooks) - 1].Balance
	}

	if err := htmltemplate.Execute(w, views); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
