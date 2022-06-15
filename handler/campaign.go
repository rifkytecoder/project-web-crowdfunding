package handler

import (
	"fmt"
	"net/http"
	"project-campaign/campaign"
	"project-campaign/helper"
	"project-campaign/user"
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

// Create campaign analisa
// tangkap parameter dari user mapping ke input struct
// ambil current user dari jwt/handler
// panggil service, parameter input struct (dan juga buat slug)
// panggil repository untuk simpan data campaign baru
func (h *campaignHandler) CreateCampaign(c *gin.Context) {
	// variabel untuk menangkap parameter yg di kirim dari user
	var input campaign.CreateCampaignInput

	err := c.ShouldBindJSON(&input)

	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Failed to create campaign", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	// mengambil data user dri kontex yg tlh di buat/set di middleware
	currentUser := c.MustGet("currentUser").(user.User)
	// memasukkan data user
	input.User = currentUser

	newCampaign, err := h.service.CreateCampaign(input)
	if err != nil {
		response := helper.APIResponse("Failed to create campaign", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("Success to create campaign", http.StatusOK, "success", campaign.FormatCampaignDetail(newCampaign)) //newCampaign
	c.JSON(http.StatusOK, response)

}

//Analisa update campaign
// user masukkan input
// handler
// mapping input ke input struct (ada 2)
// input dari user, dan juga input yang ada di uri (passing ke service)
// service (find campaign by id, mengkp parameter yg sdh ada ke struct)
// repository update data campaign
func (h *campaignHandler) UpdatedCampaign(c *gin.Context) {
	var inputID campaign.GetCampaignDetailInput //from input

	err := c.ShouldBindUri(&inputID) //binding uri :id (id params)

	if err != nil {
		response := helper.APIResponse("Failed to update campaign", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	var inputData campaign.CreateCampaignInput

	err = c.ShouldBindJSON(&inputData) //binding json
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Failed to update campaign", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	// agar hanya ID yg memiliki campaign yg bisa update data campaignnya
	// mengambil data user dri kontex yg tlh di buat/set di middleware
	currentUser := c.MustGet("currentUser").(user.User)
	// memasukkan data user
	inputData.User = currentUser

	updatedCampaign, err := h.service.UpdateCampaign(inputID, inputData)
	if err != nil {
		response := helper.APIResponse("Failed to update campaign", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("Success to update campaign", http.StatusOK, "success", campaign.FormatCampaignDetail(updatedCampaign)) //newCampaign
	c.JSON(http.StatusOK, response)
}

// Analis upload campaign image
// handler
// tangkap input dan ubah ke struct input
// save image campaign ke suatu folder
// service (kondisi panggil point 2 di repo, panggil repo point 1)
// repository :
// 1. create image/save data image ke dalam tabel campaign_images
// 2. ubah is_primary true ke false
func (h *campaignHandler) UploadImage(c *gin.Context) {
	var input campaign.CreateCampaignImageInput

	err := c.ShouldBind(&input) // karena berbentuk form

	// validasi
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Failed to update campaign", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	// agar hanya ID yg memiliki campaign yg bisa uplaod data campaignnya
	// mengambil data user dri kontex yg tlh di buat/set di middleware
	currentUser := c.MustGet("currentUser").(user.User)
	// memasukkan data user
	input.User = currentUser
	userID := currentUser.ID

	file, err := c.FormFile("file") // key form-data(multipart insomnia)
	if err != nil {
		data := gin.H{
			"is_uploaded": false,
		}
		response := helper.APIResponse("Failed to upload campaign image", http.StatusBadRequest, "error", data)

		c.JSON(http.StatusBadRequest, response)
		return
	}

	// destinations
	//path := "images/" + file.Filename //images/file-name.png
	path := fmt.Sprintf("images/%d-%s", userID, file.Filename) //<images> foldernya bisa di ganti

	// filenya di upload yg di tangkap dari file
	err = c.SaveUploadedFile(file, path)
	if err != nil {
		data := gin.H{
			"is_uploaded": false,
		}
		response := helper.APIResponse("Failed to upload campaign image", http.StatusBadRequest, "error", data)

		c.JSON(http.StatusBadRequest, response)
		return
	}

	// simpan file ke database dan file path folder images
	_, err = h.service.SaveCampaignImage(input, path)
	if err != nil {
		data := gin.H{
			"is_uploaded": false,
		}
		response := helper.APIResponse("Failed to upload campaign image", http.StatusBadRequest, "error", data)

		c.JSON(http.StatusBadRequest, response)
		return
	}

	data := gin.H{
		"is_uploaded": true,
	}
	response := helper.APIResponse("Campaign image successfully uploaded", http.StatusOK, "success", data)

	c.JSON(http.StatusOK, response)

}
