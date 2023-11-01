package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/variant-abdulaziz/initializers"
	"github.com/variant-abdulaziz/models"
	"golang.org/x/crypto/bcrypt"
)

func SignUp(c *gin.Context) {
	// get username & password
	var body struct {
		Username string
		Password string
	}

	if c.BindJSON(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to read body",
		})
		return
	}
	// hash the password
	fmt.Println(body.Username, body.Password)
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to hash password",
		})
		return
	}

	// create user
	user := models.User{Username: body.Username, Password: string(hash)}

	// persist user in db
	result := initializers.DB.Create(&user)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "server error",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{})
}
