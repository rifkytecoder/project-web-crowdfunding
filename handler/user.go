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

	// declare mapping request body from insomnia(client)
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
