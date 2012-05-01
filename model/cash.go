package model

import (
	"appengine"
	"appengine/datastore"
	"time"
)

type Cash struct {
	Key *datastore.Key // 現金データへのキー
	Date time.Time       // 登録日付
	Detail string	// 明細
	MoneySalesIn int64	// 現金売上入金
	MoneyIn int64		// その他入金
	MoneySalesOut int64	// 現金仕入
	MoneyOut int64		// その他出金
	Balance int64		// 残高
}

func QueryCash(c appengine.Context) ([]Cash, error) {
	if c == nil {
		return nil, nil
	}

	q := datastore.NewQuery("Cash")
	if count, err := q.Count(c); err != nil {
		return nil, err
	} else {
		cashes := make([]Cash, count)
		_, getErr := q.GetAll(c, &cashes)
		if getErr != nil {
			return nil, getErr
		}

		return cashes, nil
	}

	return nil, nil
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
func (cash *Cash)PutNew(c appengine.Context) (*Cash, error) {
	if c == nil {
		return nil, nil
	}

	key := datastore.NewIncompleteKey(c, "cash", nil)

	cash.Key = key
	cash.Date = time.Now()

	return cash.Put(c)
}

// データストアへのPUT
func (cash *Cash)Put(c appengine.Context) (*Cash, error) {
	if c == nil || cash.Key == nil {
		return nil, nil
	}

	if _, err := datastore.Put(c, cash.Key, cash); err != nil {
		return nil, err
	}

	return cash, nil
}
