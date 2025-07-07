package main

import (
	"fmt"
	"os"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
 "github.com/syrillings/nora-backend/Controllers"
	"github.com/syrillings/nora-backend/Models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"github.com/syrillings/nora-backend/MiddleWare"
	"github.com/syrillings/nora-backend/Services"
)

//Global DB Varaiable that'll connect to my database
var db *gorm.DB


func main() {

  
 //  authMiddleware := middleware.AuthMiddleware(jwtSecret)

    //This block confirms if the .env file exists
    if err := godotenv.Load(); err != nil{
      fmt.Println("No .env file was found")
    }

    //Getting values from the env file
    dsn := os.Getenv("DIRECT_URL")
    if dsn == " "{
      fmt.Println("Values not present")
    }
    //Connecting to my database
    var err error
    db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil{
      fmt.Println("Failed to Connect to Database", err)
    } else {
      fmt.Println("Connection to database succesful")
    }

  //Auto-migrate the models. Automigrate makes sure the data is saved properly to the supabase db
  if err := db.AutoMigrate(&Models.User{}, &Models.Sites{}, &Models.SiteCheck{} ); err != nil {
    fmt.Println("Failed to migrate database:", err)
}

   //Initializing Services with the db. It allows the signup and login functions to use the database
   Services.InitDB(db)

  server := gin.Default()
 setupRoutes(server)
  server.GET("/ping", func(c *gin.Context) {
    c.JSON(200, gin.H{
      "message": "pong",
    })
  })
      server.Run(":5000") 
}

func setupRoutes(server *gin.Engine){
  siteController := controllers.NewSiteController(db)
  server.POST("/signup", Services.Signup)
   server.POST("/login", Services.Login)
 
  //Routes that require auth
  jwtSecret := []byte(os.Getenv("JWT_SECRET"))

  protected := server.Group("/api").Use(middleware.AuthMiddleware(jwtSecret))
  protected.GET("/sites", siteController.GetSites)
  protected.POST("/sites", siteController.AddSite)
  protected.DELETE("/sites/:id", siteController.DeleteSite)
    
}
  
