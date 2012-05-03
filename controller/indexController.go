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
