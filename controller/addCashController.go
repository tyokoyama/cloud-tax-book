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
	"strconv"
	"time"
)

func addcash(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	date := r.FormValue("date")
	detail := r.FormValue("detail")
	moneysalesin := r.FormValue("moneysalesin")
	moneyin := r.FormValue("moneyin")
	moneysalesout := r.FormValue("moneysalesout")
	moneyout := r.FormValue("moneyout")
	balance := r.FormValue("balance")

	var cash model.Cash
	cash.Date, _ = time.Parse("2006-01-02", date)
	cash.Detail = detail
	cash.MoneySalesIn, _ = strconv.ParseInt(moneysalesin, 10, 64)
	cash.MoneyIn, _ = strconv.ParseInt(moneyin, 10, 64)
	cash.MoneySalesOut, _ = strconv.ParseInt(moneysalesout, 10, 64)
	cash.MoneyOut, _ = strconv.ParseInt(moneyout, 10, 64)
	cash.Balance, _ = strconv.ParseInt(balance, 10, 64)

	_, err := cash.PutNew(c)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("{error: " + err.Error() + "}"))
	} else {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "ok"}`))
	}
}