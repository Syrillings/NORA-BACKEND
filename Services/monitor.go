package Services;

import (
	"fmt"
	"gorm.io/gorm"
	"io"
	"context"
	"net/http"
	"strings"
	"crypto/tls"
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


func (ms *MonitorService) CheckSite(site Models.Sites) {
    check := Models.SiteCheck{
        SiteID: int(site.ID),
    }

    start := time.Now()

    // Ensure URL has a scheme
    urlToCheck := site.URL
    if !strings.HasPrefix(urlToCheck, "http://") && !strings.HasPrefix(urlToCheck, "https://") {
        urlToCheck = "https://" + urlToCheck
    }

    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    // Create a new request
    req, err := http.NewRequestWithContext(ctx, "GET", urlToCheck, nil)
    if err != nil {
        ms.RecordCheck(site, check, Models.StatusDown, 0, fmt.Sprintf("Request creation failed: %v", err))
        return
    }

    // Set headers to mimic a browser request
    req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
    req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
    req.Header.Set("Accept-Language", "en-US,en;q=0.5")

    // Configure the HTTP client
    client := &http.Client{
        Timeout: 30 * time.Second,
        // Follow up to 10 redirects
        CheckRedirect: func(req *http.Request, via []*http.Request) error {
            if len(via) >= 10 {
                return fmt.Errorf("stopped after 10 redirects")
            }
            return nil
        },
        // Skip SSL verification (use only for development)
        Transport: &http.Transport{
            TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
        },
    }

    // Send the request
    resp, err := client.Do(req)
    if err != nil {
        ms.RecordCheck(site, check, Models.StatusDown, 0, fmt.Sprintf("Request failed: %v", err))
        return
    }
    defer resp.Body.Close()

    // Read the response body to ensure the full request completes
    _, err = io.Copy(io.Discard, resp.Body)
    if err != nil {
        ms.RecordCheck(site, check, Models.StatusDown, 0, fmt.Sprintf("Failed to read response: %v", err))
        return
    }

    latency := time.Since(start).Milliseconds()
    status := Models.StatusUp
    if resp.StatusCode >= 400 {
        status = Models.StatusDown
    }

    ms.RecordCheck(site, check, status, resp.StatusCode, "")
    check.Latency = latency
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