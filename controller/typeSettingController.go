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

import(
	"appengine"
	"model"
	"net/http"
	"html/template"
)

var htmltypesettingtemplate = template.Must(template.ParseFiles("view/typesetting.html"))
func typesetting(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	c.Infof("typesetting")

	keys, types, _ := model.QueryType(c)

	var views = make(
		[]struct {
			Id int64
			Type model.Type
		}, len(keys))

	for i := 0; i < len(views); i++ {
		views[i].Id = keys[i].IntID()
		views[i].Type = types[i]
	}

	if err := htmltypesettingtemplate.Execute(w, views); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}