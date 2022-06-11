package campaign

import "strings"

// data response json format
type CampaignFormatter struct {
	ID               int    `json:"id"`
	UserID           int    `json:"user_id"`
	Name             string `json:"name"`
	ShortDescription string `json:"short_description"`
	ImageURL         string `json:"image_url"`
	GoalAmount       int    `json:"goal_amount"`
	CurrentAmount    int    `json:"current_amount"`
	Slug             string `json:"slug"`
}

// mapping satu campaign formatter
func FormatCampaign(campaign Campaign) CampaignFormatter {
	campaignFormatter := CampaignFormatter{}
	campaignFormatter.ID = campaign.ID
	campaignFormatter.UserID = campaign.UserID
	campaignFormatter.Name = campaign.Name
	campaignFormatter.ShortDescription = campaign.ShortDescription
	campaignFormatter.GoalAmount = campaign.GoalAmount
	campaignFormatter.CurrentAmount = campaign.CurrentAmount
	campaignFormatter.Slug = campaign.Slug
	campaignFormatter.ImageURL = "" //set nilai default jika id kosong

	// jika campaign punya gambar
	if len(campaign.CampaignImages) > 0 {
		campaignFormatter.ImageURL = campaign.CampaignImages[0].FileName
	}

	return campaignFormatter
}

// mapping untuk lebih satu campaign
func FormatCampaigns(campaigns []Campaign) []CampaignFormatter {

	//jika campaign bernilai 0
	// if len(campaigns) == 0 {
	// 	return []CampaignFormatter{} // [] data slice kosong
	// }
	// var campaignsFormatter []CampaignFormatter
	campaignsFormatter := []CampaignFormatter{}

	// setiap perulangn mendpatkan satu campaign dan memasukan ke campaigns
	for _, campaign := range campaigns {
		campaignFormatter := FormatCampaign(campaign)
		campaignsFormatter = append(campaignsFormatter, campaignFormatter)
	}

	return campaignsFormatter
}

// detail campaign formatter
type CampaignDetailFormatter struct {
	ID               int                      `json:"id"`
	Name             string                   `json:"name"`
	ShortDescription string                   `json:"short_description"`
	Description      string                   `json:"description"`
	ImageURL         string                   `json:"image_url"`
	GoalAmount       int                      `json:"goal_amount"`
	CurrentAmount    int                      `json:"current_amount"`
	UserID           int                      `json:"user_id"`
	Slug             string                   `json:"slug"`
	Perks            []string                 `json:"perks"`
	User             CampaignUserFormatter    `json:"user"`   // bentuk struct
	Images           []CampaignImageFormatter `json:"images"` // bentuk slice
}

type CampaignUserFormatter struct {
	Name     string `json:"name"`
	ImageURL string `json:"image_url"`
}

type CampaignImageFormatter struct {
	ImageURL  string `json:"image_url"`
	IsPrimary bool   `json:"is_primary"`
}

// mapping detail campaign formatter
func FormatCampaignDetail(campaign Campaign) CampaignDetailFormatter {
	campaignDetailFormatter := CampaignDetailFormatter{}
	campaignDetailFormatter.ID = campaign.ID
	campaignDetailFormatter.Name = campaign.Name
	campaignDetailFormatter.ShortDescription = campaign.ShortDescription
	campaignDetailFormatter.Description = campaign.Description
	campaignDetailFormatter.GoalAmount = campaign.GoalAmount
	campaignDetailFormatter.CurrentAmount = campaign.CurrentAmount
	campaignDetailFormatter.UserID = campaign.UserID
	campaignDetailFormatter.Slug = campaign.Slug
	campaignDetailFormatter.ImageURL = "" //set nilai default jika id kosong

	// jika campaign punya gambar
	if len(campaign.CampaignImages) > 0 {
		campaignDetailFormatter.ImageURL = campaign.CampaignImages[0].FileName
	}

	// Perks
	var perks []string // perks dibuat slice
	// loop untuk memunculkan semua data perk
	for _, perk := range strings.Split(campaign.Perks, ",") {
		perks = append(perks, strings.TrimSpace(perk)) // menghilangkan spasi pada data perk
	}
	campaignDetailFormatter.Perks = perks

	// membuat instance krna User adalah struct
	user := campaign.User
	// instance campaign user formatter
	campaignUserFormatter := CampaignUserFormatter{}
	campaignUserFormatter.Name = user.Name
	campaignUserFormatter.ImageURL = user.AvatarFileName

	// mapping user ke campaign detail formatter
	campaignDetailFormatter.User = campaignUserFormatter

	// Images
	images := []CampaignImageFormatter{}

	// looping satuan dan masukkan ke slice images
	for _, image := range campaign.CampaignImages {
		// buat instance untuk bisa mengambil data satuan/atribut
		campaignImageFormatter := CampaignImageFormatter{}
		campaignImageFormatter.ImageURL = image.FileName

		// konvert
		isPrimary := false
		if image.IsPrimary == 1 {
			isPrimary = true
		}
		campaignImageFormatter.IsPrimary = isPrimary //image.IsPrimary

		// masukkan satuan image ke slice images
		images = append(images, campaignImageFormatter)

	}

	// mapping images ke campaign detail formatter
	campaignDetailFormatter.Images = images

	return campaignDetailFormatter

}
