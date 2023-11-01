package main

import (
	"github.com/gin-gonic/gin"
	"github.com/variant-abdulaziz/controllers"
	"github.com/variant-abdulaziz/initializers"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
	initializers.SyncDB()
}

func main() {
	r := gin.Default()
	r.POST("/api/signup", controllers.SignUp)
	r.POST("/api/login", controllers.Login)
	r.Run()
}
