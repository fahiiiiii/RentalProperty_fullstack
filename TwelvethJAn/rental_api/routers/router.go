// @APIVersion 1.0.0
// @Title Booking API
// @Description This API manages bookings using the Beego framework
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
// package routers

// import (
// 	"rental_api/controllers"

// 	beego "github.com/beego/beego/v2/server/web"
// )

// func init() {
// 	ns := beego.NewNamespace("/v1",
// 		beego.NSNamespace("/booking", // Changed from /object to /booking
// 			beego.NSInclude(
// 				&controllers.BookingController{}, // Ensure you have this controller implemented
// 			),
// 		),
		
// 	)
// 	beego.AddNamespace(ns)
// }

package routers

import (
	"rental_api/controllers"
	"github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/filter/cors"
)

func init() {
	// Enable CORS
	web.InsertFilter("*", web.BeforeRouter, cors.Allow(&cors.Options{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Access-Control-Allow-Origin", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin"},
		AllowCredentials: true,
	}))

	// API version namespace
	ns := web.NewNamespace("/v1",
		// Booking routes
		web.NSNamespace("/booking",
			// Main booking operations
			web.NSRouter("/", &controllers.BookingController{}, "get:Get;post:Post"),
			web.NSRouter("/:propertyId", &controllers.BookingController{}, "get:GetPropertyDetails"),
			
			// Cities
			web.NSRouter("/cities", &controllers.BookingController{}, "get:ListCities"),
			web.NSRouter("/cities/process", &controllers.BookingController{}, "post:ProcessAllCities"),
			
			// Properties
			web.NSRouter("/properties", &controllers.BookingController{}, "get:ListProperties"),
			web.NSRouter("/properties/process", &controllers.BookingController{}, "post:ProcessAllProperties"),
			
			// Property Details
			web.NSRouter("/properties/:propertyId/details", &controllers.BookingController{}, 
				"get:GetPropertyDetails;post:ProcessAllHotelDetails"),
			
			// Images
			web.NSRouter("/properties/:propertyId/images", &controllers.BookingController{}, 
				"get:GetHotelImages;post:ProcessAllHotelImages"),
			
			// Reviews and Ratings
			web.NSRouter("/properties/:propertyId/reviews", &controllers.BookingController{}, 
				"get:GetHotelReviews;post:ProcessAllHotelRatingsAndReviews"),
			
			// Database Operations
			web.NSRouter("/save", &controllers.BookingController{}, "post:SaveToDatabase"),
			web.NSRouter("/save/properties", &controllers.BookingController{}, "post:SaveRentalPropertiesToDatabase"),
			web.NSRouter("/save/details", &controllers.BookingController{}, "post:SavePropertyDetails"),
		),

		// Health check
		web.NSRouter("/health", &controllers.BookingController{}, "get:HealthCheck"),
	)

	// Register namespace
	web.AddNamespace(ns)

	// Static files
	web.SetStaticPath("/swagger", "swagger")
	web.SetStaticPath("/static", "static")

	// Custom error handlers
	web.ErrorHandler("404", handle404)
	web.ErrorHandler("500", handle500)
}

// Error handlers
func handle404(rw web.ResponseWriter, r *web.Request) {
	rw.WriteHeader(404)
	rw.Write([]byte("Resource not found"))
}

func handle500(rw web.ResponseWriter, r *web.Request) {
	rw.WriteHeader(500)
	rw.Write([]byte("Internal server error"))
}