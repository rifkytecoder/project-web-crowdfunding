package transaction

import "project-campaign/user"

type GetCampaignTransactionsInput struct {
	ID   int `uri:"id" binding:"required"` //mengirim parameter dgn uri bkn json
	User user.User
}

type CreateTransactionInput struct {
	Amount     int `json:"amount" binding:"required"`
	CampaignID int `json:"campaign_id" binding:"required"`
	User       user.User
}
