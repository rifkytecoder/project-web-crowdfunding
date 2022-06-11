package handler

import (
	"net/http"
	"project-campaign/campaign"
	"project-campaign/helper"
	"strconv"

	"github.com/gin-gonic/gin"
)

// steps campaign
// tangkap parameter di handler
// handler ke service
// service yang menentukan repository mana yg di call
// repository : GetAll, GetUserByID
// db

type campaignHandler struct {
	service campaign.Service
}

func NewCampaignHandler(service campaign.Service) *campaignHandler {
	return &campaignHandler{service}
}

// api/v1/campaigns
func (h *campaignHandler) GetCampaigns(c *gin.Context) {
	// query path parameter
	userID, _ := strconv.Atoi(c.Query("user_id")) // convert ke int

	campaigns, err := h.service.GetCampaigns(userID)
	if err != nil {

		response := helper.APIResponse("Error to get campaigns", http.StatusBadRequest, "Error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("List of campaigns", http.StatusOK, "success", campaign.FormatCampaigns(campaigns))
	c.JSON(http.StatusOK, response)

}

// membuat detail campaign
func (h *campaignHandler) GetCampaign(c *gin.Context) {
	//api/v1/campaigns/1
	//handler : mapping id yg di url ke struct input => service, call formatter
	//service : inputnya struct input => menangkap id di ulr, call repository
	//repository : get campaign by id

	var input campaign.GetCampaignDetailInput //from input

	err := c.ShouldBindUri(&input) //binding uri

	if err != nil {
		response := helper.APIResponse("Failed to get detail of campaign", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	campaignDetail, err := h.service.GetCampaignByID(input)
	if err != nil {
		response := helper.APIResponse("Failed to get detail of campaign", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("Campaign detail", http.StatusOK, "success", campaign.FormatCampaignDetail(campaignDetail)) //campaignDetail
	c.JSON(http.StatusOK, response)
}
