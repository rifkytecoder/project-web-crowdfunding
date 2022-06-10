package campaign

type Service interface {
	FindCampaigns(userID int) ([]Campaign, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

// bukan dlm bentuk json jdi tdk mapping json
func (s *service) FindCampaigns(userID int) ([]Campaign, error) {
	// mengecek user_id ada atau tidak
	if userID != 0 {
		campaigns, err := s.repository.FindByUserID(userID)
		if err != nil {
			return campaigns, err
		}
		return campaigns, nil
	}
	// kalau tdk ada user_id nya kemblikan semua campaigns
	campaigns, err := s.repository.FindAll()
	if err != nil {
		return campaigns, err
	}
	return campaigns, nil
}
