package transaction

import "project-campaign/user"

type GetCampaignTransactionsInput struct {
	ID   int `uri:"id" binding:"required"` //mengirim parameter dgn uri bkn json
	User user.User
}

type CreateTransactionInput struct {
	Amount     int `json:"amount" binding:"required"`
	CampaignID int `json:"campaign_id" binding:"required"`
	User       user.User
}

// todo menangkap response JSON notification dari midtrans
// mendapatakan notifikasi status setelah user melakukan transfer/transaksi
type TransactionNotificationInput struct {
	// empat data essensial yg di kirim dari midtrans
	TransactionStatus string `json:"transaction_status"`
	OrderID           string `json:"order_id"`
	PaymentType       string `json:"payment_type"`
	FraudStatus       string `json:"fraud_status"`
}
