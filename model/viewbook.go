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
	"fmt"
)

type ViewBook struct {
	Id int64			// 現金データのキーのID
	Type string 		// 費用種別（通信費、事業主貸、事業主借…）
	Date string      	// 登録日付
	Detail string		// 明細
	MoneySalesIn string	// 現金売上入金
	MoneyIn string		// その他入金
	MoneySalesOut string	// 現金仕入
	MoneyOut string		// その他出金
	Balance string		// 残高
}

func (view *ViewBook) Create(data Book, id int64) {
	view.Id = id
	view.Type = "勘定科目"			// まだ未実装
	view.Date = data.Date.Format("2006-01-02")
	view.Detail = data.Detail
	view.MoneySalesIn = fmt.Sprintf("%d", data.MoneySalesIn)
	view.MoneyIn = fmt.Sprintf("%d", data.MoneyIn)
	view.MoneySalesOut = fmt.Sprintf("%d", data.MoneySalesOut)
	view.MoneyOut = fmt.Sprintf("%d", data.MoneyOut)
	view.Balance = fmt.Sprintf("%d", data.Balance)
}