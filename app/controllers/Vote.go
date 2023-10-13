package controllers

import (
	"SnappVotingBack/app/models"
	"SnappVotingBack/app/serializers"
	"github.com/gin-gonic/gin"
	"net/http"
)

type VoteController struct {
}

// Vote
// @Summary      Submit user vote
// @Tags         vote
// @Produce      json
// @Param        snapp_id   path      int  true  "Snapp id"
// @Success      200  {object}  serializers.Vote
// @Failure      400  {object}  serializers.Base
// @Failure      404  {object}  serializers.Base
// @Router       /vote/{snapp_id} [get]
func (VoteController) Vote(ctx *gin.Context) {
	var banner models.Banner
	var participant models.Participant
	var voucher models.Voucher
	var mentor models.Mentor
	var voting models.Voting
	voting.GetLastWinner()
	voucher.OwnerId = ctx.GetInt64("snappUser_id")
	err := voting.GetLast()
	if err != nil {
		print(err.Error())
	}
	voucher.GetUserVouchers()
	ctx.JSON(http.StatusOK, serializers.Vote{
		Banner:       banner.GetBanner(),
		Participants: participant.All(),
		Mentors:      mentor.All(),
		Voting:       voting,
		History: serializers.UserHistory{
			Vouchers: voucher.GetUserVouchers(),
			Voted:    voting.GetUserVotesParticipants(voucher.OwnerId),
		},
		LastWinner: voting.GetLastWinner(),
	})
}

// SubmitVote
// @Summary      Submit user vote
// @Tags         vote
// @Produce      json
// @Param        snapp_id   path      int  true  "Snapp id"
// @Param        voting_id   path      int  true  "voting id"
// @Param        vote_id   path      int  true  "vote id"
// @Success      200  {object}  serializers.Vote
// @Failure      400  {object}  serializers.Base
// @Failure      404  {object}  serializers.Base
// @Router       /vote/{snapp_id}/{voting_id}/{vote_id} [post]
func (VoteController) SubmitVote(ctx *gin.Context) {
	votingIdString, voteIdString := ctx.Param("voting_id"), ctx.Param("vote_id")

	request := new(serializers.VoteRequest)

	base, isValid := request.Validate(votingIdString, voteIdString)
	if !isValid {
		ctx.JSON(http.StatusBadRequest, base)
		return
	}

	userVoting := models.UserVoting{
		OwnerId:  ctx.GetInt64("snappUser_id"),
		VotingId: request.VotingId,
		VoteId:   request.VoteId,
	}

	alreadyVoted := userVoting.SubmitVote()
	if alreadyVoted {
		ctx.JSON(http.StatusBadRequest, serializers.Base{
			Code:    serializers.AlreadyVoted,
			Message: "already voted",
		})
		return
	}

	var banner models.Banner
	var participant models.Participant
	var voucher models.Voucher
	var mentor models.Mentor
	var voting models.Voting
	voting.GetLastWinner()
	voucher.OwnerId = ctx.GetInt64("snappUser_id")
	err := voting.GetLast()
	if err != nil {
		print(err.Error())
	}
	voucher.GetUserVouchers()
	ctx.JSON(http.StatusOK, serializers.Vote{
		Banner:       banner.GetBanner(),
		Participants: participant.All(),
		Mentors:      mentor.All(),
		Voting:       voting,
		History: serializers.UserHistory{
			Vouchers: voucher.GetUserVouchers(),
			Voted:    voting.GetUserVotesParticipants(voucher.OwnerId),
		},
		LastWinner: voting.GetLastWinner(),
	})
}
