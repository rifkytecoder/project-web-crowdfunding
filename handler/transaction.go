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
