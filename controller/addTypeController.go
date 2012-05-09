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
	"encoding/json"
	"fmt"
	"model"
	"net/http"
)

func addtype(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	name := r.FormValue("name")
	expense := r.FormValue("expense")
	detail := r.FormValue("detail")

	c.Infof("name = %s", name)
	c.Infof("expense = %s", expense)
	c.Infof("detail = %s", detail)

	var newType model.Type
	newType.Name = name
	if expense == "" {
		newType.IsExpense = false
	} else {
		newType.IsExpense = true
	}
	newType.Detail = detail

	if _, err := newType.PutNew(c); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("{\"error\": \"%s\"}", err.Error())))
	} else {
		if response, jsonerr := json.Marshal(&newType); jsonerr != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("{\"error\": \"%s\"}", jsonerr.Error())))
		} else {
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(response)
		}
	}

}