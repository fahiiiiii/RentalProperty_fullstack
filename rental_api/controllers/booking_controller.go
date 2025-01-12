// controllers/booking_controller.go

package controllers

import (
	"rental_api/models"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"sync"
	"rental_api/conf"
	// "rental_api/conf"
	beego "github.com/beego/beego/v2/server/web"
)

// Operations about Users
type BookingControllerController struct {
	beego.Controller
	rateLimiter sync.Mutex
}


// Create handles POST requests to add a new property
func (c *BookingController) Create() {
	var property models.Property
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &property); err != nil {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Ctx.Output.Body([]byte("Invalid input"))
		return
	}

	if err := conf.DB.Create(&property).Error; err != nil {
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Ctx.Output.Body([]byte("Failed to create property"))
		return
	}

	c.Data["json"] = property
	c.ServeJSON()
}

// Get handles GET requests to fetch a property by ID
func (c *BookingController) Get() {
	id, err := strconv.Atoi(c.Ctx.Input.Param(":id"))
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Ctx.Output.Body([]byte("Invalid property ID"))
		return
	}

	var property models.Property
	if err := conf.DB.First(&property, id).Error; err != nil {
		c.Ctx.Output.SetStatus(http.StatusNotFound)
		c.Ctx.Output.Body([]byte("Property not found"))
		return
	}

	c.Data["json"] = property
	c.ServeJSON()
}

// Update handles PUT requests to update a property by ID
func (c *BookingController) Update() {
	id, err := strconv.Atoi(c.Ctx.Input.Param(":id"))
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Ctx.Output.Body([]byte("Invalid property ID"))
		return
	}

	var property models.Property
	if err := conf.DB.First(&property, id).Error; err != nil {
		c.Ctx.Output.SetStatus(http.StatusNotFound)
		c.Ctx.Output.Body([]byte("Property not found"))
		return
	}

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &property); err != nil {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Ctx.Output.Body([]byte("Invalid input"))
		return
	}

	if err := conf.DB.Save(&property).Error; err != nil {
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Ctx.Output.Body([]byte("Failed to update property"))
		return
	}

	c.Data["json"] = property
	c.ServeJSON()
}

// Delete handles DELETE requests to remove a property by ID
func (c *BookingController) Delete() {
	id, err := strconv.Atoi(c.Ctx.Input.Param(":id"))
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Ctx.Output.Body([]byte("Invalid property ID"))
		return
	}

	if err := conf.DB.Delete(&models.Property{}, id).Error; err != nil {
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Ctx.Output.Body([]byte("Failed to delete property"))
		return
	}

	c.Ctx.Output.SetStatus(http.StatusNoContent)
}

// List handles GET requests to fetch all properties
func (c *BookingController) List() {
	var properties []models.Property
	if err := conf.DB.Find(&properties).Error; err != nil {
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Ctx.Output.Body([]byte("Failed to fetch properties"))
		return
	}

	c.Data["json"] = properties
	c.ServeJSON()
}

// ProcessAllCities refactored to use GET operation
func (c *BookingController) ProcessAllCities() {
	queries := []string{"New York", "Los Angeles", "Chicago"}
	var results []map[string]interface{}

	for _, query := range queries {
		cityData := map[string]interface{}{
			"city": query,
			"properties": []string{"Hotel A", "Hotel B"},
		}
		results = append(results, cityData)
	}

	c.Data["json"] = results
	c.ServeJSON()
}

// ProcessPropertyDetails as an Update operation
func (c *BookingController) ProcessPropertyDetails() {
	id, err := strconv.Atoi(c.Ctx.Input.Param(":id"))
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Ctx.Output.Body([]byte("Invalid property ID"))
		return
	}

	propertyDetails := map[string]interface{}{
		"description": "Luxury property with sea view",
		"rating": 4.5,
		"images": []string{"image1.jpg", "image2.jpg"},
	}

	log.Printf("Processing property ID %d with details: %+v", id, propertyDetails)
	c.Data["json"] = propertyDetails
	c.ServeJSON()
}

// Index handles the root endpoint
func (c *BookingController) Index() {
	c.TplName = "index.html"
}
