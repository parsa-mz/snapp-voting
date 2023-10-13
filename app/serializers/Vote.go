package serializers

import (
	"SnappVotingBack/app/models"
	"strconv"
)

type UserHistory struct {
	Vouchers []models.Voucher   `json:"vouchers"`
	Voted    []models.UserVotes `json:"voted"`
}
type Vote struct {
	Banner       *models.Banner       `json:"banner,omitempty"`
	Participants []models.Participant `json:"participants"`
	Mentors      []models.Mentor      `json:"mentors"`
	Voting       models.Voting        `json:"voting"`
	History      UserHistory          `json:"history"`
	LastWinner   *models.UserVotes    `json:"lastWinner"`
}

type VoteRequest struct {
	VotingId   int64 `json:"votingId"`
	VoteId     int64 `json:"voteId"`
	userVoting models.UserVoting
	voting     models.Voting
	contestant models.Participant
}

func (v *VoteRequest) Validate(votingIdString, voteIdString string) (Base, bool) {
	var err error
	v.VotingId, err = strconv.ParseInt(votingIdString, 10, 64)
	if err != nil {
		return Base{
			Code:    InvalidInput,
			Message: "voting id must be integer",
		}, false
	}
	v.VoteId, err = strconv.ParseInt(voteIdString, 10, 64)
	if err != nil {
		return Base{
			Code:    InvalidInput,
			Message: "vote id must be integer",
		}, false
	}
	v.voting.Id = v.VotingId

	if !v.voting.IsValid() {
		return Base{
			Code:    InvalidInput,
			Message: "voting id is invalid",
		}, false
	}
	v.contestant.Id = v.VoteId
	if !v.contestant.IsValid() {

		return Base{
			Code:    InvalidInput,
			Message: "vote id is invalid",
		}, false
	}
	return Base{}, true
}
