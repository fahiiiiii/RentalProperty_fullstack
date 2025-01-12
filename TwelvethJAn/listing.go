package models

// import (
//     "encoding/json"
// )

type Listing struct {
    ID          int      `json:"id"`
    Title       string   `json:"title"`
    Price       float64  `json:"price"`
    Rating      float64  `json:"rating"`
    Reviews     int      `json:"reviews"`
    Amenities   []string `json:"amenities"`
    Location    string   `json:"location"`
    Images      []string `json:"images"`
}

func GetListings() []Listing {
    return []Listing{
        {
            ID:          1,
            Title:       "Luxury Villa in Palm Jumeirah",
            Price:       565782,
            Rating:      9.5,
            Reviews:     15,
            Amenities:   []string{"Air Conditioner", "Swimming Pool", "Cold Storage"},
            Location:    "Dubai • Palm Jumeirah",
            Images:      []string{"villa1.jpg", "villa2.jpg", "villa3.jpg"},
        },
        {
            ID:          2,
            Title:       "Modern Apartment in Capital Bay",
            Price:       223000,
            Rating:      8.8,
            Reviews:     24,
            Amenities:   []string{"Air Conditioner", "City View", "Hot Tub"},
            Location:    "Dubai • Business Bay",
            Images:      []string{"apt1.jpg", "apt2.jpg", "apt3.jpg"},
        },
    }
}