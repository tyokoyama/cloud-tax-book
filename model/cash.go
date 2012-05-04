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

type Cash struct {
	Type int 									// 費用種別（通信費、事業主貸、事業主借…）
	Date time.Time      						// 登録日付
	IsExpense bool `datastore:",noindex"`		// 経費かどうか（true：経費／false：経費ではない）
	Detail string `datastore:",noindex"`		// 明細
	MoneySalesIn int64 `datastore:",noindex"`	// 現金売上入金
	MoneyIn int64 `datastore:",noindex"`		// その他入金
	MoneySalesOut int64 `datastore:",noindex"`	// 現金仕入
	MoneyOut int64 `datastore:",noindex"`		// その他出金
	Balance int64 `datastore:",noindex"`		// 残高
}

func QueryCash(c appengine.Context) ([]*datastore.Key, []Cash, error) {
	if c == nil {
		return nil, nil, nil
	}

	q := datastore.NewQuery("Cash")
	if count, err := q.Count(c); err != nil {
		return nil, nil, err
	} else {
		cashes := make([]Cash, 0, count)
		keys, getErr := q.GetAll(c, &cashes)
		if getErr != nil {
			return nil, nil, getErr
		}

		return keys, cashes, nil
	}

	return nil, nil, nil
}

// データストアからのGET
func GetCash(c appengine.Context, key *datastore.Key) (*Cash, error) {
	if c == nil || key == nil {
		return nil, nil
	}

	var cash Cash

	if err := datastore.Get(c, key, &cash); err != nil {
		return nil, err
	}

	return &cash, nil
}

// データストアへのPUT（新規登録）
func (cash *Cash)PutNew(c appengine.Context) (*Cash, *datastore.Key, error) {
	if c == nil {
		return nil, nil, nil
	}

	key := datastore.NewIncompleteKey(c, "Cash", nil)

	return cash.Put(c, key)
}

// データストアへのPUT
func (cash *Cash)Put(c appengine.Context, key *datastore.Key) (*Cash, *datastore.Key, error) {
	if c == nil || key == nil {
		return nil, nil, nil
	}

	if putKey, err := datastore.Put(c, key, cash); err != nil {
		return nil, nil, err
	} else {
		return cash, putKey, nil
	}

	return nil, nil, nil
}
