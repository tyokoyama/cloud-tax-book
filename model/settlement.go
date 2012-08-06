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

type Settlement struct {
	TypeId int64 `datastore:",noindex"` 	// 勘定科目種別
	Name string `datastore:",noindex"`		// 勘定科目名
	IsExpense bool `datastore:",noindex"`	// 経費かどうか（true：経費／false：経費ではない）
	MoneyIn int64 `datastore:",noindex"`	// 収入
	MoneyOut int64 `datastore:",noindex"`	// 支出
}

type SettlementSummary struct {
	MoneyIn int64 `datastore:",noindex"`
	MoneyOut int64 `datastore:",noindex"`
	StartCash int64 `datastore:",noindex"`		// 現金：前年繰越／元入金	
	StartBook int64 `datastore:",noindex"`		// 預金：前年繰越／元入金
	EndCash int64 `datastore:",noindex"`
	EndBook int64 `datastore:",noindex"`
}

func(summary *SettlementSummary) Get(c appengine.Context) error {
	key := datastore.NewKey(c, "settlementsummary", "summary", 0, nil)
	if err := datastore.Get(c, key, summary); err != nil && err != datastore.ErrNoSuchEntity {
		return err
	}

	return nil
}

func(summary SettlementSummary) Put(c appengine.Context) error {
	key := datastore.NewKey(c, "settlementsummary", "summary", 0, nil)
	if _, err := datastore.Put(c, key, &summary); err != nil {
		return err
	}

	return nil
}

func CalcSummary(c appengine.Context) error {
	// Read Initial
	var init Initial
	if err := init.ReadDatastore(c); err != nil {
		c.Errorf("Read Initial Error: %s", err.Error())
		return err
	}

	var settlementsummary SettlementSummary
	if err := settlementsummary.Get(c); err != nil {
		c.Errorf("SettlementSummary.Get Error %s", err.Error())
		return err
	}
	settlementsummary.StartCash = init.StartCash
	settlementsummary.StartBook = init.StartBook

	// Initialize MoneyIn, MoneyOut
	settlementsummary.MoneyIn = 0
	settlementsummary.MoneyOut = 0

	// Read Expense
	if expenses, err := QueryExpense(c); err != nil {
		c.Errorf("Read Expense Error: %s", err.Error())
		return err
	} else {
		var expenseMap = make(map[int64]Settlement)
		for _, expense := range expenses {
			c.Debugf("%d", expense.MoneyOut)

			var settlement Settlement
			settlement.TypeId = expense.Type
			settlement.Name = expense.TypeName
			settlement.IsExpense = true
			settlement.MoneyIn = expense.MoneySalesIn + expense.MoneyIn
			settlement.MoneyOut = expense.MoneySalesOut + expense.MoneyOut

			summary, ok := expenseMap[settlement.TypeId]
			if ok {
				summary.MoneyIn += settlement.MoneyIn
				summary.MoneyOut += settlement.MoneyOut
				expenseMap[summary.TypeId] = summary
			} else {
				expenseMap[settlement.TypeId] = settlement
			}
		}

		expenseKeys := make([]*datastore.Key, len(expenseMap))
		expenseData := make([]Settlement, len(expenseMap))
		index := 0
		for _, settlement := range expenseMap {
			c.Debugf("%s = %d, %d", settlement.Name, settlement.MoneyIn, settlement.MoneyOut)
			settlementsummary.MoneyIn += settlement.MoneyIn
			settlementsummary.MoneyOut += settlement.MoneyOut

			expenseKeys[index] = datastore.NewKey(c, "Settlement", "", settlement.TypeId, nil)
			expenseData[index] = settlement
		}

		// PUT to Datastore
		if _, err := datastore.PutMulti(c, expenseKeys, expenseData); err != nil {
			c.Errorf("Expense PutMulti Error: %s", err.Error())
			return err			
		}
	}

	// Read NonExpense
	if nonexpenses, err := QueryNonExpense(c); err != nil {
		c.Errorf("Read NonExpense Error: %s", err.Error())
		return err
	} else {
		var nonexpenseMap = make(map[int64]Settlement)
		for _, nonexpense := range nonexpenses {
			c.Debugf("%d", nonexpense.MoneyOut)

			var settlement Settlement
			settlement.TypeId = nonexpense.Type
			settlement.Name = nonexpense.TypeName
			settlement.IsExpense = false
			settlement.MoneyIn = nonexpense.MoneySalesIn + nonexpense.MoneyIn
			settlement.MoneyOut = nonexpense.MoneySalesOut + nonexpense.MoneyOut

			summary, ok := nonexpenseMap[settlement.TypeId]
			if ok {
				summary.MoneyIn += settlement.MoneyIn
				summary.MoneyOut += settlement.MoneyOut
				nonexpenseMap[summary.TypeId] = summary
			} else {
				nonexpenseMap[settlement.TypeId] = settlement
			}
		}

		expenseKeys := make([]*datastore.Key, len(nonexpenseMap))
		expenseData := make([]Settlement, len(nonexpenseMap))
		index := 0
		for _, settlement := range nonexpenseMap {
			c.Debugf("%s = %d, %d", settlement.Name, settlement.MoneyIn, settlement.MoneyOut)
			settlementsummary.MoneyIn += settlement.MoneyIn
			settlementsummary.MoneyOut += settlement.MoneyOut

			expenseKeys[index] = datastore.NewKey(c, "Settlement", "", settlement.TypeId, nil)
			expenseData[index] = settlement
		}

		// PUT to Datastore
		if _, err := datastore.PutMulti(c, expenseKeys, expenseData); err != nil {
			c.Errorf("NonExpense PutMulti Error: %s", err.Error())
			return err			
		}
	}

	c.Debugf("SettlementSummary = %d, %d", ( settlementsummary.MoneyOut + settlementsummary.EndCash + settlementsummary.EndBook), ( settlementsummary.MoneyIn + settlementsummary.StartCash + settlementsummary.StartBook))
	// Put SettlementSummary
	if err := settlementsummary.Put(c); err != nil {
		c.Errorf("SettlementSummary.Put Error %s", err.Error())
		return err
	}

	return nil
}