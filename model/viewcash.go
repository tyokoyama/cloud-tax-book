package model

import (
	"fmt"
)

type ViewCash struct {
	Type string 		// 費用種別（通信費、事業主貸、事業主借…）
	Date string      	// 登録日付
	Detail string		// 明細
	MoneySalesIn string	// 現金売上入金
	MoneyIn string		// その他入金
	MoneySalesOut string	// 現金仕入
	MoneyOut string		// その他出金
	Balance string		// 残高
}

func (view *ViewCash) Create(data Cash) {
	view.Type = "勘定科目"			// まだ未実装
	view.Date = data.Date.Format("2006-01-02")
	view.Detail = data.Detail
	view.MoneySalesIn = fmt.Sprintf("%d", data.MoneySalesIn)
	view.MoneyIn = fmt.Sprintf("%d", data.MoneyIn)
	view.MoneySalesOut = fmt.Sprintf("%d", data.MoneySalesOut)
	view.MoneyOut = fmt.Sprintf("%d", data.MoneyOut)
	view.Balance = fmt.Sprintf("%d", data.Balance)
}