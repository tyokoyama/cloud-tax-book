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

type Proceed struct {
	Date time.Time      											// 登録日付
	ViewDate string `datastore:"-"`
	Name string	   `datastore:",noindex"`							//
	Detail string  `datastore:",noindex"`							// 
	Proceed int64  `datastore:",noindex"`							//
	MoneyIn int64  `datastore:",noindex"`							//
	Balance int64  `datastore:",noindex"`							// 
}

const (
	proceedKindName = `proceed`
	TimeFormat = "2006-01-02"
)

func(proceed *Proceed) Load(c <-chan datastore.Property) error {
	if err := datastore.LoadStruct(proceed, c); err != nil {
		return err
	}

	proceed.ViewDate = proceed.Date.Format(TimeFormat)
	if proceed.Date.IsZero() {
		// データがないものは現在時刻にする。
		proceed.ViewDate = time.Now().Format(TimeFormat)
	} else {
		proceed.ViewDate = proceed.Date.Format(TimeFormat)
	}

	return nil
}

func(proceed *Proceed) Save(c chan<- datastore.Property) error {
	var parseErr error
	proceed.Date, parseErr = time.Parse(TimeFormat, proceed.ViewDate)
	if parseErr != nil {
		return parseErr
	}

	return datastore.SaveStruct(proceed, c)
}

func QueryProceed(c appengine.Context) ([] *datastore.Key, []Proceed, error) {
	if c == nil {
		return nil, nil, nil
	}

	q := datastore.NewQuery(proceedKindName).Order("Date")
	if count, err := q.Count(c); err != nil {
		return nil, nil, err
	} else {
		proceeds := make([]Proceed, 0, count)
		keys, getErr := q.GetAll(c, &proceeds)
		if getErr != nil {
			return nil, nil, getErr
		}

		return keys, proceeds, nil
	}

	return nil, nil, nil
}

func(p *Proceed) PutNew(c appengine.Context) (*Proceed, *datastore.Key, error) {

	key := datastore.NewIncompleteKey(c, proceedKindName, nil)

	return p.Put(c, key)
}

func(p *Proceed) Put(c appengine.Context, key *datastore.Key) (*Proceed, *datastore.Key, error) {
	if c == nil || key == nil {
		return nil, nil, nil
	}

	if putKey, err := datastore.Put(c, key, p); err != nil {
		return nil, nil, err
	} else {
		return p, putKey, nil
	}

	return nil, nil, nil
}