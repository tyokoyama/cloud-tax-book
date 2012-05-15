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
	"appengine/datastore"
	"model"
	"net/http"
)

func init() {
	http.HandleFunc("/cron/classification", classification)
}

// 経費帳／債権債務帳への書き込み
func classification(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	cronHeader := r.Header.Get("X-AppEngine-Cron")
	if cronHeader == "" {
		// cron jobによる起動ではない。
		c.Infof("Manual Execute")
	}

	// 現金出納帳を読み取る。
	if cashKeys, cashes, err := model.QueryCashCopyYet(c); err == nil {
		c.Infof("Cash Copy Count %d", len(cashes))

		if len(cashes) > 0 {
			var expenses []model.Expense
			var nonexpenses []model.NonExpense

			var cashToExpenseKey []*datastore.Key
			var cashToNonExpenseKey []*datastore.Key
			var cashToExpense []model.Cash
			var cashToNonExpense []model.Cash
			for pos, cash := range cashes {
				// それぞれを経費帳／債権債務帳に書き込み
				if cash.IsExpense == true {
					// 経費
					var ex model.Expense
					ex.CreateExpenseFromCash(cash, cashKeys[pos].IntID())

					expenses = append(expenses, ex)
					cashToExpenseKey = append(cashToExpenseKey, cashKeys[pos])
					cashToExpense = append(cashToExpense, cash)
				} else {
					// 経費ではない
					var ex model.NonExpense
					ex.CreateNonExpenseFromCash(cash, cashKeys[pos].IntID())

					nonexpenses = append(nonexpenses, ex)
					cashToNonExpenseKey = append(cashToNonExpenseKey, cashKeys[pos])
					cashToNonExpense = append(cashToNonExpense, cash)
				}
			}

			// 経費帳へのバッチPUT
			c.Infof("Cash Expense Batch Put Count %d", len(expenses))
			if keys, puterr := model.BatchExpensePut(c, expenses); puterr != nil {
				c.Errorf("Cash Expense Batch Put Error [%s]", puterr.Error())
			} else {
				// 経費分の現金データのバッチPUT（転記済みフラグを立てるのと、キーの退避）
				for pos, key := range keys {
					cashToExpense[pos].ExpenseKeyId = key.IntID()
					cashToExpense[pos].IsCopied = true
				}

				// PUT
				if _, cashbatchputErr := model.BatchCashPut(c, cashToExpenseKey, cashToExpense); cashbatchputErr != nil {
					c.Errorf("Cash (Expense) Batch Put Error [%s]", cashbatchputErr.Error())
				}
			}

			// 債権債務帳へのバッチPUT
			c.Infof("Cash NonExpense Batch Put Count %d", len(nonexpenses))
			if keys, puterr := model.BatchNonExpensePut(c, nonexpenses); puterr != nil {
				c.Errorf("Cash NonExpense Batch Put Error [%s]", puterr.Error())
			} else {
				// 経費にならない分の現金データのバッチPUT（転記済みフラグを立てるのと、キーの退避）
				for pos, key := range keys {
					cashToNonExpense[pos].NonExpenseKeyId = key.IntID()
					cashToNonExpense[pos].IsCopied = true
				}

				// PUT
				if _, cashbatchputErr := model.BatchCashPut(c, cashToNonExpenseKey, cashToNonExpense); cashbatchputErr != nil {
					c.Errorf("Cash (NonExpense) Batch Put Error [%s]", cashbatchputErr.Error())
				}				
			}
		}

	} else {
		c.Errorf("Cash Query Error [%s]", err.Error())
	}

	// 預金出納帳を読み取る。
	if bookKeys, books, err := model.QueryBookCopyYet(c); err == nil {
		c.Infof("Book Copy Count %d", len(books))

		if len(books) > 0 {
			var expenses []model.Expense
			var nonexpenses []model.NonExpense
			var bookToExpenseKey []*datastore.Key
			var bookToNonExpenseKey []*datastore.Key
			var bookToExpense []model.Book
			var bookToNonExpense []model.Book
			for pos, book := range books {
				// それぞれを経費帳／債権債務帳に書き込み
				if book.IsExpense == true {
					// 経費
					var ex model.Expense
					ex.CreateExpenseFromBook(book, bookKeys[pos].IntID())

					expenses = append(expenses, ex)
					bookToExpenseKey = append(bookToExpenseKey, bookKeys[pos])
					bookToExpense = append(bookToExpense, book)
				} else {
					// 経費ではない
					var ex model.NonExpense
					ex.CreateNonExpenseFromBook(book, bookKeys[pos].IntID())

					nonexpenses = append(nonexpenses, ex)

					bookToNonExpenseKey = append(bookToNonExpenseKey, bookKeys[pos])
					bookToNonExpense = append(bookToNonExpense, book)
				}
			}

			// バッチPUT
			c.Infof("Book Expense Batch Put Count %d", len(expenses))
			if keys, puterr := model.BatchExpensePut(c, expenses); puterr != nil {
				c.Errorf("Book Batch Put Error [%s]", puterr.Error())
			} else {
				// 経費分の預金データのバッチPUT（転記済みフラグを立てるのと、キーの退避）
				for pos, key := range keys {
					bookToExpense[pos].ExpenseKeyId = key.IntID()
					bookToExpense[pos].IsCopied = true
				}

				// PUT
				if _, bookbatchputErr := model.BatchBookPut(c, bookToExpenseKey, bookToExpense); bookbatchputErr != nil {
					c.Errorf("Book (Expense) Batch Put Error [%s]", bookbatchputErr.Error())
				}
			}

			// 債権債務帳へのバッチPUT
			c.Infof("Book NonExpense Batch Put Count %d", len(nonexpenses))
			if keys, puterr := model.BatchNonExpensePut(c, nonexpenses); puterr != nil {
				c.Errorf("Book NonExpense Batch Put Error [%s]", puterr.Error())
			} else {
				// 経費にならない分の現金データのバッチPUT（転記済みフラグを立てるのと、キーの退避）
				for pos, key := range keys {
					bookToNonExpense[pos].NonExpenseKeyId = key.IntID()
					bookToNonExpense[pos].IsCopied = true
				}

				// PUT
				if _, bookbatchputErr := model.BatchBookPut(c, bookToNonExpenseKey, bookToNonExpense); bookbatchputErr != nil {
					c.Errorf("Book (NonExpense) Batch Put Error [%s]", bookbatchputErr.Error())
				}				
			}
		}
	} else {
		c.Errorf("Book Query Error [%s]", err.Error())
	}

}
