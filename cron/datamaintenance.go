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
package cron

import (
	"appengine"
	"model"
	"net/http"
)

// データメンテナンス用バッチ処理
func datamaintenance(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	cronHeader := r.Header.Get("X-AppEngine-Cron")
	if cronHeader == "" {
		// cron jobによる起動ではない。
		c.Infof("Manual Execute")
	}

	// 残高調整（入力順が逆転した時など、狂ってしまったものを再計算

	// 元入金を取得
	var init model.Initial

	if err := init.ReadDatastore(c); err != nil {
		c.Errorf("model.Initialize.ReadDatastore Error %s", err.Error())
		return
	}

	// Read SettlementSummary
	var summary model.SettlementSummary
	if err := summary.Get(c); err != nil {
		c.Errorf("model.SettlementSummary.Get Error %s", err.Error())
		return
	}

	// 現金データ
	if cashKeys, cashes, err := model.QueryCash(c); err != nil {
		c.Errorf("model.QueryCash Error %s", err.Error())
	} else {
		var balance int64 = init.StartCash
		for pos, cash := range cashes {
			cashIO := (cash.MoneySalesIn + cash.MoneyIn) - (cash.MoneySalesOut + cash.MoneyOut)
			balance = balance + cashIO			// 残高に収支を加算 = 残高

			if balance != cash.Balance {
				// 残高が一致しない場合、上書き
				c.Infof("Balance Not Equal Put %s %s", cash.Date.Format("2006-01-02"), cash.Detail)
				cash.Balance = balance
				if _, _, putErr := cash.Put(c, cashKeys[pos]); putErr != nil {
					c.Errorf("Cash Put Error %s", putErr.Error())
				}
			}
		}

		summary.EndCash = balance
	}

	// 預金データ
	if bookKeys, books, err := model.QueryBook(c); err != nil {
		c.Errorf("model.QueryBook Error %s", err.Error())
	} else {
		var balance int64 = init.StartBook
		for pos, book := range books {
			bookIO := (book.MoneySalesIn + book.MoneyIn) - (book.MoneySalesOut + book.MoneyOut)
			balance = balance + bookIO			// 残高に収支を加算 = 残高

			if balance != book.Balance {
				// 残高が一致しない場合、上書き
				c.Infof("Balance Not Equal Put %s %s", book.Date.Format("2006-01-02"), book.Detail)
				book.Balance = balance
				if _, _, putErr := book.Put(c, bookKeys[pos]); putErr != nil {
					c.Errorf("Book Put Error %s", putErr.Error())
				}
			}
		}

		summary.EndBook = balance
	}

	// Put SettlementSummary
	if err := summary.Put(c); err != nil {
		c.Errorf("model.SettlementSummary.Put Error %s", err.Error())
		return		
	}

	// Summary
	if err := model.CalcSummary(c); err != nil {
		c.Errorf("CalcSummary Error %s", err.Error())
	}

}
