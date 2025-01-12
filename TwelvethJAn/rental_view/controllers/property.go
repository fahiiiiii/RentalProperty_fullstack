package controllers

import (
    "encoding/json"
    "log"
    
    beego "github.com/beego/beego/v2/server/web"
    "rental_view/models"
)

type PropertyController struct {
    beego.Controller
}

// List handles property listing
func (c *PropertyController) List() {
    // Default parameters
    city := c.GetString("city", "Dubai")
    page, _ := c.GetInt("page", 1)
    pageSize, _ := c.GetInt("pageSize", 10)
    propertyType := c.GetString("type", "")

    // Create model
    model := models.NewRentalPropertyModel()
    defer model.Close()

    var properties []models.RentalProperty
    var totalCount int
    var err error

    // Flexible property fetching
    if propertyType != "" {
        // Fetch properties by specific type
        properties, totalCount, err = model.GetPropertiesByType(city, propertyType, page, pageSize)
    } else {
        // Default listing
        properties, totalCount, err = model.GetProperties(city, page, pageSize)
    }

    // Error handling
    if err != nil {
        log.Printf("Error fetching properties: %v", err)
        c.Data["json"] = map[string]interface{}{
            "success": false,
            "error":   "Failed to retrieve properties",
            "details": err.Error(),
        }
        c.ServeJSON()
        return
    }

    // Prepare response with comprehensive details
    response := map[string]interface{}{
        "success": true,
        "data": map[string]interface{}{
            "properties": properties,
            "pagination": map[string]interface{}{
                "page":       page,
                "pageSize":   pageSize,
                "totalCount": totalCount,
                "totalPages": (totalCount + pageSize - 1) / pageSize,
            },
            "filters": map[string]interface{}{
                "city": city,
                "type": propertyType,
            },
        },
    }

    // Send JSON response
    c.Data["json"] = response
    c.ServeJSON()
}

// Filter handles advanced property filtering
func (c *PropertyController) Filter() {
    // Comprehensive filter parameters
    var filterParams struct {
        City       string   `json:"city"`
        Type       string   `json:"type"`
        MinPrice   float64  `json:"minPrice"`
        MaxPrice   float64  `json:"maxPrice"`
        Bedrooms   int      `json:"bedrooms"`
        Amenities  []string `json:"amenities"`
        Page       int      `json:"page"`
        PageSize   int      `json:"pageSize"`
    }

    // Default values
    filterParams.City = "Dubai"
    filterParams.Page = 1
    filterParams.PageSize = 10

    // Try to parse JSON body
    if err := json.Unmarshal(c.Ctx.Input.RequestBody, &filterParams); err != nil {
        c.Data["json"] = map[string]interface{}{
            "success": false,
            "error":   "Invalid request body",
            "details": err.Error(),
        }
        c.ServeJSON()
        return
    }

    // Create model
    model := models.NewRentalPropertyModel()
    defer model.Close()

    var properties []models.RentalProperty
    var err error
    var totalCount int

    // Flexible filtering logic
    if filterParams.Type != "" {
        // Filter by property type
        properties, totalCount, err = model.GetPropertiesByType(
            filterParams.City, 
            filterParams.Type, 
            filterParams.Page, 
            filterParams.PageSize,
        )
    } else {
        // Advanced filtering
        properties, err = model.FilterProperties(
            filterParams.City, 
            filterParams.MinPrice, 
            filterParams.MaxPrice, 
            filterParams.Bedrooms, 
            filterParams.Amenities,
        )
        totalCount = len(properties)
    }

    // Error handling
    if err != nil {
        log.Printf("Error filtering properties: %v", err)
        c.Data["json"] = map[string]interface{}{
            "success": false,
            "error":   "Failed to filter properties",
            "details": err.Error(),
        }
        c.ServeJSON()
        return
    }

    // Prepare comprehensive response
    response := map[string]interface{}{
        "success": true,
        "data": map[string]interface{}{
            "properties": properties,
            "pagination": map[string]interface{}{
                "page":       filterParams.Page,
                "pageSize":   filterParams.PageSize,
                "totalCount": totalCount,
                "totalPages": (totalCount + filterParams.PageSize - 1) / filterParams.PageSize,
            },
            "appliedFilters": map[string]interface{}{
                "city":       filterParams.City,
                "type":       filterParams.Type,
                "minPrice":   filterParams.MinPrice,
                "maxPrice":   filterParams.MaxPrice,
                "bedrooms":   filterParams.Bedrooms,
                "amenities":  filterParams.Amenities,
            },
        },
    }

    // Send JSON response
    c.Data["json"] = response
    c.ServeJSON()
}

// GetPropertyTypes retrieves available property types
func (c *PropertyController) GetPropertyTypes() {
    // Create model
    model := models.NewRentalPropertyModel()
    defer model.Close()

    // Fetch unique property types
    types, err := model.GetUniquePropertyTypes()
    if err != nil {
        c.Data["json"] = map[string]interface{}{
            "success": false,
            "error":   "Failed to retrieve property types",
            "details": err.Error(),
        }
        c.ServeJSON()
        return
    }

    // Send response
    c.Data["json"] = map[string]interface{}{
        "success": true,
        "data":    types,
    }
    c.ServeJSON()
}
// package controllers

// import (
//     "encoding/json"
    
//     beego "github.com/beego/beego/v2/server/web"
//     "rental_view/models"
// )

// type PropertyController struct {
//     beego.Controller
// }

// // List handles property listing
// func (c *PropertyController) List() {
//     // Default parameters
//     city := c.GetString("city", "Dubai")
//     page, _ := c.GetInt("page", 1)
//     pageSize, _ := c.GetInt("pageSize", 10)

//     // Fetch properties
//     properties, totalCount, err := models.GetRentalProperties(city, page, pageSize)
//     if err != nil {
//         c.Data["json"] = map[string]interface{}{
//             "success": false,
//             "error":   err.Error(),
//         }
//         c.ServeJSON()
//         return
//     }

//     // Prepare response
//     response := map[string]interface{}{
//         "success": true,
//         "data": map[string]interface{}{
//             "properties": properties,
//             "pagination": map[string]interface{}{
//                 "page":       page,
//                 "pageSize":   pageSize,
//                 "totalCount": totalCount,
//                 "totalPages": (totalCount + pageSize - 1) / pageSize,
//             },
//         },
//     }

//     // Send JSON response
//     c.Data["json"] = response
//     c.ServeJSON()
// }

// // Filter handles property filtering
// func (c *PropertyController) Filter() {
//     // Parse filter parameters
//     var filterParams struct {
//         City      string   `json:"city"`
//         MinPrice  float64  `json:"minPrice"`
//         MaxPrice  float64  `json:"maxPrice"`
//         Bedrooms  int      `json:"bedrooms"`
//         Amenities []string `json:"amenities"`
//     }

//     // Try to parse JSON body
//     if err := json.Unmarshal(c.Ctx.Input.RequestBody, &filterParams); err != nil {
//         c.Data["json"] = map[string]interface{}{
//             "success": false,
//             "error":   "Invalid request body",
//         }
//         c.ServeJSON()
//         return
//     }

//     // Default city if not provided
//     if filterParams.City == "" {
//         filterParams.City = "Dubai"
//     }

//     // Filter properties
//     properties, err := models.FilterRentalProperties(
//         filterParams.City, 
//         filterParams.MinPrice, 
//         filterParams.MaxPrice, 
//         filterParams.Bedrooms, 
//         filterParams.Amenities,
//     )
//     if err != nil {
//         c.Data["json"] = map[string]interface{}{
//             "success": false,
//             "error":   err.Error(),
//         }
//         c.ServeJSON()
//         return
//     }

//     // Send filtered properties
//     c.Data["json"] = map[string]interface{}{
//         "success": true,
//         "data":    properties,
//     }
//     c.ServeJSON()
// }