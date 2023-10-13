package models

import (
	databases "SnappVotingBack/app"
	"database/sql"
	"fmt"
	"github.com/getsentry/sentry-go"
)

type Banner struct {
	Id    int64  `json:"id,omitempty"`
	Image string `json:"image"`
	Link  string `json:"link,omitempty"`
}

func (b *Banner) TableName() string {
	return "banners"
}
func (b *Banner) GetBanner() *Banner {
	query := fmt.Sprintf("SELECT image,link FROM %s WHERE is_active = true order by id desc LIMIT 1", b.TableName())
	row := databases.PostgresDB.QueryRow(query)
	if row.Err() != nil {
		if row.Err() != sql.ErrNoRows {
			sentry.CaptureException(row.Err())
		}
		return nil
	}

	err := row.Scan(&b.Image, &b.Link)
	if err != nil {
		sentry.CaptureException(err)
		return nil
	}

	return b
}
