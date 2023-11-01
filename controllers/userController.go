package controllers

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
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

func Login(c *gin.Context) {
	// get data from body
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

	// get user specified from db
	var user = models.User{Username: body.Username}
	initializers.DB.First(&user)

	if user.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "0: incorrect username/password",
		})
		return
	}

	// compare the stored hash with the provided password
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "1: incorrect username/password",
		})
		return
	}

	// generate a jwt token
	key := os.Getenv("JWT_SECRET")
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256, jwt.MapClaims{
			"sub": user.ID,
		})
	tokenString, err := token.SignedString([]byte(key))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "server error",
		})
		return
	}

	// respond
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("token", tokenString, 60*60*24*30, "", "", false, true)
}
