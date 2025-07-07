package middleware

import (
	"strings"
	"fmt"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

//AuthMiddleware verifies the jwt token
func AuthMiddleware(jwtSecret []byte) gin.HandlerFunc{
    
    return func (c *gin.Context){
//This block get tokens from authorization 
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" { 
          c.JSON(http.StatusUnauthorized, gin.H{"error":"Auth Token not found"})
          c.Abort()
		  return
		}  

		// This block extracts the token from "Bearer"
      parts := strings.SplitN(authHeader, " ", 2)
      if len(parts) != 2 || parts[0] != "Bearer" {
      c.JSON(http.StatusUnauthorized, gin.H{"error":"Invalid Auth Header"})
      c.Abort()
      return
}

//This block extracts the token from "Bearer"	
    tokenString := parts[1]
	//Here, I'm validating that the auth header is
	//  more than seven characters long and come after
	//the word Bearer
    if len(authHeader)>7 && authHeader [:7] == "Bearer"{
		//"Bearer" comes before the token and in the authHeader 
		// and this line below tells the program to jump seven places 
		// which is "Bearer" and the space after it and just get the 
		// token  straight
		tokenString = authHeader[7:]
	}
 
	//This block parses and validates json tokens
	token, err:= jwt.Parse(tokenString, func (token *jwt.Token) (interface{}, error){
    if  token.Method != jwt.SigningMethodHS256{
		return nil, jwt.ErrSignatureInvalid
	}
       return jwtSecret, nil
	})
   
   //This block checks for token validation errors
	if err != nil || !token.Valid{
		fmt.Println("Authorization header:", c.GetHeader("Authorization"))
fmt.Println("JWT Secret:", string(jwtSecret))
fmt.Println("Token string:", tokenString)
		c.JSON(http.StatusUnauthorized, gin.H{"error":"Invalid Token"})
		c.Abort()
		return
	}

    if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid{
		c.Set("userID", uint(claims["user_id"].(float64)))
		c.Set("email", claims["email"].(string))
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error":"Invalid Credentials"})
		c.Abort()
		return
	}
       c.Next()
}

}