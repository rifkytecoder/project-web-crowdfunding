package handler

import (
	"fmt"
	"net/http"
	"project-campaign/auth"
	"project-campaign/helper"
	"project-campaign/user"

	"github.com/gin-gonic/gin"
)

type userHandler struct {
	userService user.Service
	authService auth.Service //tambah jika service jwt sdh ada **
}

func NewUserHandler(userService user.Service, authService auth.Service) *userHandler {
	return &userHandler{userService, authService}
}

func (h *userHandler) RegisterUser(c *gin.Context) {

	// declare mapping binding request body JSON from insomnia(client)
	var input user.RegisterUserInput
	err := c.ShouldBindJSON(&input)
	if err != nil {

		// message validation errors
		// var errors []string
		// for _, e := range err.(validator.ValidationErrors) {
		// 	errors = append(errors, e.Error())
		// }
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Register account failed", http.StatusUnprocessableEntity, "Error", errorMessage) // strings.Split(err.Error(), "\n")
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	// input request to service
	newUser, err := h.userService.RegisterUser(input)
	if err != nil {
		response := helper.APIResponse("Register account failed", http.StatusBadRequest, "Error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// token, err := h.jwtService.GenerateToken() **
	token, err := h.authService.GenerateToken(newUser.ID)
	if err != nil {
		response := helper.APIResponse("Register account failed", http.StatusBadRequest, "Error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// response json format user
	formatter := user.FormatUser(newUser, token) //"tokentokentoken" **

	// has helper meta response
	response := helper.APIResponse("Account has been registered", http.StatusOK, "Success", formatter)

	// show response body
	c.JSON(http.StatusOK, response)
}

func (h *userHandler) Login(c *gin.Context) {

	var input user.LoginInput

	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Login failed", http.StatusUnprocessableEntity, "Error", errorMessage) // strings.Split(err.Error(), "\n")
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}
	// logged in user session
	loggedinUser, err := h.userService.Login(input)
	if err != nil {
		errorMessage := gin.H{"errors": err.Error()}

		response := helper.APIResponse("Login failed", http.StatusUnprocessableEntity, "Error", errorMessage) // strings.Split(err.Error(), "\n")
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	// add Token **
	token, err := h.authService.GenerateToken(loggedinUser.ID)
	if err != nil {
		response := helper.APIResponse("Login failed", http.StatusBadRequest, "Error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	formatter := user.FormatUser(loggedinUser, token) //"tokentokentoken" **

	response := helper.APIResponse("Successfully Login", http.StatusOK, "Success", formatter)

	c.JSON(http.StatusOK, response)
}

func (h *userHandler) CheckEmailAvailability(c *gin.Context) {

	var input user.CheckEmailInput

	err := c.ShouldBindJSON(&input)

	if err != nil {
		errors := helper.FormatValidationError(err)

		errorMessage := gin.H{"errors": errors}
		response := helper.APIResponse("Email checkings failed", http.StatusUnprocessableEntity, "Error", errorMessage) // strings.Split(err.Error(), "\n")
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	isEmailAvailable, err := h.userService.IsEmailAvailable(input)
	if err != nil {
		errorMessage := gin.H{"errors": "Server error"}
		response := helper.APIResponse("Email checkings failed", http.StatusUnprocessableEntity, "Error", errorMessage) // strings.Split(err.Error(), "\n")
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	data := gin.H{
		"is_available": isEmailAvailable,
	}

	// meta response has email is available
	metaMessage := "Email has been registered"
	if isEmailAvailable {
		metaMessage = "Email is available"
	}

	response := helper.APIResponse(metaMessage, http.StatusOK, "Success", data) // strings.Split(err.Error(), "\n")
	c.JSON(http.StatusOK, response)

	// ada input email dari user
	// input email di mapping ke struct input
	// struct input di passing ke service
	// service akan panggil repository - email sudah ada atau belum
	// repository - db
}

func (h *userHandler) UploadAvatar(c *gin.Context) {

	file, err := c.FormFile("avatar") // key form-data
	if err != nil {
		data := gin.H{
			"is_uploaded": false,
		}
		response := helper.APIResponse("Failed to upload avatar image", http.StatusBadRequest, "error", data)

		c.JSON(http.StatusBadRequest, response)
		return
	}

	// harusnya dapat dari JWT
	userID := 1

	// destinations
	//path := "images/" + file.Filename //images/file-name.png
	path := fmt.Sprintf("images/%d-%s", userID, file.Filename) //images/1-file-name.png

	err = c.SaveUploadedFile(file, path)
	if err != nil {
		data := gin.H{
			"is_uploaded": false,
		}
		response := helper.APIResponse("Failed to upload avatar image", http.StatusBadRequest, "error", data)

		c.JSON(http.StatusBadRequest, response)
		return
	}

	// harusnya dapat dari JWT
	//userID := 1

	_, err = h.userService.SaveAvatar(userID, path)
	if err != nil {
		data := gin.H{
			"is_uploaded": false,
		}
		response := helper.APIResponse("Failed to upload avatar image", http.StatusBadRequest, "error", data)

		c.JSON(http.StatusBadRequest, response)
		return
	}

	data := gin.H{
		"is_uploaded": true,
	}
	response := helper.APIResponse("Avatar successfully uploaded", http.StatusOK, "success", data)

	c.JSON(http.StatusOK, response)

	// input dari user
	// simpan gambarnya di folder "images/"
	// di service kita panggil repo
	// JWT (hardcode dummy user login ID = 1)
	//  repo ambil data user yg di ID = 1
	//  repo update data user simpan lokasi file
}
