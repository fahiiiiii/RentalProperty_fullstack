package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
	"context"
	"github.com/beego/beego/v2/server/web"
	"golang.org/x/time/rate"
	 "github.com/beego/beego/v2/core/config"
    "strconv"
	// "github.com/joho/godotenv"
	"os"
	"log"
	"rental_api/models"
	"rental_api/conf"  // Add this import
	
)
// BookingController handles all booking.com API related operations
type BookingController struct {
    web.Controller
    rapidAPIKey     string
    uniqueCountries map[string]bool
    uniqueCities    map[string]bool
    countryCities   map[string][]string
    cityProperties  map[string][]string
    rateLimiter    *rate.Limiter
    mutex          sync.RWMutex
}

// BookingController handles all booking.com API related operations
// Prepare is Beego's initialization method that runs before any request
func (c *BookingController) Prepare() {
    // Load configuration
    if err := conf.LoadConfig(); err != nil {
        log.Fatalf("Error loading configuration: %v", err)
    }

    // Initialize controller fields
    c.rapidAPIKey = conf.AppConfig.API.RapidAPIKey
    if c.rapidAPIKey == "" {
        log.Fatalf("RAPIDAPI_KEY is not set in the configuration")
    }

    // Initialize maps
    c.uniqueCountries = make(map[string]bool)
    c.uniqueCities = make(map[string]bool)
    c.countryCities = make(map[string][]string)
    c.cityProperties = make(map[string][]string)

    // Initialize rate limiter using configuration
    requestsPerSecond := rate.Limit(conf.AppConfig.RateLimit.RequestsPerSecond)
    burstSize := conf.AppConfig.RateLimit.BurstSize
    c.rateLimiter = rate.NewLimiter(requestsPerSecond, burstSize)
}

// Prepare is Beego's initialization method that runs before any request
// Prepare is Beego's initialization method that runs before any request
func (c *BookingController) Prepare() {
    // Get configuration
    if err := conf.LoadConfig(); err != nil {
        log.Fatalf("Error loading configuration: %v", err)
    }

    // Initialize controller fields
    c.rapidAPIKey = conf.AppConfig.API.RapidAPIKey
    if c.rapidAPIKey == "" {
        log.Fatalf("RAPIDAPI_KEY is not set in the configuration")
    }

    // Initialize maps
    c.uniqueCountries = make(map[string]bool)
    c.uniqueCities = make(map[string]bool)
    c.countryCities = make(map[string][]string)
    c.cityProperties = make(map[string][]string)

    // Initialize rate limiter using configuration
    requestsPerSecond := rate.Limit(conf.AppConfig.RateLimit.RequestsPerSecond)
    burstSize := conf.AppConfig.RateLimit.BurstSize
    c.rateLimiter = rate.NewLimiter(requestsPerSecond, burstSize)
}

// Get handles GET requests for the booking resource
func (c *BookingController) Get() {
	action := c.GetString("action")
	switch action {
	case "cities":
		c.ListCities()
	case "properties":
		c.ListProperties()
	case "summary":
		c.GetSummary()
	default:
		c.Data["json"] = map[string]string{"error": "Invalid action"}
		c.ServeJSON()
	}
}

// Post handles POST requests for the booking resource
func (c *BookingController) Post() {
	action := c.GetString("action")
	switch action {
	case "process-cities":
		if err := c.ProcessAllCities(); err != nil {
			c.Data["json"] = map[string]string{"error": err.Error()}
		} else {
			c.Data["json"] = map[string]string{"status": "success"}
		}
	case "process-properties":
		if err := c.ProcessAllProperties(); err != nil {
			c.Data["json"] = map[string]string{"error": err.Error()}
		} else {
			c.Data["json"] = map[string]string{"status": "success"}
		}
	case "process-details":
		if err := c.ProcessAllHotelDetails(); err != nil {
			c.Data["json"] = map[string]string{"error": err.Error()}
		} else {
			c.Data["json"] = map[string]string{"status": "success"}
		}
	default:
		c.Data["json"] = map[string]string{"error": "Invalid action"}
	}
	c.ServeJSON()
}

// ListCities returns all available cities
func (c *BookingController) ListCities() {
	var cities []string
	if err := models.GetAllCities(&cities); err != nil {
		c.Data["json"] = map[string]string{"error": err.Error()}
	} else {
		c.Data["json"] = cities
	}
	c.ServeJSON()
}

// ListProperties returns properties for a given location
func (c *BookingController) ListProperties() {
	location := c.GetString("location")
	if location == "" {
		c.CustomAbort(400, "Location is required")
		return
	}

	var properties []models.RentalProperty
	if err := models.GetPropertiesByLocation(location, &properties); err != nil {
		c.CustomAbort(500, fmt.Sprintf("Failed to fetch properties: %v", err))
		return
	}

	c.Data["json"] = properties
	c.ServeJSON()
}

// GetPropertyDetails returns details for a specific property
func (c *BookingController) GetPropertyDetails() {
	propertyId, err := c.GetInt(":propertyId")
	if err != nil {
		c.CustomAbort(400, "Invalid property ID")
		return
	}

	var property models.Property
	if err := models.GetPropertyById(propertyId, &property); err != nil {
		c.CustomAbort(404, "Property not found")
		return
	}

	c.Data["json"] = property
	c.ServeJSON()
}

// GetSummary returns a summary of all data
func (c *BookingController) GetSummary() {
	summary := struct {
		Countries      map[string]bool     `json:"countries"`
		Cities         map[string]bool     `json:"cities"`
		CountryCities  map[string][]string `json:"country_cities"`
		CityProperties map[string][]string `json:"city_properties"`
	}{
		Countries:      c.uniqueCountries,
		Cities:         c.uniqueCities,
		CountryCities:  c.countryCities,
		CityProperties: c.cityProperties,
	}

	c.Data["json"] = summary
	c.ServeJSON()
}

// makeRateLimitedRequest handles rate-limited API requests
func (c *BookingController) makeRateLimitedRequest(req *http.Request) (*http.Response, error) {
    if err := c.rateLimiter.Wait(context.Background()); err != nil {
        return nil, fmt.Errorf("rate limiter error: %v", err)
    }

    timeout := time.Duration(conf.AppConfig.API.RequestTimeout) * time.Second
    client := &http.Client{Timeout: timeout}
    
    // Add RapidAPI headers
    req.Header.Add("X-RapidAPI-Key", c.rapidAPIKey)
    req.Header.Add("X-RapidAPI-Host", "booking-com.p.rapidapi.com")
    
    return client.Do(req)
}

// SaveToDatabase saves all data to the database
func (c *BookingController) SaveToDatabase() {
	if err := models.SaveLocations(c.countryCities); err != nil {
		c.CustomAbort(500, fmt.Sprintf("Failed to save locations: %v", err))
		return
	}

	if err := models.SaveRentalProperties(c.cityProperties); err != nil {
		c.CustomAbort(500, fmt.Sprintf("Failed to save rental properties: %v", err))
		return
	}

	c.Data["json"] = map[string]string{"status": "success"}
	c.ServeJSON()
}

// Additional helper methods for processing data
func (c *BookingController) ProcessAllCities() error {
	queries := c.generateQueries()
	results := make(chan struct{}, len(queries))
	semaphore := make(chan struct{}, 1)

	for _, query := range queries {
		semaphore <- struct{}{}
		go func(q string) {
			defer func() { <-semaphore }()
			c.processCities(q, results)
		}(query)
	}

	for range queries {
		<-results
	}

	return nil
}

func (c *BookingController) generateQueries() []string {
	queries := []string{}
	for char := 'A'; char <= 'Z'; char++ {
		queries = append(queries, string(char))
	}
	
	prefixes := []string{
		"a", "the", "new", "old", "big", "small", 
		"north", "south", "east", "west", "central",
	}
	
	queries = append(queries, prefixes...)
	return queries
}

// Initialize registers the controller with Beego
func init() {
	web.Router("/v1/booking/", &BookingController{})
	web.Router("/v1/booking/:propertyId", &BookingController{}, "get:GetPropertyDetails")
}