package Services

import (
	"net/http"
	"os"
	"time"
   	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/syrillings/nora-backend/Models"
	"gorm.io/gorm"
)

var (
	db        *gorm.DB //This line creates the database instance
	jwtSecret = []byte(os.Getenv("JWT_SECRET"))
)

// InitDB initializes the database connection for the Services package
func InitDB(database *gorm.DB) {
	db = database
	db.AutoMigrate(&Models.User{})
	db.AutoMigrate(&Models.Sites{})
}

func generateToken(userID uint, email string) (string, error) {
	// Set token claims
	claims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
		"iat":     time.Now().Unix(),
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token
	return token.SignedString(jwtSecret)
}

// Signing up Functionanlity
func Signup(c *gin.Context) {
	//Extracts and converts the requested data from my api into usable variables
	var req Models.SignupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	//Be sure whether this user's email has been used before
	var user Models.User
	if db.Where("email = ?", req.Email).First(&user).Error == nil {
	c.JSON(401, gin.H{"error": "Oops! Email is already in use"})
		return
	}

	//Creates a new user
	newUser := Models.User{
		Username: req.Username,
		Email:    req.Email,
	}

	//Hashes Password
	if err := newUser.HashPassword(req.Password); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating user"})
		return
	}

	//This block will save the user to my database
	if err := db.Create(&newUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error saving user"})
	}

	//This block generates the token the user will use to sign into Nora
	token, err := generateToken(newUser.ID, user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating token"})
		return
	}

	//Returns the response that logging in would generate
	c.JSON(http.StatusCreated, Models.AuthResponse{
		Token:    token,
		Username: newUser.Username,
		Email:    newUser.Email,
	})

}

// Login Functionality
func Login(c *gin.Context) {
	var req Models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(401, gin.H{"error": "Invalid Credentials"})
		return
	}

	//This block finds user by email by skimming through Nora's database
	var user Models.User
	if err := db.Where("email=?", req.Email).First(&user).Error; err != nil {
		//This error is returned when a user isn't found
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Login Credetials"})
		return
	}

	//This block verifies if the password is right or wrong
	if err := user.CheckPassword(req.Password); err != nil {
		//This error is returned when a password is wrong
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Password or email is incorrect"})
		return
	}

	token, err := generateToken(user.ID, user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating token"})
		return
	}

	c.JSON(http.StatusOK, Models.AuthResponse{
		Username: user.Username,
		Token:    token,
		Email:    user.Email,
		Sites:    user.Sites,
	})

}
