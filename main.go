package main

import (
	"log"
	"project-campaign/handler"
	"project-campaign/user"

	"github.com/gin-gonic/gin"
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
	userService := user.NewService(userRepository)

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

	userHandler := handler.NewUserHandler(userService)

	router := gin.Default()
	api := router.Group("/api/v1")
	{
		api.POST("/users", userHandler.RegisterUser)
		api.POST("/sessions", userHandler.Login)
		api.POST("/email_checkers", userHandler.CheckEmailAvailability)
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
