package models

import (
	"github.com/beego/beego/v2/client/orm"
)

// Location represents a geographical location
type Location struct {
	ID      int    `orm:"auto;pk" json:"id"`
	City    string `orm:"size(100)" json:"city"`
	Country string `orm:"size(100)" json:"country"`
	Properties []*RentalProperty `orm:"reverse(many)" json:"properties,omitempty"`
}

func (l *Location) TableName() string {
	return "locations"
}

// GetAllLocations retrieves all locations from the database
func GetAllLocations() ([]Location, error) {
	var locations []Location
	o := orm.NewOrm()
	_, err := o.QueryTable("locations").All(&locations)
	return locations, err
}

// GetLocationByID retrieves a location by its ID
func GetLocationByID(id int) (*Location, error) {
	location := &Location{ID: id}
	o := orm.NewOrm()
	err := o.Read(location)
	return location, err
}