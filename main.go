package main

import (
	"log"
	"net/http"
	"project-campaign/auth"
	"project-campaign/campaign"
	"project-campaign/handler"
	"project-campaign/helper"
	"project-campaign/user"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	dsn := "root:admin@tcp(127.0.0.1:3306)/demo_db?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal(err.Error())
	}

	//fmt.Println("Connection database succeedd")

	userRepository := user.NewRepository(db)
	campaignRepository := campaign.NewRepository(db) //**
	//testing preload gambar is_primary
	// campaigns, _ := campaignRepository.FindByUserID(1)
	// fmt.Println("debug")
	// fmt.Println("debug")
	// fmt.Println("debug")
	// fmt.Println(len(campaigns)) // rentang/jumlah data
	// // tampilkan rentang data dari field name
	// for _, campaign := range campaigns {
	// 	fmt.Println(campaign.Name)
	// 	if len(campaign.CampaignImages) > 0 {
	// 		fmt.Println("jumlah gambar")
	// 		fmt.Println(len(campaign.CampaignImages))
	// 		fmt.Println(campaign.CampaignImages[0].FileName)
	// 	}
	// }

	userService := user.NewService(userRepository)
	campaignService := campaign.NewService(campaignRepository)
	// userByEmail, err := userRepository.FindByEmail("xgopers@gmail.com")
	// if err != nil {
	// 	fmt.Println(err.Error())
	// }
	// if userByEmail.ID == 0 {
	// 	fmt.Println("User tidak ditemukan")
	// } else {
	// 	fmt.Println(userByEmail.Name)
	// }

	// input := user.LoginInput{
	// 	Email:    "gopers@gmail.com",
	// 	Password: "passwordz",
	// }
	// user, err := userService.Login(input)
	// if err != nil {
	// 	fmt.Println("Terjadi kesalahan")
	// 	fmt.Println(err.Error())
	// }
	// fmt.Println(user.Name)
	// fmt.Println(user.Email)

	//userService.SaveAvatar(1, "images/1-profile.png")

	// testing service load campaign data dengan user_id
	// campaign, _ := campaignService.GetCampaigns(1)
	// fmt.Println(len(campaign))

	authService := auth.NewService()
	//fmt.Println(authService.GenerateToken(1001)) //hasilnya copy ke jwt.io

	// Testing Validation token
	// token, err := authService.ValidateToken("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoyMX0.Mbv-7hj16P7mVtBK3mD_zX9CoTz6yzUhzjuAZxSYP5I")
	// if err != nil {
	// 	fmt.Println("ERROR")
	// }
	// if token.Valid {
	// 	fmt.Println("VALID")
	// } else {
	// 	fmt.Println("INVALID")
	// } // cek di JWT.io di bagian Signature

	userHandler := handler.NewUserHandler(userService, authService)
	campaignHandler := handler.NewCampaignHandler(campaignService)

	router := gin.Default()
	api := router.Group("/api/v1")
	{
		api.POST("/users", userHandler.RegisterUser)
		api.POST("/sessions", userHandler.Login)
		api.POST("/email_checkers", userHandler.CheckEmailAvailability)
		//api.POST("/avatars", userHandler.UploadAvatar) // none middleware
		api.POST("/avatars", authMiddleware(authService, userService), userHandler.UploadAvatar)

		api.GET("/campaigns", campaignHandler.GetCampaigns)
	}

	router.Run()
	// userInput := user.RegisterUserInput{}
	// userInput.Name = "Dart"
	// userInput.Occupation = "Mobile developer"
	// userInput.Email = "flutter@gmail.com"
	// userInput.Password = "password"
	// userService.RegisterUser(userInput)

	// user := user.User{
	// 	Name:           "springboot",
	// 	Occupation:     "backend",
	// 	Email:          "javaspring@gmail.com",
	// 	PasswordHash:   "java123",
	// 	AvatarFileName: "logo.jpg",
	// 	Role:           "user",
	// 	CreatedAt:      time.Time{},
	// 	UpdatedAt:      time.Time{},
	// }
	// userRepository.Save(user)
}

// Middleware
func authMiddleware(authService auth.Service, userService user.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Ambil header
		authHeader := c.GetHeader("Authorization")
		if !strings.Contains(authHeader, "Bearer") {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		// Bearer <tokentokentoken>
		var tokenString string
		arrayToken := strings.Split(authHeader, " ")
		if len(arrayToken) == 2 {
			tokenString = arrayToken[1]
		}

		// Ambil token
		token, err := authService.ValidateToken(tokenString)
		if err != nil {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		// Validation token
		claim, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		// Ambil user id
		userID := int(claim["user_id"].(float64))
		user, err := userService.GetUserByID(userID)
		if err != nil {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}
		// map "key" : value
		c.Set("currentUser", user) // kirim context(key) ke user handler

	}

	// Middleware steps
	// ambil nilai header Authorization: <Bearer tokentokentoken>
	// dari header Authorization, kita ambil nilai tokennya saja
	// kita validasi token
	// kita ambil user_id
	// ambil user dari db berdasarkan user_id lewat service
	// kita set context isinya user
}
