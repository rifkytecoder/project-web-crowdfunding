package handler

import (
	"net/http"
	"project-campaign/helper"
	"project-campaign/transaction"
	"project-campaign/user"

	"github.com/gin-gonic/gin"
)

// Analisa transaction
// parameter di uri
// tangkap parameter mapping input struct
// panggil service, input struct sebagai parameter
// service, berbekal campaign id bisa panggil repo
// repo mencari data transaction suatu campaign
type transactionHandler struct {
	service transaction.Service
}

func NewTransactionHandler(service transaction.Service) *transactionHandler {
	return &transactionHandler{service}
}

func (h *transactionHandler) GetCampaignTransactions(c *gin.Context) {
	// tangkap input ID
	var input transaction.GetCampaignTransactionsInput

	err := c.ShouldBindUri(&input) // dari input
	if err != nil {
		response := helper.APIResponse("Failed to get campaign's transactions", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	//buat authorization hanya id yg login saat ini bisa akses**
	currentUser := c.MustGet("currentUser").(user.User)
	input.User = currentUser

	transactions, err := h.service.GetTransactionsByCampaignID(input) // dari service
	if err != nil {
		response := helper.APIResponse("Failed to get campaign's transactionsz", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("Campaign's List transactions", http.StatusOK, "success", transaction.FormatCampaignTransactions(transactions)) //dari formatter
	c.JSON(http.StatusOK, response)

}

// todo GetUserTransactions
// handler
// ambil nilai user dari jwt/middleware
// service
// repo => ambil data transaction (preload campaign)
func (h *transactionHandler) GetUserTransactions(c *gin.Context) {
	// mengambil data siapa user yg melakukan request
	currentUser := c.MustGet("currentUser").(user.User)
	userID := currentUser.ID

	transactions, err := h.service.GetTransactionsByUserID(userID)
	if err != nil {
		response := helper.APIResponse("Failed to get user's transactions", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("User's transactions", http.StatusOK, "success", transaction.FormatUserTransactions(transactions))
	c.JSON(http.StatusOK, response)
}

// Midtrans analis
// ada input dari user `amount`
// handler tangkap input terus di-mapping ke input struct
// panggil service buat transaction, manggil sistem midtrans
// panggil repository create new transaction data
func (h *transactionHandler) CreateTransaction(c *gin.Context) {
	// deklarasi var input `untuk menampung nilai request/inputan body json`
	var input transaction.CreateTransactionInput
	// binding data json ke var input
	err := c.ShouldBindJSON(&input)
	// cek error validasi
	if err != nil {
		// validasi
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}
		// response saat validasi error
		response := helper.APIResponse("Failed to create transaction", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	// hanya user yg saat ini login bisa request
	currentUser := c.MustGet("currentUser").(user.User)
	input.User = currentUser

	// service create new data ke database
	newTransaction, err := h.service.CreateTransaction(input)

	// cek response kondisi saat create data error
	if err != nil {
		response := helper.APIResponse("Failed to create transactionX", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	// kondisi response saat new data create success
	response := helper.APIResponse("Success to create transaction", http.StatusOK, "success", transaction.FormatTransaction(newTransaction))
	c.JSON(http.StatusOK, response)
}
