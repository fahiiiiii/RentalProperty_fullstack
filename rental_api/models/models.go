// models/models.go
package models

import (
     "github.com/lib/pq"
    "gorm.io/gorm"
)

type Location struct {
    gorm.Model
    City     string `gorm:"not null"`
    Country  string `gorm:"not null"`
    // Property string `gorm:"not null"`
}

// type RentalProperty struct {
//     ID         int      `db:"id" json:"id"`
//     Name       string   `db:"name" json:"name"`
//     Type       string   `db:"type" json:"type"`
//     Bedrooms   int      `db:"bedrooms" json:"bedrooms"`
//     Bathrooms  int      `db:"bathrooms" json:"bathrooms"`
//     Amenities  []string `db:"amenities" json:"amenities"`
//     LocationID int      `db:"location_id" json:"location_id"`
// }
type RentalProperty struct {
    ID           int            `db:"id" json:"id"`
    PropertyName string         `db:"property_name" json:"property_name"`
    Type         string         `db:"type" json:"type"`
    Bedrooms     int            `db:"bedrooms" json:"bedrooms"`
    Bathrooms    int            `db:"bathrooms" json:"bathrooms"`
    Amenities    pq.StringArray `gorm:"type:text[]" db:"amenities" json:"amenities"`
    LocationID   int            `db:"location_id" json:"location_id"`
}
type PropertyDetail struct {
    ID          int                    `db:"id" json:"id"`
    PropertyID  int                    `db:"property_id" json:"property_id"`
    Description string                 `db:"description" json:"description"`
    Images      []string               `db:"images" json:"images"`
    OtherDetails map[string]interface{} `db:"other_details" json:"other_details"`
}


type PropertyDetails struct {
    gorm.Model
    PropertyID          int            `gorm:"not null" json:"property_id"`
    Description         string         `gorm:"type:text" json:"description"`
    Images              pq.StringArray `gorm:"type:text[]" json:"images"`
    Rating              float64        `gorm:"type:double precision" json:"rating"`
    Reviews             pq.StringArray `gorm:"type:text[]" json:"reviews"`
}