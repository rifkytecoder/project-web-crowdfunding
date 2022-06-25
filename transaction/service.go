package transaction

import (
	"errors"
	"project-campaign/campaign"
	"project-campaign/payment"
)

type Service interface {
	GetTransactionsByCampaignID(input GetCampaignTransactionsInput) ([]Transaction, error) //campaignID int = parameternya dibungkus input struct
	// dapat id tdk dari user langsung tpi dari jwt siapa yg melakukan request `mknya tdk pke input `
	GetTransactionsByUserID(userID int) ([]Transaction, error)
	CreateTransaction(input CreateTransactionInput) (Transaction, error)
}

type service struct {
	repository         Repository
	campaignRepository campaign.Repository // get campaign repository **
	paymentService     payment.Service     // midtrans **
}

func NewService(repository Repository, campaignRepository campaign.Repository, paymentService payment.Service) *service {
	return &service{repository, campaignRepository, paymentService}
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

func (s *service) GetTransactionsByUserID(userID int) ([]Transaction, error) {

	transactions, err := s.repository.GetByUserID(userID)

	if err != nil {
		return transactions, err
	}

	return transactions, nil
}

// todo create payment transaction
func (s *service) CreateTransaction(input CreateTransactionInput) (Transaction, error) {
	// object transaction
	transaction := Transaction{}
	// mapping data ke input struct
	transaction.CampaignID = input.CampaignID
	transaction.Amount = input.Amount
	transaction.UserID = input.User.ID
	transaction.Status = "pending"
	//transaction.Code = "Code-Unik-001"

	// simpan new data transaksi
	newTransaction, err := s.repository.Save(transaction)
	if err != nil {
		return newTransaction, err
	}
	// mapping newTransaction ke paymentTransaction `krna error cycle
	paymentTransaction := payment.Transaction{
		ID:     newTransaction.ID,
		Amount: newTransaction.Amount,
	}

	// mendapatkan data payment url
	paymentURL, err := s.paymentService.GetPaymentURL(paymentTransaction, input.User)
	if err != nil {
		return newTransaction, err
	}
	// mapping data payment url
	newTransaction.PaymentURL = paymentURL

	// memasukkan/update data payment url ke field transaction
	newTransaction, err = s.repository.Update(newTransaction)
	if err != nil {
		return newTransaction, err
	}

	return newTransaction, nil

}
