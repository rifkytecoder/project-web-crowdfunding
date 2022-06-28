package payment

import (
	"project-campaign/user"
	"strconv"

	midtrans "github.com/veritrans/go-midtrans"
)

type service struct {
	// transactionRepository transaction.Repository // notification midtrans call repo** //dapatkan
	// campaignRepository    campaign.Repository    // call/dapatkan campaign**
}

type Service interface {
	GetPaymentURL(transaction Transaction, user user.User) (string, error)
	// ProcessPayment(input transaction.TransactionNotificationInput) error // notification midtrans**
}

func NewService() *service {
	return &service{} // notification midtrans call** & campaign
}

func (s *service) GetPaymentURL(transaction Transaction, user user.User) (string, error) {
	midclient := midtrans.NewClient()
	midclient.ServerKey = "SB-Mid-server-jF76yh7Dj1u47w2uA2TV6u01" //"YOUR-VT-SERVER-KEY"
	midclient.ClientKey = "SB-Mid-client-UUGRfsIHH34cUgqy"         //"YOUR-VT-CLIENT-KEY"
	midclient.APIEnvType = midtrans.Sandbox                        // ganti ke Production klo udh launching/delploy

	snapGateway := midtrans.SnapGateway{
		Client: midclient,
	}

	snapReq := &midtrans.SnapReq{
		CustomerDetail: &midtrans.CustDetail{
			Email: user.Email,
			FName: user.Name,
		},
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  strconv.Itoa(transaction.ID),
			GrossAmt: int64(transaction.Amount),
		},
	}

	// get token untuk mendapatkan redirect url
	snapTokenResp, err := snapGateway.GetToken(snapReq)
	if err != nil {
		return "", err
	}

	return snapTokenResp.RedirectURL, nil
}

// // todo tangkap response json midtrans notification transaction
// func (s *service) ProcessPayment(input transaction.TransactionNotificationInput) error {
// 	// var menampung inputan midtrans, parsing response json dari midtrans
// 	transaction_id, _ := strconv.Atoi(input.OrderID) //string to int

// 	// mapping data midtrans ke entitas transaction dgn field ID
// 	transaction, err := s.transactionRepository.GetByID(transaction_id)
// 	if err != nil {
// 		return err
// 	}

// 	// Cek kondisi status <value> dari midtrans transaction notification dan assign status ke field status
// 	if input.PaymentType == "credit_card" && input.TransactionStatus == "capture" && input.FraudStatus == "accept" {
// 		transaction.Status = "paid"
// 	} else if input.TransactionStatus == "settlement" {
// 		transaction.Status = "paid"
// 	} else if input.TransactionStatus == "deny" || input.TransactionStatus == "expire" || input.TransactionStatus == "cancel" {
// 		transaction.Status = "cancelled"
// 	}

// 	// update field Status di entitas transaction
// 	updatedTransaction, err := s.transactionRepository.Update(transaction)
// 	if err != nil {
// 		return err
// 	}

// 	// todo menga-update data campaign field backer_count (menambah jumlah) jika transaction ber-status paid
// 	// dapatkan field campaign_id dari entitas campaign bersangkutan
// 	campaign, err := s.campaignRepository.FindByID(updatedTransaction.CampaignID)
// 	if err != nil {
// 		return err
// 	}

// 	// jika ada kondisi field status paid maka field backer dan amount akan bertambah
// 	if updatedTransaction.Status == "paid" {
// 		campaign.BackerCount = campaign.BackerCount + 1                             // bertambah 1 jika status transaksi paid
// 		campaign.CurrentAmount = campaign.CurrentAmount + updatedTransaction.Amount // bertambah jumlah amount jika status transaksi paid

// 		// simpan dan update data perubahan campaign
// 		_, err := s.campaignRepository.Update(campaign)
// 		if err != nil {
// 			return err
// 		}
// 	}

// 	return nil

// }
