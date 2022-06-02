package handler

import (
	"net/http"
	"project-campaign/helper"
	"project-campaign/user"

	"github.com/gin-gonic/gin"
)

type userHandler struct {
	userService user.Service
}

func NewUserHandler(userService user.Service) *userHandler {
	return &userHandler{userService}
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

	// token, err := h.jwtService.GenerateToken()

	// response json format user
	formatter := user.FormatUser(newUser, "tokentokentoken")

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

	formatter := user.FormatUser(loggedinUser, "tokentokentoken")

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
