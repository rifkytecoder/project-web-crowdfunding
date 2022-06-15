package campaign

import (
	"errors"
	"fmt"

	"github.com/gosimple/slug"
)

type Service interface {
	GetCampaigns(userID int) ([]Campaign, error)
	GetCampaignByID(userID GetCampaignDetailInput) (Campaign, error) //from input
	CreateCampaign(input CreateCampaignInput) (Campaign, error)      //from input
	UpdateCampaign(inputID GetCampaignDetailInput, inputData CreateCampaignInput) (Campaign, error)
	SaveCampaignImage(input CreateCampaignImageInput, fileLocation string) (CampaignImage, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

// bukan dlm bentuk json jdi tdk mapping json
func (s *service) GetCampaigns(userID int) ([]Campaign, error) {
	// mengecek user_id ada atau tidak
	if userID != 0 {
		campaigns, err := s.repository.FindByUserID(userID)
		if err != nil {
			return campaigns, err
		}
		return campaigns, nil
	}
	// kalau tdk ada user_id nya kembalikan semua campaigns
	campaigns, err := s.repository.FindAll()
	if err != nil {
		return campaigns, err
	}
	return campaigns, nil
}

// mendapatkan detail campaign dengan ID
func (s *service) GetCampaignByID(input GetCampaignDetailInput) (Campaign, error) {
	campaign, err := s.repository.FindByID(input.ID)

	if err != nil {
		return campaign, err
	}

	return campaign, nil
}

// mapping entity campaign to input
func (s *service) CreateCampaign(input CreateCampaignInput) (Campaign, error) {
	campaign := Campaign{}

	campaign.Name = input.Name
	campaign.ShortDescription = input.ShortDescription
	campaign.Description = input.Description
	campaign.Perks = input.Perks
	campaign.GoalAmount = input.GoalAmount
	campaign.UserID = input.User.ID // dibutuhkan user ID nya sja

	slugCandidate := fmt.Sprintf("%s %d", input.Name, input.User.ID)
	campaign.Slug = slug.Make(slugCandidate) // agar bisa unik nama slug mknya gnkan idnya

	// simpan/create data mapping tadi
	newCampaign, err := s.repository.Save(campaign)
	if err != nil {
		return newCampaign, err
	}
	return newCampaign, nil
}

// update campaign- parameter dri id uri dari (input)
func (s *service) UpdateCampaign(inputID GetCampaignDetailInput, inputData CreateCampaignInput) (Campaign, error) {
	// ambil data dari ID uri
	campaign, err := s.repository.FindByID(inputID.ID)
	if err != nil {
		return campaign, err
	}

	// jika bukan ID pemilik campaign jangan update data
	if campaign.UserID != inputData.User.ID {
		return campaign, errors.New("Not an owner of the campaign")
	}

	// mapping data dari input (data baru re- assign)
	campaign.Name = inputData.Name
	campaign.ShortDescription = inputData.ShortDescription
	campaign.Description = inputData.Description
	campaign.Perks = inputData.Perks
	campaign.GoalAmount = inputData.GoalAmount

	// simpan update data yg baru
	updatedCampaign, err := s.repository.Update(campaign)
	if err != nil {
		return updatedCampaign, err
	}
	return updatedCampaign, nil

}

func (s *service) SaveCampaignImage(input CreateCampaignImageInput, fileLocation string) (CampaignImage, error) {
	// ambil data campaign ID // untuk cek siap yg membuat campaign image tersebut
	campaign, err := s.repository.FindByID(input.CampaignID)
	if err != nil {
		return CampaignImage{}, err
	}

	// jika bukan ID pemilik campaign jangan upload data images
	if campaign.UserID != input.User.ID {
		return CampaignImage{}, errors.New("Not an owner of the campaign")
	}

	isPrimary := 0

	// apakah nilai isPrimary true or false
	if input.IsPrimary {
		isPrimary = 1

		_, err := s.repository.MarkAllImagesAsNonPrimary(input.CampaignID)
		if err != nil {
			return CampaignImage{}, err
		}
	}

	// mapping data input campaign image
	campaignImage := CampaignImage{}

	campaignImage.CampaignID = input.CampaignID
	campaignImage.IsPrimary = isPrimary   //input.IsPrimary
	campaignImage.FileName = fileLocation //dari parameter

	// create image
	newCampaignImage, err := s.repository.CreateImage(campaignImage)
	if err != nil {
		return newCampaignImage, err
	}

	return newCampaignImage, nil

}
