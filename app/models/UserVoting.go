package models

import (
	databases "SnappVotingBack/app"
	"database/sql"
	"github.com/getsentry/sentry-go"
)

type UserVoting struct {
	Id       int64 `json:"id"`
	VotingId int64 `json:"voting_id"`
	OwnerId  int64 `json:"owner"`
	VoteId   int64 `json:"vote"`
}

func (u *UserVoting) SubmitVote() bool {
	//var total
	//databases.PostgresDB.QueryRow("SELECT COUNT(*) FROM user_voting WHERE voting_id = $1 AND owner = $2", u.VotingId, u.Owner).Scan(&u.Vote)
	var totalCount int64
	err := databases.PostgresDB.QueryRow("SELECT COUNT(id) FROM user_voting WHERE voting_id = $1 AND owner_id= $2", u.VotingId, u.OwnerId).Scan(&totalCount)
	if err != nil {
		if err != sql.ErrNoRows {
			return true
		}
		sentry.CaptureException(err)
	}
	if totalCount > 0 {
		return true
	}
	_, err = databases.PostgresDB.Query("INSERT INTO user_voting (owner_id, voting_id,vote_id) VALUES ($1,$2,$3)", u.OwnerId, u.VotingId, u.VoteId)
	if err != nil {
		sentry.CaptureException(err)
		return false
	}

	return false
}

func (u *UserVoting) GetUserVoteCounts() int64 {
	var totalCount int64
	err := databases.PostgresDB.QueryRow("SELECT COUNT(id) FROM user_voting WHERE voting_id = $1 AND vote_id= $2", u.VotingId, u.VoteId).Scan(&totalCount)
	if err != nil {
		if err != sql.ErrNoRows {
			return 0
		}
		sentry.CaptureException(err)
	}
	return totalCount
}
