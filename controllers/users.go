package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ibamibrhm/donation-server/helpers"
	"github.com/ibamibrhm/donation-server/models"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

//UserController ...
type UserController struct{}

func verifyPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

type registerInput struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Phone    string `json:"phone"`
}

type loginInput struct {
	EmailUsername string `json:"emailUsername"`
	Password      string `json:"password"`
}

type updateUserInput struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
	Phone    string `json:"phone"`
}

// Register -> create user for route POST /users/regsiter
func (ctrl UserController) Register(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	// Validate input
	var input registerInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create user
	user := models.User{Name: input.Name, Email: input.Email, Username: input.Username, Password: input.Password, Phone: input.Phone}
	if err := db.Create(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": "User created"})
}

// Login -> login user for route POST /users/login
func (ctrl UserController) Login(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	// Validate input
	var input loginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := db.Where("email = ?", input.EmailUsername).Or("username = ?", input.EmailUsername).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Akun belum terdaftar"})
		return
	}

	match := verifyPassword(input.Password, user.Password)

	if !match {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Wrong Password!"})
		return
	}

	token, _ := helpers.CreateToken(user)
	refreshToken, _ := helpers.CreateRefreshToken(user)

	// SetCookie(name string, value string, maxAge int, path string, domain string, secure bool, httpOnly bool)
	// SetCookie adds a Set-Cookie header to the ResponseWriter's headers.
	c.SetCookie("jid", refreshToken, 60*60*24*7 /* 7 days */, "/", "localhost", false, true)

	c.JSON(http.StatusOK, gin.H{"data": gin.H{
		"user":  user,
		"token": token,
	}})
}

// FindUsers -> get all users for route GET /users
func (ctrl UserController) FindUsers(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var users []models.User
	db.Find(&users)

	c.JSON(http.StatusOK, gin.H{"data": users})
}

// FindUser -> get single user for route GET /users/:id
func (ctrl UserController) FindUser(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	authUserID := fmt.Sprint(c.MustGet("userId"))

	if authUserID != c.Param("id") {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized!"})
		return
	}

	// Get model if exist
	var user models.User
	if err := db.Where("id = ?", c.Param("id")).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": user})
}

// UpdateUser -> update single user for route -> PATCH /users/:id
func (ctrl UserController) UpdateUser(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	authUserID := fmt.Sprint(c.MustGet("userId"))

	if authUserID != c.Param("id") {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized!"})
		return
	}

	var user models.User

	if err := db.Where("id = ?", c.Param("id")).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	// Validate input
	var input updateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.Model(&user).Updates(input).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": user})
}

// DeleteUser -> delete single user for route DELETE /users/:id
func (ctrl UserController) DeleteUser(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	authUserID := fmt.Sprint(c.MustGet("userId"))

	if authUserID != c.Param("id") {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized!"})
		return
	}

	// Get model if exist
	var user models.User
	if err := db.Where("id = ?", c.Param("id")).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	db.Delete(&user)

	c.JSON(http.StatusOK, gin.H{"data": true})
}

// Logout -> logout and destroy cookies
func (ctrl UserController) Logout(c *gin.Context) {
	c.SetCookie("jid", "", 1, "/", "localhost", false, true)
	c.JSON(http.StatusOK, gin.H{"data": ""})
}
