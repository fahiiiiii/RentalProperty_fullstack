// // controllers/booking_controller.go
// package controllers

// import (
// 	"encoding/json"
// 	"fmt"
// 	"io"
// 	"log"
// 	"net/http"
// 	"net/url"
// 	"strings"
// 	"sync"
// 	"time"
// 	"math"
//     "rental_api/conf"
// 	"rental_api/models" // You'll need to adjust this import path based on your module name
//     "golang.org/x/time/rate"
// 	"context"
// 	"github.com/joho/godotenv"
// 	"os"


//     "github.com/beego/beego/v2/server/web"
//     // "github.com/astaxie/beego"

//     // "gorm.io/gorm"
//     // "regexp"
//     // "strconv"
   
    
// )

// // BookingController handles all booking.com API related operations
// // type BookingController struct {
// // 	uniqueCountries map[string]bool
// // 	uniqueCities    map[string]bool
// // 	countryCities   map[string][]string
// // 	cityProperties  map[string][]string
// // 	mutex           sync.Mutex
// // }
// type BookingController struct {
//     uniqueCountries map[string]bool
//     uniqueCities    map[string]bool
//     countryCities   map[string][]string
//     cityProperties  map[string][]string
//     mutex           sync.Mutex
//     rateLimiter    *rate.Limiter
// 	rapidAPIKey     string //!added rapidApiKey
// 	web.Controller
    
// }
// // NewBookingController creates a new instance of BookingController
// func NewBookingController() *BookingController {
// 	// Load environment variables from .env file
// 	err := godotenv.Load()
// 	if err != nil {
// 		log.Fatalf("Error loading .env file")
// 	}

// 	// Get the RapidAPI key from the environment variable
// 	rapidAPIKey := os.Getenv("RAPIDAPI_KEY")
// 	if rapidAPIKey == "" {
// 		log.Fatalf("RAPIDAPI_KEY is not set in the environment")
// 	}
	
// 	//! Created a rate limiter with 5 requests per minute
//     // The first parameter (5) is the rate limit
//     // The second parameter (1) is the burst size
//     return &BookingController{
//         uniqueCountries: make(map[string]bool),
//         uniqueCities:    make(map[string]bool),
//         countryCities:   make(map[string][]string),
//         cityProperties:  make(map[string][]string),
//         rateLimiter:    rate.NewLimiter(rate.Every(12*time.Second), 1), // 5 requests per minute = 1 request per 12 seconds
// 		rapidAPIKey:     rapidAPIKey, //! Initialize the rapidAPIKey field
//     }
// }
// //! helper method for rate-limited requests
// func (c *BookingController) makeRateLimitedRequest(req *http.Request) (*http.Response, error) {
//     // Wait for rate limiter
//     err := c.rateLimiter.Wait(context.Background())
//     if err != nil {
//         return nil, fmt.Errorf("rate limiter error: %v", err)
//     }

//     client := &http.Client{
//         Timeout: 10 * time.Second,
//     }
//     return client.Do(req)
// }
// // GetSummary returns the current state of all data
// func (c *BookingController) GetSummary() interface{} {
// 	return struct {
// 		Countries      map[string]bool            `json:"countries"`
// 		Cities         map[string]bool            `json:"cities"`
// 		CountryCities  map[string][]string        `json:"country_cities"`
// 		CityProperties map[string][]string        `json:"city_properties"`
// 	}{
// 		Countries:      c.uniqueCountries,
// 		Cities:         c.uniqueCities,
// 		CountryCities:  c.countryCities,
// 		CityProperties: c.cityProperties,
// 	}
// }



// //! Modified ProcessAllCities to limit concurrent requests
// func (c *BookingController) ProcessAllCities() error {
//     queries := c.generateQueries()
//     results := make(chan struct{}, len(queries))
    
//     // Reduce concurrent requests to 1 to better control rate limiting
//     semaphore := make(chan struct{}, 1)

//     for _, query := range queries {
//         semaphore <- struct{}{}
//         go func(q string) {
//             defer func() { <-semaphore }()
//             c.processCities(q, results)
//         }(query)
//     }

//     // Wait for all queries to complete
//     for range queries {
//         <-results
//     }

//     return nil
// }

// // ProcessAllProperties processes properties for all cities
// func (c *BookingController) ProcessAllProperties() error {
// 	return c.processPropertiesForCities()
// }

// // Private methods

// func (c *BookingController) generateQueries() []string {
// 	queries := []string{}
	
// 	// Alphabet queries
// 	for char := 'A'; char <= 'Z'; char++ {
// 		queries = append(queries, string(char))
// 	}

// 	// Common prefixes and patterns
// 	prefixes := []string{
// 		"a", "the", "new", "old", "big", "small", 
// 		"north", "south", "east", "west", "central",
// 	}

// 	for _, prefix := range prefixes {
// 		queries = append(queries, prefix)
// 	}

// 	return queries
// }

// func (c *BookingController) processCities(query string, results chan<- struct{}) {
// 	defer func() { results <- struct{}{} }()

// 	cities, err := c.fetchCities(query)
// 	if err != nil {
// 		log.Printf("Error fetching cities for query '%s': %v", query, err)
// 		return
// 	}

// 	c.mutex.Lock()
// 	defer c.mutex.Unlock()

// 	for _, city := range cities {
// 		country := strings.TrimSpace(strings.ToUpper(city.Country))
// 		cityName := strings.TrimSpace(strings.ToUpper(city.CityName))

// 		if country != "" {
// 			c.uniqueCountries[country] = true
// 		}

// 		if cityName != "" {
// 			c.uniqueCities[cityName] = true
// 		}

// 		if country != "" && cityName != "" {
// 			if _, exists := c.countryCities[country]; !exists {
// 				c.countryCities[country] = []string{}
// 			}
			
// 			cityExists := false
// 			for _, existingCity := range c.countryCities[country] {
// 				if existingCity == cityName {
// 					cityExists = true
// 					break
// 				}
// 			}
			
// 			if !cityExists {
// 				c.countryCities[country] = append(c.countryCities[country], cityName)
// 			}
// 		}
// 	}
// }

// func (c *BookingController) fetchCities(query string) ([]models.City, error) {
	
// 	//!using rate limit
// 	apiURL := fmt.Sprintf("https://booking-com18.p.rapidapi.com/stays/auto-complete?query=%s", query)
    
//     req, err := http.NewRequest("GET", apiURL, nil)
//     if err != nil {
//         return nil, fmt.Errorf("error creating request: %v", err)
//     }

//     req.Header.Add("x-rapidapi-host", "booking-com18.p.rapidapi.com")
// 	req.Header.Add("x-rapidapi-key", c.rapidAPIKey) // Use the stored RapidAPI key
//     // req.Header.Add("x-rapidapi-key", "79d933f58amsh0baa13f673b03f0p16d4a2jsnb299a967d295")

//     resp, err := c.makeRateLimitedRequest(req)
//     if err != nil {
//         return nil, fmt.Errorf("error sending request: %v", err)
//     }
//     defer resp.Body.Close()

// 	body, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		return nil, fmt.Errorf("error reading response body: %v", err)
// 	}

// 	if resp.StatusCode != http.StatusOK {
// 		return nil, fmt.Errorf("API request failed with status code: %d, body: %s", 
// 			resp.StatusCode, string(body))
// 	}

// 	var citiesResp struct {
// 		Data []models.City `json:"data"`
// 	}
// 	err = json.Unmarshal(body, &citiesResp)
// 	if err != nil {
// 		return nil, fmt.Errorf("error parsing JSON: %v", err)
// 	}

// 	return citiesResp.Data, nil
// }

// func (c *BookingController) processPropertiesForCities() error {
// 	//!using limit for req
// 	propertyResults := make(chan struct {
//         City       models.CityKey
//         Properties []models.Property
//         Err        error
//     }, len(c.uniqueCities))

//     // Reduce concurrent requests to 1 to better control rate limiting
//     semaphore := make(chan struct{}, 1)
//     var wg sync.WaitGroup

//     for country, cities := range c.countryCities {
//         for _, cityName := range cities {
//             wg.Add(1)
            
//             go func(city, country string) {
//                 defer wg.Done()
                
//                 semaphore <- struct{}{} 
//                 defer func() { <-semaphore }()

//                 properties, err := c.fetchPropertiesWithRetry(city, country, 3)
                
//                 propertyResults <- struct {
//                     City       models.CityKey
//                     Properties []models.Property
//                     Err        error
//                 }{
//                     City:       models.CityKey{Name: city, Country: country},
//                     Properties: properties,
//                     Err:        err,
//                 }
//             }(cityName, country)
//         }
//     }

// 	go func() {
// 		wg.Wait()
// 		close(propertyResults)
// 	}()

// 	for result := range propertyResults {
// 		if result.Err != nil {
// 			log.Printf("Error fetching properties for %s, %s: %v", 
// 				result.City.Name, result.City.Country, result.Err)
// 			continue
// 		}

// 		c.processPropertyResult(result)
// 	}

// 	return nil
// }

// func (c *BookingController) processPropertyResult(result struct {
// 	City       models.CityKey
// 	Properties []models.Property
// 	Err        error
// }) {
// 	if len(result.Properties) == 0 {
// 		return
// 	}

// 	c.mutex.Lock()
// 	defer c.mutex.Unlock()

// 	c.cityProperties[result.City.Name] = []string{}
	
// 	maxProperties := 20
// 	if len(result.Properties) < maxProperties {
// 		maxProperties = len(result.Properties)
// 	}

	
// 	for _, prop := range result.Properties[:maxProperties] {
// 		// Store only the property name without additional details
// 		c.cityProperties[result.City.Name] = append(
// 			c.cityProperties[result.City.Name], 
// 			prop.Name, // Only store the property name
// 		)
// 	}
// }

// func (c *BookingController) fetchPropertiesWithRetry(cityName, country string, maxRetries int) ([]models.Property, error) {
//     for attempt := 0; attempt < maxRetries; attempt++ {
//         properties, err := c.fetchPropertiesForCity(cityName, country)
        
//         if err == nil {
//             return properties, nil
//         }

//         if strings.Contains(err.Error(), "Too many requests") || 
//            strings.Contains(err.Error(), "You are not subscribed") {
//             waitTime := time.Duration(math.Pow(2, float64(attempt))) * time.Second
//             time.Sleep(waitTime)
//             continue
//         }

//         return nil, err
//     }

//     return nil, fmt.Errorf("failed to fetch properties after %d attempts", maxRetries)
// }



// func (c *BookingController) fetchPropertiesForCity(cityName, country string) ([]models.Property, error) {
//     uniqueProperties := make(map[string]models.Property)
//     searchQueries := []string{
//         cityName,
//         fmt.Sprintf("%s hotels", cityName),
//         fmt.Sprintf("%s accommodation", cityName),
//     }

//     for _, query := range searchQueries {
//         encodedQuery := url.QueryEscape(query)
//         apiURL := fmt.Sprintf("https://booking-com18.p.rapidapi.com/stays/auto-complete?query=%s", encodedQuery)

//         properties, err := c.fetchPropertyData(apiURL)
//         if err != nil {
//             continue // Skip errors and proceed with other queries
//         }

//         for _, prop := range properties {
//             if prop.DestID != "" { // Ensure destId exists
//                 uniqueProperties[prop.DestID] = prop
//             }
//         }
//     }

//     result := make([]models.Property, 0, len(uniqueProperties))
//     for _, prop := range uniqueProperties {
//         result = append(result, prop)
//     }

//     return result, nil
// }
// func (c *BookingController) fetchPropertyData(apiURL string) ([]models.Property, error) {
//     req, err := http.NewRequest("GET", apiURL, nil)
//     if err != nil {
//         return nil, err
//     }

//     req.Header.Add("x-rapidapi-host", "booking-com18.p.rapidapi.com")
//     req.Header.Add("x-rapidapi-key", c.rapidAPIKey) // Use the stored RapidAPI key

//     resp, err := c.makeRateLimitedRequest(req)
//     if err != nil {
//         return nil, err
//     }
//     defer resp.Body.Close()

//     body, err := io.ReadAll(resp.Body)
//     if err != nil {
//         return nil, err
//     }

//     if resp.StatusCode == 429 || strings.Contains(string(body), "Too many requests") {
//         return nil, fmt.Errorf("rate limit exceeded")
//     }

//     // Unmarshal the JSON response
//     var response struct {
//         Data []struct {
//             DestID   string `json:"dest_id"`
//             Name     string `json:"name"`
//             // Address  string `json:"address"`
//             CityName string `json:"city_name"`
//             // Add other fields if needed
//         } `json:"data"`
//     }

//     err = json.Unmarshal(body, &response)
//     if err != nil {
//         return nil, err
//     }

//     // Map the API response to your models.Property struct
//     properties := make([]models.Property, 0, len(response.Data))
//     for _, item := range response.Data {
//         properties = append(properties, models.Property{
//             DestID:   item.DestID,
//             Name:     item.Name,
//             // Address:  item.Address,
//             CityName: item.CityName,
//             // Map other fields as necessary
//         })
//     }

//     return properties, nil
// }


// func (c *BookingController) SaveToDatabase() error {
//     // Create a map to store unique city-country pairs
//     uniqueLocations := make(map[string]string)

//     // Iterate through cityProperties to get properties
//     for city := range c.cityProperties {
//         // Find the country for this city
//         var country string
//         for countryName, cities := range c.countryCities {
//             for _, cityName := range cities {
//                 if cityName == city {
//                     country = countryName
//                     break
//                 }
//             }
//         }

//         // Store the unique city-country pair
//         if country != "" {
//             uniqueLocations[city] = country
//         }
//     }

//     // Create a slice to store locations for batch insert
//     var locations []models.Location
//     for city, country := range uniqueLocations {
//         locations = append(locations, models.Location{
//             City:    city,
//             Country: country,
//         })
//     }

//     // Batch insert locations
//     result := conf.DB.CreateInBatches(locations, 100)
//     if result.Error != nil {
//         return fmt.Errorf("failed to save locations: %v", result.Error)
//     }

//     log.Printf("Successfully saved %d unique locations to database", len(locations))
//     return nil
// }

// // ProcessAllHotelDetails processes all hotel details
// func (c *BookingController) ProcessAllHotelDetails() error {
//     log.Println("Starting to process hotel details...")

//     c.mutex.Lock()
//     var allDestIds []string
//     processedIds := make(map[string]bool)

//     // Collect all unique destIds from each city
//     for _, properties := range c.cityProperties {
//         for _, propID := range properties {
//             if !processedIds[propID] {
//                 allDestIds = append(allDestIds, propID)
//                 processedIds[propID] = true
//             }
//         }
//     }
//     c.mutex.Unlock()

//     log.Printf("Found %d unique hotels to process", len(allDestIds))

//     // Create a channel to store results``
//     results := make(chan *models.HotelDetails, len(allDestIds))
//     var wg sync.WaitGroup

//     // Create semaphore for rate limiting
//     semaphore := make(chan struct{}, 1) // Process one request at a time

//     // Process each destId
//     for _, destID := range allDestIds {
//         wg.Add(1)
//         go func(id string) {
//             defer wg.Done()

//             // Acquire semaphore
//             semaphore <- struct{}{}
//             defer func() { <-semaphore }()

//             // Wait for rate limiter
//             err := c.rateLimiter.Wait(context.Background())
//             if err != nil {
//                 log.Printf("Rate limiter error for hotel %s: %v", id, err)
//                 return
//             }

//             details, err := c.fetchHotelDetailsFromAPI(id)
//             if err != nil {
//                 log.Printf("Error fetching details for hotel %s: %v", id, err)
//                 return
//             }

//             results <- details
//         }(destID)
//     }

//     // Close results channel when all goroutines are done
//     go func() {
//         wg.Wait()
//         close(results)
//     }()

//     // Collect and store results
//     hotelDetails := make(map[string]*models.HotelDetails)
//     for detail := range results {
//         if detail != nil {
//             hotelDetails[detail.HotelID] = detail
//             log.Printf("Successfully processed hotel: %s - %s (Bedrooms: %d)", detail.HotelID, detail.PropertyName, detail.Bedrooms)
//         }
//     }

//     // Save results to file
//     outputData, err := json.MarshalIndent(hotelDetails, "", "    ")
//     if err != nil {
//         return fmt.Errorf("error marshaling hotel details: %v", err)
//     }

//     err = os.WriteFile("hotel_details.json", outputData, 0644)
//     if err != nil {
//         return fmt.Errorf("error saving hotel details to file: %v", err)
//     }

//     log.Printf("Successfully processed and saved details for %d hotels", len(hotelDetails))
//     return nil
// }

// func (c *BookingController) fetchHotelDetailsFromAPI(hotelID string) (*models.HotelDetails, error) {
//     url := fmt.Sprintf("https://booking-com18.p.rapidapi.com/stays/detail?hotelId=%s&checkinDate=2025-01-09&checkoutDate=2025-01-23&units=metric", hotelID)
    
//     req, err := http.NewRequest("GET", url, nil)
//     if err != nil {
//         return nil, fmt.Errorf("error creating request: %v", err)
//     }

//     req.Header.Add("x-rapidapi-host", "booking-com18.p.rapidapi.com")
//     req.Header.Add("x-rapidapi-key", c.rapidAPIKey)

//     client := &http.Client{Timeout: 10 * time.Second}
//     resp, err := client.Do(req)
//     if err != nil {
//         return nil, fmt.Errorf("error making request: %v", err)
//     }
//     defer resp.Body.Close()

//     body, err := io.ReadAll(resp.Body)
//     if err != nil {
//         return nil, fmt.Errorf("error reading response: %v", err)
//     }

//     var apiResponse struct {
//         Data struct {
//             HotelID              string `json:"hotel_id"`
//             HotelName            string `json:"hotel_name"`
//             AccommodationTypeName string `json:"accommodation_type_name"`
//             Rooms                map[string]struct {
//                 PrivateBathroomCount int `json:"private_bathroom_count"`
//             } `json:"rooms"`
//             FacilitiesBlock struct {
//                 Facilities []struct {
//                     Name string `json:"name"`
//                 } `json:"facilities"`
//             } `json:"facilities_block"`
// 			BlockCount int `json:"block_count"` // Field for bedroom count
//         } `json:"data"`
//     }

//     if err := json.Unmarshal(body, &apiResponse); err != nil {
//         return nil, fmt.Errorf("error parsing JSON: %v", err)
//     }

//     details := &models.HotelDetails{
//         HotelID:     apiResponse.Data.HotelID,
//         PropertyName: apiResponse.Data.HotelName,
//         Type:        apiResponse.Data.AccommodationTypeName,
// 		Bedrooms:    apiResponse.Data.BlockCount, // Assign block_count to Bedrooms
//         Bathroom:    0, // Default to 0, as we'll extract it from the rooms
// 		Amenities:   make([]models.Facility, 0),
//     }

//     // Extract bathroom count from first room
//     for _, room := range apiResponse.Data.Rooms {
//         details.Bathroom = room.PrivateBathroomCount
//         break
//     }

//     // Extract facilities
//     for _, facility := range apiResponse.Data.FacilitiesBlock.Facilities {
//         if facility.Name != "" {
//             details.Amenities = append(details.Amenities, models.Facility{Name: facility.Name})
//         }
//     }

//     return details, nil
// }



// // Fetch hotel description from the API
// func (c *BookingController) fetchHotelDescriptionFromAPI(hotelID string) ([]models.HotelDetails, error) {
//     url := fmt.Sprintf("https://booking-com18.p.rapidapi.com/stays/get-description?hotelId=%s", hotelID)
    
//     req, err := http.NewRequest("GET", url, nil)
//     if err != nil {
//         return nil, fmt.Errorf("error creating request: %v", err)
//     }

//     req.Header.Add("x-rapidapi-host", "booking-com18.p.rapidapi.com")
//     req.Header.Add("x-rapidapi-key", c.rapidAPIKey)

//     client := &http.Client{Timeout: 10 * time.Second}
//     resp, err := client.Do(req)
//     if err != nil {
//         return nil, fmt.Errorf("error making request: %v", err)
//     }
//     defer resp.Body.Close()

//     body, err := io.ReadAll(resp.Body)
//     if err != nil {
//         return nil, fmt.Errorf("error reading response: %v", err)
//     }
//     var apiResponse struct {
//         Data struct {
//             Description string `json:"description"`
//             PropertyName string `json:"property_name"` // Assuming this field exists in the API response
//         } `json:"data"`
//     }

//     if err := json.Unmarshal(body, &apiResponse); err != nil {
//         return nil, fmt.Errorf("error parsing JSON: %v", err)
//     }

//     // Create an instance of HotelDetails
//     hotelDetails := models.HotelDetails{
//         HotelID:     hotelID,
//         Description: apiResponse.Data.Description,
//         PropertyName: apiResponse.Data.PropertyName, // Assuming this field exists
//     }

//     return []models.HotelDetails{hotelDetails}, nil // Return as a slice

// }

// // Process all hotel descriptions
// // Process all hotel descriptions
// func (c *BookingController) ProcessAllHotelDescriptions() error {
//     log.Println("Starting to process hotel descriptions...")

//     c.mutex.Lock()
//     var allDestIds []string
//     processedIds := make(map[string]bool)

//     // Collect all unique destIds from each city
//     for city, properties := range c.cityProperties {
//         log.Printf("Processing properties for city: %s", city)
//         for _, propID := range properties {
//             if !processedIds[propID] {
//                 allDestIds = append(allDestIds, propID)
//                 processedIds[propID] = true
//             }
//         }
//     }
//     c.mutex.Unlock()

//     log.Printf("Found %d unique hotels to process", len(allDestIds))

//     // Create a channel to store results
//     results := make(chan []models.HotelDetails, len(allDestIds))
//     var wg sync.WaitGroup

//     // Create semaphore for rate limiting
//     semaphore := make(chan struct{}, 1) // Process one request at a time

//     // Process each destId
//     for _, destID := range allDestIds {
//         wg.Add(1)
//         go func(id string) {
//             defer wg.Done()

//             // Acquire semaphore
//             semaphore <- struct{}{}
//             defer func() { <-semaphore }()

//             // Wait for rate limiter
//             err := c.rateLimiter.Wait(context.Background())
//             if err != nil {
//                 log.Printf("Rate limiter error for hotel %s: %v", id, err)
//                 return
//             }

//             details, err := c.fetchHotelDescriptionFromAPI(id)
//             if err != nil {
//                 log.Printf("Error fetching description for hotel %s: %v", id, err)
//                 return
//             }

//             results <- details
//         }(destID)
//     }

//     // Close results channel when all goroutines are done
//     go func() {
//         wg.Wait()
//         close(results)
//     }()

//     // Collect and store results
//     hotelDescriptions := make(map[string]*models.HotelDetails)
//     for detailSlice := range results {
//         for _, detail := range detailSlice { // Iterate over the slice
//             hotelDescriptions[detail.HotelID] = &detail // Store the pointer to detail
//             log.Printf("Successfully processed hotel: %s - %s", detail.HotelID, detail.PropertyName)
//         }
//     }

//     // Save results to file
//     outputData, err := json.MarshalIndent(hotelDescriptions, "", "    ")
//     if err != nil {
//         return fmt.Errorf("error marshaling hotel descriptions: %v", err)
//     }

//     err = os.WriteFile("hotel_descriptions.json", outputData, 0644)
//     if err != nil {
//         return fmt.Errorf("error saving hotel descriptions to file: %v", err)
//     }

//     log.Printf("Successfully processed and saved descriptions for %d hotels", len(hotelDescriptions))
//     return nil
// }

// //!saving Rental Property
// func (c *BookingController) SaveRentalPropertiesToDatabase() error {
//     // Create a slice to store rental properties
//     var rentalProperties []models.RentalProperty

//     // Iterate through cityProperties to get property details
//     for city, properties := range c.cityProperties {
//         for _, propertyName := range properties {
//             // Fetch additional details for each property
//             propertyDetails, err := c.fetchPropertyDetails(propertyName, city)
//             if err != nil {
//                 log.Printf("Error fetching details for property %s in city %s: %v", propertyName, city, err)
//                 continue
//             }

//             // Create a RentalProperty instance
//             rentalProperty := models.RentalProperty{
//                 PropertyName: propertyDetails.PropertyName,
//                 Type:         propertyDetails.Type,
//                 Bedrooms:     propertyDetails.Bedrooms,
//                 Bathrooms:    propertyDetails.Bathrooms,
//                 Amenities:    propertyDetails.Amenities,
//                 LocationID:   propertyDetails.LocationID,
//             }

//             // Append to the slice
//             rentalProperties = append(rentalProperties, rentalProperty)
//         }
//     }

//     // Batch insert rental properties
//     result := conf.DB.CreateInBatches(rentalProperties, 100)
//     if result.Error != nil {
//         return fmt.Errorf("failed to save rental properties: %v", result.Error)
//     }

//     log.Printf("Successfully saved %d rental properties to database", len(rentalProperties))
//     return nil
// }

// // Example function to fetch property details (you need to implement this)
// func (c *BookingController) fetchPropertyDetails(propertyName, city string) (*models.RentalProperty, error) {
//     // Implement the logic to fetch property details based on propertyName and city
//     // This is a placeholder function and should return the necessary details
//     return &models.RentalProperty{
//         PropertyName: propertyName,
//         Type:         "Apartment", // Example type
//         Bedrooms:     2,           // Example number of bedrooms
//         Bathrooms:    1,           // Example number of bathrooms
//         Amenities:    []string{"WiFi", "Pool"}, // Example amenities
//         LocationID:   1,           // Example location ID (you need to determine how to get this)
//     }, nil
// }



// // ListCities returns the list of available cities
// func (c *BookingController) ListCities() {
//     var cities []string
//     conf.DB.Model(&models.RentalProperty{}).Distinct().Pluck("city", &cities)
//     c.Data["json"] = cities
//     c.ServeJSON()
// }


// //!Fetching images:
// // FetchHotelImagesFromAPI fetches images for a given hotelId and categorizes them
// func (c *BookingController) fetchHotelImagesFromAPI(hotelID string) (*models.CategorizedImages, error) {
//     url := fmt.Sprintf("https://booking-com18.p.rapidapi.com/stays/get-photos?hotelId=%s", hotelID)
    
//     req, err := http.NewRequest("GET", url, nil)
//     if err != nil {
//         return nil, fmt.Errorf("error creating request: %v", err)
//     }

//     req.Header.Add("x-rapidapi-host", "booking-com18.p.rapidapi.com")
//     req.Header.Add("x-rapidapi-key", c.rapidAPIKey)

//     client := &http.Client{Timeout: 10 * time.Second}
//     resp, err := client.Do(req)
//     if err != nil {
//         return nil, fmt.Errorf("error making request: %v", err)
//     }
//     defer resp.Body.Close()

//     body, err := io.ReadAll(resp.Body)
//     if err != nil {
//         return nil, fmt.Errorf("error reading response: %v", err)
//     }

//     var apiResponse struct {
//         Data []struct {
//             ID    int      `json:"id"`
//             Tag   string   `json:"tag"`
//             Images []string `json:"images"`
//         } `json:"data"`
//     }

//     if err := json.Unmarshal(body, &apiResponse); err != nil {
//         return nil, fmt.Errorf("error parsing JSON: %v", err)
//     }

//     categorizedImages := &models.CategorizedImages{}
//     for _, item := range apiResponse.Data {
//         switch item.Tag {
//         case "Property building":
//             categorizedImages.PropertyBuilding = append(categorizedImages.PropertyBuilding, item.Images...)
//         case "Property":
//             categorizedImages.Property = append(categorizedImages.Property, item.Images...)
//         case "Room":
//             categorizedImages.Room = append(categorizedImages.Room, item.Images...)
//         // Add more categories as needed
//         default:
//             // Handle other categories or ignore
//         }
//     }

//     return categorizedImages, nil
// }

// // Function to process all hotel images
// func (c *BookingController) ProcessAllHotelImages() error {
//     log.Println("Starting to process hotel images...")

//     c.mutex.Lock()
//     var allDestIds []string
//     processedIds := make(map[string]bool)

//     // Collect all unique destIds from each city
//     for _, properties := range c.cityProperties {
//         for _, propID := range properties {
//             if !processedIds[propID] {
//                 allDestIds = append(allDestIds, propID)
//                 processedIds[propID] = true
//             }
//         }
//     }
//     c.mutex.Unlock()

//     log.Printf("Found %d unique hotels to process images for", len(allDestIds))

//     // Create a channel to store results
//     results := make(chan struct {
//         HotelID          string
//         CategorizedImages *models.CategorizedImages
//         Err              error
//     }, len(allDestIds))
//     var wg sync.WaitGroup

//     // Create semaphore for rate limiting
//     semaphore := make(chan struct{}, 1) // Process one request at a time

//     // Process each destId
//     for _, destID := range allDestIds {
//         wg.Add(1)
//         go func(id string) {
//             defer wg.Done()

//             // Acquire semaphore
//             semaphore <- struct{}{}
//             defer func() { <-semaphore }()

//             // Wait for rate limiter
//             err := c.rateLimiter.Wait(context.Background())
//             if err != nil {
//                 log.Printf("Rate limiter error for hotel %s: %v", id, err)
//                 return
//             }

//             categorizedImages, err := c.fetchHotelImagesFromAPI(id)
//             results <- struct {
//                 HotelID          string
//                 CategorizedImages *models.CategorizedImages
//                 Err              error
//             }{
//                 HotelID:          id,
//                 CategorizedImages: categorizedImages,
//                 Err:              err,
//             }
//         }(destID)
//     }

//     // Close results channel when all goroutines are done
//     go func() {
//         wg.Wait()
//         close(results)
//     }()

//     // Collect and store results
//     hotelImages := make(map[string]*models.CategorizedImages)
//     for result := range results {
//         if result.Err != nil {
//             log.Printf("Error fetching images for hotel %s: %v", result.HotelID, result.Err)
//             continue
//         }
//         hotelImages[result.HotelID] = result.CategorizedImages
//         log.Printf("Successfully processed images for hotel: %s", result.HotelID)
//     }

//     // Save results to file or database as needed
//     // Here you would implement the logic to save hotelImages

//     log.Printf("Successfully processed and saved images for %d hotels", len(hotelImages))
//     return nil
// }

// //!for rating and review:
// // FetchHotelRatingAndReviewsFromAPI fetches rating and review count for a given hotelId
// func (c *BookingController) fetchHotelRatingAndReviewsFromAPI(hotelID string) (*models.RatingAndReviews, error) {
//     url := fmt.Sprintf("https://booking-com18.p.rapidapi.com/stays/review-featured?hotelId=%s", hotelID)
    
//     req, err := http.NewRequest("GET", url, nil)
//     if err != nil {
//         return nil, fmt.Errorf("error creating request: %v", err)
//     }

//     req.Header.Add("x-rapidapi-host", "booking-com18.p.rapidapi.com")
//     req.Header.Add("x-rapidapi-key", c.rapidAPIKey)

//     client := &http.Client{Timeout: 10 * time.Second}
//     resp, err := client.Do(req)
//     if err != nil {
//         return nil, fmt.Errorf("error making request: %v", err)
//     }
//     defer resp.Body.Close()

//     body, err := io.ReadAll(resp.Body)
//     if err != nil {
//         return nil, fmt.Errorf("error reading response: %v", err)
//     }

//     var apiResponse struct {
//         VpmFavorableReviewCount int `json:"vpm_favorable_review_count"`
//         VpmFeaturedReviews      []struct {
//             AverageScoreOutOf10 float64 `json:"average_score_out_of_10"`
//         } `json:"vpm_featured_reviews"`
//     }

//     if err := json.Unmarshal(body, &apiResponse); err != nil {
//         return nil, fmt.Errorf("error parsing JSON: %v", err)
//     }

//     ratingAndReviews := &models.RatingAndReviews{
//         VpmFavorableReviewCount: apiResponse.VpmFavorableReviewCount,
//     }

//     if len(apiResponse.VpmFeaturedReviews) > 0 {
//         ratingAndReviews.AverageScoreOutOf10 = apiResponse.VpmFeaturedReviews[0].AverageScoreOutOf10
//     }

//     return ratingAndReviews, nil
// }

// // Function to process all hotel ratings and reviews
// func (c *BookingController) ProcessAllHotelRatingsAndReviews() error {
//     log.Println("Starting to process hotel ratings and reviews...")

//     c.mutex.Lock()
//     var allDestIds []string
//     processedIds := make(map[string]bool)

//     // Collect all unique destIds from each city
//     for _, properties := range c.cityProperties {
//         for _, propID := range properties {
//             if !processedIds[propID] {
//                 allDestIds = append(allDestIds, propID)
//                 processedIds[propID] = true
//             }
//         }
//     }
//     c.mutex.Unlock()

//     log.Printf("Found %d unique hotels to process ratings and reviews for", len(allDestIds))

//     // Create a channel to store results
//     results := make(chan struct {
//         HotelID          string
//         RatingAndReviews *models.RatingAndReviews
//         Err              error
//     }, len(allDestIds))
//     var wg sync.WaitGroup

//     // Create semaphore for rate limiting
//     semaphore := make(chan struct{}, 1) // Process one request at a time

//     // Process each destId
//     for _, destID := range allDestIds {
//         wg.Add(1)
//         go func(id string) {
//             defer wg.Done()

//             // Acquire semaphore
//             semaphore <- struct{}{}
//             defer func() { <-semaphore }()

//             // Wait for rate limiter
//             err := c.rateLimiter.Wait(context.Background())
//             if err != nil {
//                 log.Printf("Rate limiter error for hotel %s: %v", id, err)
//                 return
//             }

//             ratingAndReviews, err := c.fetchHotelRatingAndReviewsFromAPI(id)
//             results <- struct {
//                 HotelID          string
//                 RatingAndReviews *models.RatingAndReviews
//                 Err              error
//             }{
//                 HotelID:          id,
//                 RatingAndReviews: ratingAndReviews,
//                 Err:              err,
//             }
//         }(destID)
//     }

//     // Close results channel when all goroutines are done
//     go func() {
//         wg.Wait()
//         close(results)
//     }()

//     // Collect and store results
//     hotelRatingsAndReviews := make(map[string]*models.RatingAndReviews)
//     for result := range results {
//         if result.Err != nil {
//             log.Printf("Error fetching ratings and reviews for hotel %s: %v", result.HotelID, result.Err)
//             continue
//         }
//         hotelRatingsAndReviews[result.HotelID] = result.RatingAndReviews
//         log.Printf("Successfully processed ratings and reviews for hotel: %s", result.HotelID)
//     }

//     // Save results to file or database as needed
//     // Here you would implement the logic to save hotelRatingsAndReviews

//     log.Printf("Successfully processed and saved ratings and reviews for %d hotels", len(hotelRatingsAndReviews))
//     return nil
// }

// //! SavePropertyDetails saves property details including description, images, rating, and reviews to the PropertyDetails table
// func (c *BookingController) SavePropertyDetails(propertyID int, description string, images []string, rating float64, reviews []string) error {
//     // Create a PropertyDetails instance
//     propertyDetails := models.PropertyDetails{
//         PropertyID:        propertyID,
//         Description:       description,
//         Images:            images,
//         Rating:            rating,
//         Reviews:           reviews,
//     }

//     // Save to the database
//     result := conf.DB.Create(&propertyDetails)
//     if result.Error != nil {
//         return result.Error
//     }

//     log.Printf("Successfully saved property details for property ID %d", propertyID)
//     return nil
// }

// // Example usage of SavePropertyDetails
// func (c *BookingController) ProcessPropertyDetails() {
//     // Example data
//     propertyID := 1
//     description := "A beautiful property with an amazing view."
//     images := []string{
//         "/xdata/images/hotel/square60/484656765.jpg?k=aaa04cb98edd1c367ceb4bca7b73867d501154a2ea4f3ef39a3ddcfeeb1e8fe8&o=",
//         "/xdata/images/hotel/max1024x768/484656765.jpg?k=aaa04cb98edd1c367ceb4bca7b73867d501154a2ea4f3ef39a3ddcfeeb1e8fe8&o=",
//     }
//     rating := 9.0
//     reviews := []string{
//         "Staff was really helpful and polite.",
//         "Rooms very clean and big enough, amazing view.",
//         "Proximity to core NYC.",
//     }

//     err := c.SavePropertyDetails(propertyID, description, images, rating, reviews)
//     if err != nil {
//         log.Fatalf("Failed to save property details: %v", err)
//     }
// }
// //! ListProperties handles the /v1/property/list endpoint
// func (c *BookingController) ListProperties() {
//     location := c.GetString("location")
//     if location == "" {
//         c.Ctx.Output.SetStatus(400)
//         c.Ctx.Output.Body([]byte("Location is required"))
//         return
//     }

//     var properties []models.RentalProperty
//     result := conf.DB.Preload("Location").Where("location_id IN (SELECT id FROM locations WHERE city = ?)", location).Find(&properties)
//     if result.Error != nil {
//         c.Ctx.Output.SetStatus(500)
//         c.Ctx.Output.Body([]byte("Failed to fetch properties"))
//         return
//     }

//     c.Data["json"] = properties
//     c.ServeJSON()
    
// }

// func (c *BookingController) Index() {
//     c.TplName = "index.html"
// }

// //!GET property details:
// func (c *BookingController) GetPropertyDetails() {
//     propertyIdStr := c.Ctx.Input.Param(":propertyId")
//     propertyId, err := strconv.Atoi(propertyIdStr)
//     if err != nil {
//         c.Ctx.Output.SetStatus(400)
//         c.Ctx.Output.Body([]byte("Invalid property ID"))
//         return
//     }

//     var property models.Property
//     result := conf.DB.First(&property, propertyId)
//     if result.Error != nil {
//         c.Ctx.Output.SetStatus(404)
//         c.Ctx.Output.Body([]byte("Property not found"))
//         return
//     }

//     c.Data["json"] = property
//     c.ServeJSON()
// }