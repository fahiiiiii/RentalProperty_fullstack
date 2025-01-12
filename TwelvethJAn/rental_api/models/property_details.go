package models

// import (
// 	"github.com/beego/beego/v2/client/orm"
// )

// PropertyDetails represents detailed information about a property
type PropertyDetails struct {
	ID          int    `orm:"auto;pk" json:"id"`
	PropertyID  int    `orm:"unique" json:"property_id"`
	Property    *RentalProperty `orm:"rel(one)" json:"property,omitempty"`
	Description string `orm:"type(text)" json:"description"`
	Rating      float64 `json:"rating"`
}

// Image represents a property image
type Image struct {
	ID         int    `orm:"auto;pk" json:"id"`
	PropertyID int    `json:"property_id"`
	Property   *RentalProperty `orm:"rel(fk)" json:"property,omitempty"`
	URL        string `orm:"size(500)" json:"url"`
	Category   string `orm:"size(50)" json:"category"` // e.g., "exterior", "interior", "room"
}

// Review represents a property review
type Review struct {
	ID         int    `orm:"auto;pk" json:"id"`
	PropertyID int    `json:"property_id"`
	Property   *RentalProperty `orm:"rel(fk)" json:"property,omitempty"`
	Rating     float64 `json:"rating"`
	Comment    string  `orm:"type(text)" json:"comment"`
	UserName   string  `orm:"size(100)" json:"user_name"`
	CreatedAt  string  `orm:"auto_now_add;type(datetime)" json:"created_at"`
}

func (pd *PropertyDetails) TableName() string {
	return "property_details"
}

func (i *Image) TableName() string {
	return "images"
}

func (r *Review) TableName() string {
	return "reviews"
}