package campaign

import "project-campaign/user"

type GetCampaignDetailInput struct {
	ID int `uri:"id" binding:"required"` //mengirim parameter dgn uri bkn json

}

// request backer new campaign `dalam bentuk/body json`
type CreateCampaignInput struct {
	Name             string `json:"name" binding:"required"`
	ShortDescription string `json:"short_description" binding:"required"`
	Description      string `json:"description" binding:"required"`
	GoalAmount       int    `json:"goal_amount" binding:"required"`
	Perks            string `json:"perks" binding:"required"`
	// tdk perlu json krn untk mncri tau siapa usernya dan untuk pembuatn user ID nya untuk slug
	User user.User
}

// upload images `dalam bentuk form`
type CreateCampaignImageInput struct {
	CampaignID int       `form:"campaign_id" binding:"required"`
	IsPrimary  bool      `form:"is_primary"`
	User       user.User // untuk mencari tau siapa user untuk upload campaign image
}
