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
package model

import (
	"appengine"
	"appengine/datastore"
	"time"
)

type Expense struct {
	Type int64 									// 費用種別（通信費、事業主貸、事業主借…）（勘定科目データのキーのID）
	TypeName string `datastore:",noindex"`		// 費用種別名
	Date time.Time      						// 登録日付
	CashId int64 `datastore:",noindex"`			// 現金出納帳	へのキー
	BookId int64 `datastore:",noindex"` 		// 預金出納帳	へのキー
	Detail string `datastore:",noindex"`		// 明細
	MoneySalesIn int64 `datastore:",noindex"`	// 売上入金
	MoneyIn int64 `datastore:",noindex"`		// その他入金
	MoneySalesOut int64 `datastore:",noindex"`	// 仕入
	MoneyOut int64 `datastore:",noindex"`		// その他出金
}

const expenseKindName = "expense"

func (ex *Expense) CreateExpenseFromCash(cash Cash, keyId int64) {
	ex.Type = cash.Type
	ex.TypeName = cash.TypeName
	ex.Date = cash.Date
	ex.CashId = keyId			// 現金出納帳へのキーのID
	ex.BookId = 0				// 預金出納帳へはキーなし
	ex.Detail = cash.Detail
	ex.MoneySalesIn = cash.MoneySalesIn
	ex.MoneyIn = cash.MoneyIn
	ex.MoneySalesOut = cash.MoneySalesOut
	ex.MoneyOut = cash.MoneyOut
}

func (ex *Expense) CreateExpenseFromBook(book Book, keyId int64) {
	ex.Type = book.Type
	ex.TypeName = book.TypeName
	ex.Date = book.Date
	ex.CashId = 0					// 現金出納帳へはキーなし
	ex.BookId = keyId				// 預金出納帳へのキーのID
	ex.Detail = book.Detail
	ex.MoneySalesIn = book.MoneySalesIn
	ex.MoneyIn = book.MoneyIn
	ex.MoneySalesOut = book.MoneySalesOut
	ex.MoneyOut = book.MoneyOut	
}

func GetExpense(c appengine.Context, key *datastore.Key) (*Expense, error) {
	if c == nil || key == nil {
		return nil, nil
	}

	var expense Expense
	if err := datastore.Get(c, key, &expense); err != nil {
		return nil, err
	}

	return &expense, nil
}

func QueryExpense(c appengine.Context) ([]Expense, error) {
	q := datastore.NewQuery(expenseKindName)
	if count, err := q.Count(c); err != nil {
		return nil, err
	} else {
		expenses := make([]Expense, 0, count)
		_, getErr := q.GetAll(c, &expenses)
		if getErr != nil {
			return nil, getErr
		}

		return expenses, nil
	}

	return nil, nil
}

func BatchExpensePut(c appengine.Context, expenses []Expense) ([]*datastore.Key, error) {
	if c == nil {
		return nil, nil
	}

	var keys []*datastore.Key
	for _, expense := range expenses {
		// ここでは、CashIdかBookIdのどちらかにしかIDはセットされていない前提とする。
		// ID=0の時はIncomplete()がtrueになる。
		key := datastore.NewKey(c, expenseKindName, "", expense.CashId, nil)
		if key.Incomplete() == true {
			key = datastore.NewKey(c, expenseKindName, "", expense.BookId, nil)
		}
		keys = append(keys, key)
	}

	return datastore.PutMulti(c, keys, expenses);
}

func DeleteExpense(c appengine.Context, id int64) error {
	if c == nil {
		return nil
	}

	key := datastore.NewKey(c, expenseKindName, "", id, nil)
	if key.Incomplete() {
		c.Infof("DeleteExpense: key id incomplete")
		return nil
	}
	return datastore.Delete(c, key)
}

func (ex *Expense) PutNew(c appengine.Context) (*datastore.Key, *Expense, error) {
	if c == nil {
		return nil, nil, nil
	}

	key := datastore.NewIncompleteKey(c, expenseKindName, nil)
	return ex.Put(c, key)
}

func (ex *Expense) Put(c appengine.Context, key *datastore.Key) (*datastore.Key, *Expense, error) {
	if c == nil || key == nil {
		return nil, nil, nil
	}

	if key, err := datastore.Put(c, key, ex); err != nil {
		return nil, nil, err
	} else {
		return key, ex, nil
	}

	return nil, nil, nil
}
