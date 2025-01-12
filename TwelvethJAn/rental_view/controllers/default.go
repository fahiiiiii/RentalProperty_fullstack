package controllers

import (
    "encoding/json"
    "log"

    beego "github.com/beego/beego/v2/server/web"
    "rental_view/models"
)

type MainController struct {
    beego.Controller
}

func (c *MainController) Get() {
    // Create model
    model := models.NewRentalPropertyModel()
    defer model.Close()

    // Get city from query parameter
    city := c.GetString("city", "Dubai")

    // Fetch properties
    properties, totalCount, err := model.GetProperties(city, 1, 10)
    if err != nil {
        log.Println("Error fetching properties:", err)
        c.Data["PropertiesJSON"] = "[]"
        c.Data["Listings"] = []models.RentalProperty{}
    } else {
        // Convert properties to JSON for JavaScript
        propertiesJSON, _ := json.Marshal(properties)
        c.Data["PropertiesJSON"] = string(propertiesJSON)
        c.Data["Listings"] = properties
        c.Data["TotalCount"] = totalCount
    }

    // Render the template
    c.TplName = "property/index.tpl"
}


// package controllers

// import (
//     "encoding/json"
//     "log"

//     beego "github.com/beego/beego/v2/server/web"
//     "rental_view/models"
// )

// type MainController struct {
//     beego.Controller
// }

// func (c *MainController) Get() {
//     // Create model
//     model := models.NewRentalPropertyModel()
//     defer model.Close()

//     // Get city from query parameter
//     city := c.GetString("city", "Dubai")

//     // Fetch properties
//     properties, totalCount, err := model.GetProperties(city, 1, 10)
//     if err != nil {
//         log.Println("Error fetching properties:", err)
//         c.Data["PropertiesJSON"] = "[]"
//         c.Data["Listings"] = []models.RentalProperty{}
//     } else {
//         // Convert properties to JSON for JavaScript
//         propertiesJSON, _ := json.Marshal(properties)
//         c.Data["PropertiesJSON"] = string(propertiesJSON)
//         c.Data["Listings"] = properties
//         c.Data["TotalCount"] = totalCount
//     }

//     // Render the template
//     c.TplName = "property/index.tpl"
// }