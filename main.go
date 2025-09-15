package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"

	"github.com/kiemduong-dev/task-manager/config"
	"github.com/kiemduong-dev/task-manager/models"
	"github.com/kiemduong-dev/task-manager/routes"
)

func main() {
	// load .env
	_ = godotenv.Load()

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL not set")
	}
	config.Connect(dsn)

	// auto migrate
	config.DB.AutoMigrate(&models.User{}, &models.Category{}, &models.Task{})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r := routes.SetupRouter()
	fmt.Println("Server started at :" + port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("server run error:", err)
	}
}
