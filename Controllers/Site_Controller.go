package controllers

import (
	"fmt"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/syrillings/nora-backend/Models"
	"gorm.io/gorm"
)


type SiteController struct{
	db*gorm.DB
}

func NewSiteController(db*gorm.DB) *SiteController{
		return &SiteController{db:db}
}

//Adds a new site to the database
// sc*SiteController means the function connects to the SiteController db
func (sc*SiteController) AddSite (c*gin.Context){
	var input struct{
		Name string `json:"name"`
		URL string `json:"url"`
	}
	if err := c.ShouldBindJSON(&input); err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error":err.Error()})
		return
	}

	userID := c.GetUint("userID")
	if userID == 0{
	   c.JSON(http.StatusUnauthorized, gin.H{"err":"User not authenticated"})
	   return
	} 

	//This block creates a new site with the user's ID
	//A model is a blueprint on the data an object is supposed to cary; much like a schema
	newSite:= Models.Sites{
		URL : input.URL,
		Name : input.Name,
		UserID: uint(userID),
	}

    if err := sc.db.Create(&newSite). Error; err != nil{
      c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create new site"})
	  return
	}
     
      c.JSON(http.StatusCreated, newSite)

}

       
 func (sc*SiteController) GetSites (c*gin.Context){
	
	userID := c.GetUint("userID")
	if userID == 0{
	   c.JSON(http.StatusUnauthorized, gin.H{"err":"User not authenticated"})
	} 

	var Sites []Models.Sites
	if err := sc.db.Where("user_id = ?", userID).Find(&Sites).Error; err != nil{
		//Always remember it's gin.H{ curly braces }
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user's sites"});
		c.Abort()
		return
	}

      c.JSON(http.StatusOK, Sites)

  }

  func (sc *SiteController) DeleteSite(c*gin.Context){
       
      siteID := c.Param("id")
	  userID := c.GetUint("userID")

	  result := sc.db.Where("id = ? AND user_id = ?", siteID, userID).Delete(&Models.Sites{})
      if result.Error != nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error":"Failed to delete site"})
	    c.Abort()
		return
	}

	if result.RowsAffected == 0{
		c.JSON(http.StatusNotFound, gin.H{"error":"Site not found"})
	}

      c.String(http.StatusNoContent, "No Content")
      fmt.Println("No Content")
}

