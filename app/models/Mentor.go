package models

import (
	databases "SnappVotingBack/app"
	"database/sql"
	"fmt"
	"github.com/getsentry/sentry-go"
)

type Mentor struct {
	Id    int64  `json:"id"`
	Name  string `json:"name"`
	Photo string `json:"photo"`
}

func (c *Mentor) TableName() string {
	return "mentors"
}
func (c *Mentor) All() []Mentor {
	query := fmt.Sprintf("SELECT id,name,photo FROM %s", c.TableName())
	rows, err := databases.PostgresDB.Query(query)
	if err != nil {
		sentry.CaptureException(err)
		return nil
	}
	mentors := make([]Mentor, 0)
	for rows.Next() {
		var mentor Mentor
		var PhotoNull sql.NullString
		err = rows.Scan(&mentor.Id, &mentor.Name, &PhotoNull)
		if PhotoNull.Valid {
			mentor.Photo = PhotoNull.String
		}
		if err == nil {
			mentors = append(mentors, mentor)
		} else {
			sentry.CaptureException(err)
		}
	}
	return mentors
}
