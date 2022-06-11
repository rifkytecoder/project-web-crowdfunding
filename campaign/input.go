package campaign

type GetCampaignDetailInput struct {
	ID int `uri:"id" binding:"required"` //mengirim parameter dgn uri bkn json

}
