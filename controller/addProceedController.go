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
	"model"
	"net/http"
)

func addProceed(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	var body = make([]byte, r.ContentLength)
	r.Body.Read(body)

	c.Debugf("%s", string(body))

	var proceed model.Proceed
	if err := json.Unmarshal(body, &proceed); err != nil {
		c.Errorf("JSON Unmarshal Error %s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if _, _, err := proceed.PutNew(c); err != nil {
		c.Errorf("Proceed Put Error %s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return		
	}

}