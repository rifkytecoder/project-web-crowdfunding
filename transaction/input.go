package transaction

import "project-campaign/user"

type GetCampaignTransactionsInput struct {
	ID   int `uri:"id" binding:"required"` //mengirim parameter dgn uri bkn json
	User user.User
}
