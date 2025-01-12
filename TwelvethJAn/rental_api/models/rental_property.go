package models

import (
	"github.com/beego/beego/v2/client/orm"
)

// RentalProperty represents a rental property
type RentalProperty struct {
	ID          int      `orm:"auto;pk" json:"id"`
	LocationID  int      `orm:"column(location_id)" json:"location_id"`
	Location    *Location `orm:"rel(fk)" json:"location,omitempty"`
	PropertyName string   `orm:"size(200)" json:"property_name"`
	Type        string   `orm:"size(50)" json:"type"`
	Bedrooms    int      `json:"bedrooms"`
	Bathrooms   int      `json:"bathrooms"`
	Amenities   []string `orm:"type(json)" json:"amenities"`
	Details     *PropertyDetails `orm:"reverse(one)" json:"details,omitempty"`
	Images      []*Image    `orm:"reverse(many)" json:"images,omitempty"`
	Reviews     []*Review   `orm:"reverse(many)" json:"reviews,omitempty"`
}

func (rp *RentalProperty) TableName() string {
	return "rental_properties"
}

// GetAllProperties retrieves all rental properties
func GetAllProperties() ([]RentalProperty, error) {
	var properties []RentalProperty
	o := orm.NewOrm()
	_, err := o.QueryTable("rental_properties").RelatedSel().All(&properties)
	return properties, err
}

// GetPropertiesByLocation retrieves properties by location
func GetPropertiesByLocation(location string) ([]RentalProperty, error) {
	var properties []RentalProperty
	o := orm.NewOrm()
	_, err := o.QueryTable("rental_properties").
		Filter("Location__City", location).
		RelatedSel().
		All(&properties)
	return properties, err
}