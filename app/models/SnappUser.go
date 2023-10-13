package models

import (
	databases "SnappVotingBack/app"
	"database/sql"
	"github.com/getsentry/sentry-go"
)

type SnappUser struct {
	Id      int64  `json:"id"`
	SnappId string `json:"snapp_id"`
}

func (u *SnappUser) GetUser() (bool, error) {
	var total int64
	err := databases.PostgresDB.QueryRow("SELECT id FROM snapp_users WHERE snapp_id = $1", u.SnappId).Scan(&u.Id)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		sentry.CaptureException(err)
		return false, err
	}
	if total > 0 {
		return true, nil
	}
	return true, err
}
