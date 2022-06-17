package transaction

import (
	"errors"
	"project-campaign/campaign"
)

type Service interface {
	GetTransactionsByCampaignID(input GetCampaignTransactionsInput) ([]Transaction, error) //campaignID int = parameternya dibungkus input struct
}

type service struct {
	repository         Repository
	campaignRepository campaign.Repository // get campaign repository **
}

func NewService(repository Repository, campaignRepository campaign.Repository) *service {
	return &service{repository, campaignRepository}
}

func (s *service) GetTransactionsByCampaignID(input GetCampaignTransactionsInput) ([]Transaction, error) {
	// get campaign id**
	campaign, err := s.campaignRepository.FindByID(input.ID)
	if err != nil {
		return []Transaction{}, err
	}
	// check campaign user_id**
	if campaign.UserID != input.User.ID {
		return []Transaction{}, errors.New("Not an owner of the campaign")
	}

	transactions, err := s.repository.GetByCampaignID(input.ID) //input.ID dari atribut input struct
	if err != nil {
		return transactions, err
	}
	return transactions, nil
}
