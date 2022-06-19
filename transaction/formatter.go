package transaction

import "time"

// todo response JSON campaign transaction format
type CampaignTransactionFormatter struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Amount    int       `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
}

// format single data campaign transaction
func FormatCampaignTransaction(transaction Transaction) CampaignTransactionFormatter {
	// deklarasi object struct
	formatter := CampaignTransactionFormatter{}

	// mapping attribute
	formatter.ID = transaction.ID
	formatter.Name = transaction.User.Name
	formatter.Amount = transaction.Amount
	formatter.CreatedAt = transaction.CreatedAt
	return formatter
}

// format list [] data campaign transaction
func FormatCampaignTransactions(transactions []Transaction) []CampaignTransactionFormatter {

	// jika tidak ada transaction balikan array [] `kosong`
	if len(transactions) == 0 {
		return []CampaignTransactionFormatter{}
	}

	// variabel slice []CampaignTransactionFormatter
	var transactionFormatter []CampaignTransactionFormatter

	// perulangan untuk menampilkan list variabel slice transactionFormatter
	for _, transaction := range transactions {
		formatter := FormatCampaignTransaction(transaction) //single data
		// menambah setiap single data ke variabel slice[]
		transactionFormatter = append(transactionFormatter, formatter) //list data
	}

	return transactionFormatter

}

// todo format user transaction
type UserTransactionFormatter struct {
	ID        int               `json:"id"`
	Amount    int               `json:"amount"`
	Status    string            `json:"status"`
	CreatedAt time.Time         `json:"created_at"`
	Campaign  CampaignFormatter `json:"campaign"`
}

type CampaignFormatter struct {
	Name     string `json:"name"`
	ImageURL string `json:"image_url"`
}

// format user response json transaction `single object data`
func FormatUserTransaction(transaction Transaction) UserTransactionFormatter {

	formatter := UserTransactionFormatter{}
	formatter.ID = transaction.ID
	formatter.Amount = transaction.Amount
	formatter.Status = transaction.Status
	formatter.CreatedAt = transaction.CreatedAt

	campaignFormatter := CampaignFormatter{}
	campaignFormatter.Name = transaction.Campaign.Name
	campaignFormatter.ImageURL = "" // nilai default jika tdk ada gambar
	if len(transaction.Campaign.CampaignImages) > 0 {
		campaignFormatter.ImageURL = transaction.Campaign.CampaignImages[0].FileName // hirarki
	}

	formatter.Campaign = campaignFormatter

	return formatter
}

// format list [] slice data campaign transaction
func FormatUserTransactions(transactions []Transaction) []UserTransactionFormatter {

	// jika tidak ada transaction balikan array [] `kosong`
	if len(transactions) == 0 {
		return []UserTransactionFormatter{}
	}

	// variabel slice []CampaignTransactionFormatter
	var transactionFormatter []UserTransactionFormatter

	// perulangan untuk menampilkan list variabel slice transactionFormatter
	// setiap perulangn kita mendapatkan single data dari transaction
	for _, transaction := range transactions {
		formatter := FormatUserTransaction(transaction) //single data
		// menambah setiap single data ke variabel slice[]
		transactionFormatter = append(transactionFormatter, formatter) //list data
	}

	return transactionFormatter

}
