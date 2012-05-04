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

type Initial struct {
	StartCash int64 `datasotre:",noindex"`		// 現金：前年繰越／元入金	
	StartBook int64 `datastore:",noindex"`		// 預金：前年繰越／元入金
}

func createInitialKey(c appengine.Context) *datastore.Key {
	return datastore.NewKey(c, "Initial", "initial", 0, nil)
}

func (init *Initial) ReadDatastore(c appengine.Context) error {
	key := createInitialKey(c)

	if err := datastore.Get(c, key, init); err != nil || err != datastore.ErrNoSuchEntity {
		return err
	} else if(err == datastore.ErrNoSuchEntity) {
		init.StartCash = 0
		init.StartBook = 0
	}

	return nil
}

func (init *Initial) Put(c appengine.Context) error {
	key := createInitialKey(c)

	if _, err := datastore.Put(c, key, init); err != nil {
		return err
	}

	return nil
}
