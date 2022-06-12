package campaign

import "project-campaign/user"

type GetCampaignDetailInput struct {
	ID int `uri:"id" binding:"required"` //mengirim parameter dgn uri bkn json

}

// request backer new campaign
type CreateCampaignInput struct {
	Name             string `json:"name" binding:"required"`
	ShortDescription string `json:"short_description" binding:"required"`
	Description      string `json:"description" binding:"required"`
	GoalAmount       int    `json:"goal_amount" binding:"required"`
	Perks            string `json:"perks" binding:"required"`
	// tdk perlu json krn untk mncri tau siapa usernya dan untuk pembuatn user ID nya untuk slug
	User user.User
}
