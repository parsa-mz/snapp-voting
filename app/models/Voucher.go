package models

import (
	databases "SnappVotingBack/app"
	"github.com/getsentry/sentry-go"
)

type Voucher struct {
	Id          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	OwnerId     int64  `json:"owner_id"`
	Icon        string `json:"icon"`
	IsNew       bool   `json:"isNew"`
}

func (v Voucher) GetUserVouchers() []Voucher {
	query, err := databases.PostgresDB.Query("SELECT name,description,icon,is_new FROM vouchers WHERE owner_id = $1", v.OwnerId)
	if err != nil {
		sentry.CaptureException(err)
		return nil
	}
	vouchers := make([]Voucher, 0)
	for query.Next() {
		var voucher Voucher
		err = query.Scan(&voucher.Name, &voucher.Description, &voucher.Icon, &voucher.IsNew)
		if err == nil {
			vouchers = append(vouchers, voucher)
		} else {
			sentry.CaptureException(err)
		}
	}
	return vouchers
}
