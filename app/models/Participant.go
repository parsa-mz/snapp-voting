package models

import (
	databases "SnappVotingBack/app"
	"fmt"
	"github.com/getsentry/sentry-go"
)

type Participant struct {
	Id        int64  `json:"id"`
	Name      string `json:"name"`
	Photo     string `json:"photoUrl"`
	Code      string `json:"code,omitempty"`
	IsActive  bool   `json:"isActive"`
	VoteCount int64  `json:"voteCount,omitempty"`
	MentorId  int64  `json:"mentorId,omitempty"`
}

func (c *Participant) TableName() string {
	return "participants"
}
func (c *Participant) IsValid() bool {
	var total int64
	query := fmt.Sprintf("SELECT COUNT(id) FROM %s WHERE id = $1 and is_active = true", c.TableName())
	err := databases.PostgresDB.QueryRow(query, c.Id).Scan(&total)
	if err != nil || total == 0 {
		if err != nil {
			print(err.Error())
		}
		return false
	}
	return true
}

func (c *Participant) All() []Participant {
	query := fmt.Sprintf("SELECT id,name,photo,code,is_active,mentor_id FROM %s", c.TableName())
	rows, err := databases.PostgresDB.Query(query)
	if err != nil {
		sentry.CaptureException(err)
		return nil

	}
	contestants := make([]Participant, 0)
	for rows.Next() {
		var contestant Participant
		err = rows.Scan(&contestant.Id, &contestant.Name, &contestant.Photo, &contestant.Code, &contestant.IsActive, &contestant.MentorId)
		if err == nil {
			contestants = append(contestants, contestant)
		} else {
			sentry.CaptureException(err)
		}
	}
	return contestants
}
