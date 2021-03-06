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
	"encoding/json"
	"net/http"
	"model"
)

// var proceedTemplate = template.Must(template.ParseFiles("view/proceeds.html"))

func proceedController(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	c.Debugf("proceedController")
	_, proceeds, err := model.QueryProceed(c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return		
	}

	if response, jsonerr := json.Marshal(&proceeds); jsonerr != nil {
		c.Errorf("JSON Error %s", jsonerr.Error())
		http.Error(w, "Error", http.StatusInternalServerError)
	} else {
		w.Header().Add("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		w.Write(response)
	}

}