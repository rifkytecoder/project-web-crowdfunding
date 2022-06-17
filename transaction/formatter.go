package transaction

import "time"

// response JSON campaign transaction format
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
