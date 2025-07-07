package Models

import(
	"gorm.io/gorm"
	"golang.org/x/crypto/bcrypt"
	
)

//User Model: This describes the attributes a user should have
type User struct{
	gorm.Model
	Username   string `gorm:"unique;not null"`
	Email      string `gorm:"unique; not null"` 
	Sites      []Sites
	PasswordHash string `gorm:"unique"`
}

type SignupRequest struct {
	Username   string `json:"Username" binding:"required"`
	Email      string `json:"email" binding:"required,email"`
	Password   string `json:"password" binding:"required"`
  } //this is what trying to get into the system will send, the things required to get into

type LoginRequest struct{
	Email string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}//this is the data required to login


type AuthResponse struct {
	Token string
    Username string
	Email string
	Sites []Sites
} //this is what will be returned after logging in succesfully

// This function hashes/encrypts the password before saving it
func (u *User) HashPassword(password string) error {
   bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
   if err != nil{
	  return err
   }
   u.PasswordHash = string(bytes)
   return nil
}

// The checkPassword function verifies the passwords
func (u *User) CheckPassword(password string) error {
    return bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
}

