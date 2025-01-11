// models/booking.go
package models

// City represents a city from the Booking.com API
type City struct {
	CC1        string  `json:"cc1"`
	ImageURL   string  `json:"image_url"`
	Longitude  float64 `json:"longitude"`
	CityName   string  `json:"city_name"`
	DestID     string  `json:"dest_id"`
	Timezone   string  `json:"timezone"`
	Hotels     int     `json:"hotels"`
	Label      string  `json:"label"`
	Country    string  `json:"country"`
	Region     string  `json:"region"`
	DestType   string  `json:"dest_type"`
	Name       string  `json:"name"`
	Latitude   float64 `json:"latitude"`
	Type       string  `json:"type"`
}

// Property represents a property from the Booking.com API
type Property struct {
	UFI               int64   `json:"ufi"`
	CheckoutDate      string  `json:"checkoutDate"`
	ReviewScoreWord   string  `json:"reviewScoreWord"`
	Longitude         float64 `json:"longitude"`
	IsPreferred       bool    `json:"isPreferred"`
	CountryCode       string  `json:"countryCode"`
	Latitude          float64 `json:"latitude"`
	WishlistName      string  `json:"wishlistName"`
	Name              string  `json:"name"`
	PropertyClass     float64 `json:"accuratePropertyClass"`
	DestID            string  `json:"dest_id"`
	CityName          string  `json:"city_name"`
	Country           string  `json:"country"`

	
}

// CityKey represents a unique identifier for a city
type CityKey struct {
	Name    string
	Country string
}

// HotelDetails represents the full hotel information
type HotelDetails struct {
    HotelID    string   `json:"hotel_id"`
    PropertyName string  `json:"property_name"`
    Type        string  `json:"type"`
	Bedrooms    int     `json:"bedrooms"` // New field for block_count
    Bathroom    int     `json:"bathroom"`
    Amenities   []Facility `json:"amenities"`
	Description  string `json:"description"`
}

// Facility represents a hotel facility/amenity
type Facility struct {
    Name string `json:"name"`
}


type CategorizedImages struct {
    PropertyBuilding []string `json:"property_building"`
    Property         []string `json:"property"`
    Room             []string `json:"room"`
    // Add more categories as needed
}

type RatingAndReviews struct {
    AverageScoreOutOf10       float64 `json:"average_score_out_of_10"`
    VpmFavorableReviewCount   int     `json:"vpm_favorable_review_count"`
}