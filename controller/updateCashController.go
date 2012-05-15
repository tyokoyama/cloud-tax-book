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
	"appengine/datastore"
	"appengine/taskqueue"
	"fmt"
	"model"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

func updatecash(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	id := r.FormValue("id")
	date := r.FormValue("date")
	detail := r.FormValue("detail")
	moneyType := r.FormValue("type")
	moneysalesin := r.FormValue("moneysalesin")
	moneyin := r.FormValue("moneyin")
	moneysalesout := r.FormValue("moneysalesout")
	moneyout := r.FormValue("moneyout")
	balance := r.FormValue("balance")

	c.Infof("id = %s", id)
	c.Infof("date = %s", date)
	c.Infof("detail = %s", detail)
	c.Infof("moneyType = %s", moneyType)
	c.Infof("moneysalesin = %s", moneysalesin)
	c.Infof("moneyin = %s", moneyin)
	c.Infof("moneysalesout = %s", moneysalesout)
	c.Infof("moneyout = %s", moneyout)
	c.Infof("balance = %s", balance)

	var keyid int64
	keyid, _ = strconv.ParseInt(id, 10, 64)

	cash, getErr := model.GetCash(c, datastore.NewKey(c, "Cash", "", keyid, nil))
	if getErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("{error: " + getErr.Error() + "}"))
		return 
	}

	cash.Date, _ = time.Parse("2006-01-02", date)
	cash.Type, _ = strconv.ParseInt(moneyType, 10, 64)

	if data, err := model.GetType(c, datastore.NewKey(c, "Type", "", cash.Type, nil)); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("{error: " + err.Error() + "}"))
		return
	} else {
		cash.TypeName = data.Name
		cash.IsExpense = data.IsExpense
	}

	cash.Detail = detail
	cash.MoneySalesIn, _ = strconv.ParseInt(moneysalesin, 10, 64)
	cash.MoneyIn, _ = strconv.ParseInt(moneyin, 10, 64)
	cash.MoneySalesOut, _ = strconv.ParseInt(moneysalesout, 10, 64)
	cash.MoneyOut, _ = strconv.ParseInt(moneyout, 10, 64)
	cash.Balance, _ = strconv.ParseInt(balance, 10, 64)

	// タスクキューに転記先削除のリクエストを追加する。（追加先：Default Queue）
	var tasks []*taskqueue.Task
	paramExpense := url.Values {"id": {fmt.Sprintf("%d", cash.ExpenseKeyId)}}
	tExpense := taskqueue.NewPOSTTask("/task/deleteexpense", paramExpense)
	tasks = append(tasks, tExpense)

	paramNonExpense := url.Values {"id": {fmt.Sprintf("%d", cash.NonExpenseKeyId)}}
	tNonExpense := taskqueue.NewPOSTTask("/task/deletenonexpense", paramNonExpense)
	tasks = append(tasks, tNonExpense)

	if _, err := taskqueue.AddMulti(c, tasks, ""); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("{error: " + err.Error() + "}"))
		return 
	}

	// 転記済フラグを落とす。
	cash.IsCopied = false

	// 転記済みの場合、転記先を削除する。（Keyの参照もなくす）
	cash.ExpenseKeyId = 0
	cash.NonExpenseKeyId = 0

	_, key, err := cash.Put(c, datastore.NewKey(c, "Cash", "", keyid, nil))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("{error: " + err.Error() + "}"))
	} else {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("{\"status\": \"ok\", \"id\": %d}", key.IntID())))
	}
}