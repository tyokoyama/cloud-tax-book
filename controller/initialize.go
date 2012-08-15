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
	"net/http"
)

func init() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/addcash", addcash)
	http.HandleFunc("/updatecash", updatecash)
	http.HandleFunc("/addbook", addbook)
	http.HandleFunc("/updatebook", updatebook)
	http.HandleFunc("/updateinitial", updateinitial)
	http.HandleFunc("/typesetting", typesetting)
	http.HandleFunc("/addtype", addtype)
	http.HandleFunc("/updatetype", updatetype)
	http.HandleFunc("/deletetype", deletetype)
	http.HandleFunc("/proceeddata", proceedController)
	http.HandleFunc("/addproceed", addProceed)
}
