package models

import (
	databases "SnappVotingBack/app"
	"database/sql"
	"fmt"
	"github.com/getsentry/sentry-go"
	"time"
)

type Voting struct {
	Id                 int64     `json:"id"`
	Name               string    `json:"name"`
	Description        string    `json:"description"`
	WinnerId           int64     `json:"winnerId,omitempty"`
	WinnersVotesCounts int64     `json:"winnersVotesCounts,omitempty"`
	StartedAt          time.Time `json:"startedAt"`
	EndedAt            time.Time `json:"endedAt"`
}

type UserVotes struct {
	Title         string `json:"title"`
	ParticipantId int64  `json:"participantId"`
	VotesCount    int64  `json:"votesCount"`
	IsWinner      bool   `json:"isWinner,omitempty"`
}

// WinnersVotesCountsMap map[ContestantId]VotesCount
var WinnersVotesCountsMap = make(map[int64]int64)

func init() {
	currentTime := time.Now().UTC()
	rows, err := databases.PostgresDB.Query("select COUNT(user_voting.id),voting.id from user_voting inner join voting on voting.winner_id = user_voting.vote_id where ended_at > $1 group by voting.id", currentTime)
	if err != nil {
		if err != sql.ErrNoRows {
			sentry.CaptureException(err)
		}
		return
	}
	for rows.Next() {
		var votesCount int64
		var votingId int64
		err = rows.Scan(&votesCount, &votingId)
		if err != nil {
			sentry.CaptureException(err)
			return
		}
		WinnersVotesCountsMap[votingId] = votesCount
	}
}
func (v *Voting) TableName() string {
	return "voting"
}
func (v *Voting) IsValid() bool {
	currentTime := time.Now().UTC()
	var total int64
	err := databases.PostgresDB.QueryRow("SELECT COUNT(id) FROM voting WHERE id = $1 AND ended_at > $2 AND $3 > started_at", v.Id, currentTime, currentTime).Scan(&total)
	if err != nil || total == 0 {
		if err != nil && err != sql.ErrNoRows {
			sentry.CaptureException(err)
		}
		return false
	}
	return true
}

func (v *Voting) GetLastWinner() *UserVotes {
	var userVote UserVotes
	row := databases.PostgresDB.QueryRow("SELECT participants.id,voting.name,voting_votes_cache.votes FROM voting INNER JOIN participants ON voting.winner_id = participants.id inner join voting_votes_cache on voting_votes_cache.participant_id = participants.id  WHERE voting.winner_id is not null ORDER BY ended_at desc LIMIT 1")
	if row.Err() != nil {
		if row.Err() != sql.ErrNoRows {
			sentry.CaptureException(row.Err())
		}
		return nil
	}
	err := row.
		Scan(&userVote.ParticipantId, &userVote.Title, &userVote.VotesCount)
	if err != nil {
		if err != sql.ErrNoRows {
			sentry.CaptureException(err)
		}

		return nil
	}

	return &userVote
}

func (v *Voting) GetLast() error {
	query := fmt.Sprintf("SELECT id,name,description,winner_id,started_at,ended_at FROM %s WHERE started_at < $1 order by ended_at desc LIMIT 1", v.TableName())
	currentTime := time.Now().UTC()
	row := databases.PostgresDB.QueryRow(query, currentTime)
	if row.Err() != nil {
		if row.Err() != sql.ErrNoRows {
			sentry.CaptureException(row.Err())
		}
		return row.Err()
	}
	var winnerId sql.NullInt64
	err := row.
		Scan(&v.Id, &v.Name, &v.Description, &winnerId, &v.StartedAt, &v.EndedAt)
	if err != nil {
		if err != sql.ErrNoRows {
			sentry.CaptureException(err)
		}
		return err
	}
	if winnerId.Valid {
		v.WinnerId = winnerId.Int64
	}
	return nil
}

func (v *Voting) GetUserVotesParticipants(SnappUserId int64) []UserVotes {
	rows, err := databases.PostgresDB.Query("SELECT user_voting.vote_id, voting_votes_cache.votes,voting.name,voting_votes_cache.is_winner FROM user_voting INNER JOIN voting_votes_cache ON user_voting.vote_id = voting_votes_cache.participant_id inner join voting on voting.id = user_voting.voting_id WHERE user_voting.owner_id = $1 WHERE voting_votes_cache.voting_id = user_voting.voting_id", SnappUserId)
	if err != nil {
		if err != sql.ErrNoRows {
			sentry.CaptureException(err)
		}
		return nil
	}
	userVotes := make([]UserVotes, 0)
	for rows.Next() {
		var userVote UserVotes
		err = rows.Scan(&userVote.ParticipantId, &userVote.VotesCount, &userVote.Title, &userVote.IsWinner)
		if err != nil {
			sentry.CaptureException(err)
			return nil
		}
		userVotes = append(userVotes, userVote)
	}
	return userVotes
}
