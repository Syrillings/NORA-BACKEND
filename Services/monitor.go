package Services;

import (
	"fmt"
	"gorm.io/gorm"
	"context"
	"net/http"
	"time"
	"github.com/syrillings/nora-backend/Models"
)

//This struct contains a reference to the db
type MonitorService struct{
	db *gorm.DB
}

//Creates a new monitor instance for a new user.....I think
func NewMonitorService(db *gorm.DB) *MonitorService{
   return &MonitorService{ db : db }
}

func (ms *MonitorService) StartMonitoring(){
	//This link ticks every ten minutes, when the sites are checked
	ticker := time.NewTicker(1*time.Minute)
	//Loop runs after every tick
	for range ticker.C{
		ms.CheckAllWebsites()
	}
}

//This block Checks the sites to confirm 
func (ms* MonitorService)CheckAllWebsites(){
      var Sites []Models.Sites
	  if err := ms.db.Where("is_active = ?", true).Find(&Sites).Error; err != nil{
		fmt.Println("Error fetching site for monitoring: ", err)
		return
	  }
 
	  for _, Sites := range Sites {
		//GoRoutine !!!!
	    go ms.CheckSite(Sites)
	}
}


func (ms *MonitorService)CheckSite(Sites Models.Sites){
   check := Models.SiteCheck{
	SiteID: int(Sites.ID),
   }

   start := time.Now()
  
   ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
   defer cancel()

   //Sends a get request to url 
   req, err := http.NewRequestWithContext(ctx, "GET", Sites.URL, nil )
   if err != nil{
	ms.RecordCheck(Sites, check, Models.StatusDown, 0, err.Error())
	return
   }

   //Sends the request
   client := &http.Client{
	Timeout: 60*time.Second,
   }

   resp,err := client.Do(req)
   if err != nil{
	ms.RecordCheck(Sites, check, Models.StatusDown, 0, err.Error())
	return
   }
    defer resp.Body.Close()

	Latency := time.Since(start).Milliseconds()
	status := Models.StatusUp
	if resp.StatusCode >= 400 {
		status = Models.StatusDown
	}

	ms.RecordCheck(Sites, check, status, resp.StatusCode, "")
	check.Latency = Latency
}

func (ms *MonitorService)RecordCheck(Sites Models.Sites, check Models.SiteCheck, status Models.SiteStatus, statusCode int, errorMsg string){
	now := time.Now()

	check.SiteStatus = status
	check.StatusCode = statusCode
	check.Error = errorMsg

	//Saves the check to my db....wetin I dey actually write?
	ms.db.Create(&check)

	//Update site status
	ms.db.Model(&Sites).Updates(map[string]interface{}{
        "last_checked": now,
        "last_status":  status,
    })
    
    // If website is down, send notification
    if status == Models.StatusDown {
        go ms.notifyUser(Sites)
    }

   if status ==Models.StatusUp{
       go ms.applaudUser(Sites)
   }
}


//Todo; Figure out how I'm gonna send the notifications
func (ms *MonitorService)notifyUser(Sites Models.Sites){
   fmt.Println(Sites.Name+" @ "+Sites.URL +" is down ")
}

func (ms *MonitorService)applaudUser(Sites Models.Sites){
	fmt.Println(Sites.Name+" @ "+Sites.URL +" is active ")
 }