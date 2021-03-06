package handler

import (
	"bwastartup/app/user"
	"bwastartup/auth"
	"bwastartup/helper"
	"crypto/rand"
	"fmt"
	"math/big"
	"net/http"

	"github.com/gin-gonic/gin"
)

type userHandler struct {
	userService user.Service
	authService auth.Service
}

func NewUserHandler(userService user.Service, authService auth.Service) *userHandler{
	return &userHandler{userService, authService}
}

func (h *userHandler) RegisterUser(c *gin.Context){
	// tangkap input dari user
	var input user.RegisterUserInput

	err := c.ShouldBindJSON(&input)
	if err != nil{
		// handler error validation
		errors := helper.FormatValidationError(err)

		// masukan var errors ke dalam object
		errorMessage := gin.H{"errors": errors}

		// response
		response := helper.APIResponse("Register account failed", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}
	
	// mapping input dari user ke struct RegisterUserInput
	// save input user ke service
	// passing struct sebagai parameter service
	newUser, err := h.userService.RegisterUser(input)
	if err != nil{
		// error message
		errorMessage := gin.H{"errors": err.Error()}
		// response
		response := helper.APIResponse("Register account failed", http.StatusBadGateway, "failed", errorMessage)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	token, err := h.authService.GenerateToken(newUser.ID)
	if err != nil{
		// error message
		errorMessage := gin.H{"errors": err.Error()}
		// response
		response := helper.APIResponse("Register account failed", http.StatusBadGateway, "failed", errorMessage)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// melakukan format user
	formatUser := user.FormatUser(newUser, token)

	// format response API
	response := helper.APIResponse("Account has been registered", http.StatusOK, "success", formatUser)
	// response
	c.JSON(http.StatusOK, response)
}

func (h *userHandler) Login(c *gin.Context){
	// tampung input dari user
	var input user.LoginUserInput

	// tangkap input dari user lalu masukan ke input
	err := c.ShouldBindJSON(&input)
	// handle error input user
	if err != nil {
		// validation err message
		errors := helper.FormatValidationError(err)
		// masukan error message ke dalam objek errors
		errorMessage := gin.H{"errors": errors}
		// response
		response := helper.APIResponse("Login failed", http.StatusUnprocessableEntity, "failed", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	// mapping input dari user ke struct RegisterUserInput
	// save input user ke service
	// passing struct sebagai parameter service
	loginUser, err := h.userService.Login(input)
	// handle error login
	if err != nil{
		// error message
		errorMessage := gin.H{"errors": err.Error()}
		// response
		response := helper.APIResponse("Login failed", http.StatusBadRequest, "failed", errorMessage)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	token, err := h.authService.GenerateToken(loginUser.ID)
	if err != nil{
		// error message
		errorMessage := gin.H{"errors": err.Error()}
		// response
		response := helper.APIResponse("Login failed", http.StatusBadRequest, "failed", errorMessage)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	formatter := user.FormatUser(loginUser, token)

	response := helper.APIResponse("Successfuly loggin", http.StatusOK, "success", formatter)

	c.JSON(http.StatusOK, response)
}

func (h *userHandler) IsEmailAvailable(c *gin.Context){
	var input user.CheckEmailInput

	// ambil email dari user
	err := c.ShouldBindJSON(&input)
	if err != nil {
		// handler error validation
		errors := helper.FormatValidationError(err)

		// masukan var errors ke dalam object
		errorMessage := gin.H{"errors": errors}

		// response
		response := helper.APIResponse("Email checking failed", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	isEmailAvailable, err := h.userService.IsEmailAvailable(input)
	if err != nil {
		// error message
		errorMessage := gin.H{"errors": err.Error()}
		// response
		response := helper.APIResponse("Email checking failed", http.StatusBadGateway, "failed", errorMessage)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	data := gin.H{
		"is_available": isEmailAvailable,
	}

	metaMessage := "Email has been registerd"

	if isEmailAvailable {
		metaMessage = "Email is available"
	}

	response := helper.APIResponse(metaMessage, http.StatusOK, "success", data)
	c.JSON(http.StatusOK, response)
}

func (h *userHandler) UploadAvatar(c *gin.Context){
	// ambil input dari user
	// simpan gambarnya di folder "images/"
	// di service panggil reponya
	// JWT (sementara hardcode, seakan2 user yang login)
	// repo ambil data dari user yang ID = 1
	// repo update data user simpan lokasi file

	// ambil avatar gambar dari user
	file, err := c.FormFile("avatar")
	if err != nil {
		data := gin.H{"is_uploaded": false}
		response := helper.APIResponse("Failed to upload avatar image", http.StatusBadRequest, "error", data)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// ambil ID user dari middlware yang sudah di set
	// lalu ubah ke user ".(user.User)"
	currentUser := c.MustGet("currentUser").(user.User)
	userID := currentUser.ID

	// set lokasi avatar
	randomCrypto, _ := rand.Int(rand.Reader, big.NewInt(9999999999))

	// gabungkan beberapa menjadi string
	path := fmt.Sprintf("images/user/%d-%v-%s", userID, randomCrypto, file.Filename)

	// upload file dari user ke folder images
	err = c.SaveUploadedFile(file, path)
	if err != nil {
		data := gin.H{"is_uploaded": false}
		response := helper.APIResponse("Failed to upload avatar image", http.StatusBadRequest, "error", data)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// save avatar ke db berdasarkan userID
	_, err = h.userService.SaveAvatar(userID, path)
	if err != nil {
		data := gin.H{"is_uploaded": false}
		response := helper.APIResponse("Failed to upload avatar image", http.StatusBadRequest, "error", data)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	data := gin.H{"is_uploaded": true}
	response := helper.APIResponse("Avatar successfuly uploaded", http.StatusOK, "success", data)
	c.JSON(http.StatusOK, response)
}

func (h *userHandler) FetchUser(c *gin.Context){
	// collect data user from login
	currentUser := c.MustGet("currentUser").(user.User)
	formatter := user.FormatUser(currentUser, "")

	response := helper.APIResponse("Successfuly fetch data user", http.StatusOK, "success", formatter)
	c.JSON(http.StatusOK, response)
}