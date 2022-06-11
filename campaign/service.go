package campaign

type Service interface {
	GetCampaigns(userID int) ([]Campaign, error)
	GetCampaignByID(userID GetCampaignDetailInput) (Campaign, error) //from input
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
