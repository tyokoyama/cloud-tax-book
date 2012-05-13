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

type Book struct {
	Type int64 									// 費用種別（通信費、事業主貸、事業主借…）（勘定科目データのキーのID）
	TypeName string `datastore:",noindex"`		// 費用種別名
	Date time.Time      						// 登録日付
	IsCopied bool 								// 転記済みかどうか。
	ExpenseKeyId int64 `datastore:",noindex"`	// 経費帳へのキー
	NonExpenseKeyId int64 `datastore:",noindex"` // 債権債務帳へのキー
	IsExpense bool `datastore:",noindex"`		// 経費かどうか（true：経費／false：経費ではない）
	Detail string `datastore:",noindex"`		// 明細
	MoneySalesIn int64 `datastore:",noindex"`	// 預金売上入金
	MoneyIn int64 `datastore:",noindex"`		// その他入金
	MoneySalesOut int64 `datastore:",noindex"`	// 預金仕入
	MoneyOut int64 `datastore:",noindex"`		// その他出金
	Balance int64 `datastore:",noindex"`		// 残高
}

const kindName = "Book"

func QueryBook(c appengine.Context) ([]*datastore.Key, []Book, error) {
	if c == nil {
		return nil, nil, nil
	}

	q := datastore.NewQuery(kindName)
	if count, err := q.Count(c); err != nil {
		return nil, nil, err
	} else {
		books := make([]Book, 0, count)
		keys, getErr := q.GetAll(c, &books)
		if getErr != nil {
			return nil, nil, getErr
		}

		return keys, books, nil
	}

	return nil, nil, nil
}

// 転記していないエンティティを取得する。
func QueryBookCopyYet(c appengine.Context) ([]*datastore.Key, []Book, error) {
	if c == nil {
		return nil, nil, nil
	}

	q := datastore.NewQuery(kindName).Filter("IsCopied = ", false)
	if count, err := q.Count(c); err != nil {
		return nil, nil, err
	} else {
		books := make([]Book, 0, count)
		keys, getErr := q.GetAll(c, &books)
		if getErr != nil {
			return nil, nil, getErr
		}

		return keys, books, nil
	}

	return nil, nil, nil
}

// データストアからのGET
func GetBook(c appengine.Context, key *datastore.Key) (*Book, error) {
	if c == nil || key == nil {
		return nil, nil
	}

	var book Book

	if err := datastore.Get(c, key, &book); err != nil {
		return nil, err
	}

	return &book, nil
}

// データストアへのPUT（新規登録）
func (book *Book)PutNew(c appengine.Context) (*Book, *datastore.Key, error) {
	if c == nil {
		return nil, nil, nil
	}

	key := datastore.NewIncompleteKey(c, kindName, nil)

	return book.Put(c, key)
}

// データストアへのPUT
func (book *Book)Put(c appengine.Context, key *datastore.Key) (*Book, *datastore.Key, error) {
	if c == nil || key == nil {
		return nil, nil, nil
	}

	if putKey, err := datastore.Put(c, key, book); err != nil {
		return nil, nil, err
	} else {
		return book, putKey, nil
	}

	return nil, nil, nil
}