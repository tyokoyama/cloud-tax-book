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
)

type Type struct {
	Name string `datasotre:",noindex"`		// 勘定科目名
	Detail string `datasotre:",noindex"`	// 勘定科目の説明
	IsExpense bool `datastore:",noindex"`	// 経費かどうか（true：経費／false：経費ではない）
}

func QueryType(c appengine.Context) ([]*datastore.Key, []Type, error) {
	if c == nil {
		return nil, nil, nil
	}

	q := datastore.NewQuery("Type")
	if count, err := q.Count(c); err != nil {
		return nil, nil, err
	} else {
		settings := make([]Type, 0, count)
		if keys, geterr := q.GetAll(c, &settings); geterr != nil {
			return nil, nil, geterr
		} else {
			return keys, settings, nil
		}
	}

	return nil, nil, nil
}

func GetType(c appengine.Context, key *datastore.Key) (*Type, error) {
	if c == nil || key == nil {
		return nil, nil
	}

	var setting Type
	if err := datastore.Get(c, key, &setting); err != nil && err != datastore.ErrNoSuchEntity {
		return nil, err
	} else if err == datastore.ErrNoSuchEntity {
		return nil, nil
	}

	return &setting, nil
}

func (setting *Type) PutNew(c appengine.Context) (*Type, error) {
	if c == nil {
		return nil, nil
	}

	key := datastore.NewIncompleteKey(c, "Type", nil)

	return setting.Put(c, key)

}

func (setting *Type) Put(c appengine.Context, key *datastore.Key) (*Type, error) {
	if c == nil || key == nil {
		return nil, nil
	}

	if _, err := datastore.Put(c, key, setting); err != nil {
		return nil, err
	} else {
		return setting, nil
	}

	return nil, nil

}

