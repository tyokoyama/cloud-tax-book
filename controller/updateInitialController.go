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
	"fmt"
	"model"
	"net/http"
	"strconv"
)

func updateinitial(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	startcash := r.FormValue("startcash")
	startbook := r.FormValue("startbook")

	c.Infof("startcash = %s", startcash)
	c.Infof("startbook = %s", startbook)

	var initial model.Initial
	initial.StartCash, _ = strconv.ParseInt(startcash, 10, 64)
	initial.StartBook, _ = strconv.ParseInt(startbook, 10, 64)

	if err := initial.Put(c); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("{error: " + err.Error() + "}"))
	} else {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("{\"status\": \"ok\"}")))
	}
}
