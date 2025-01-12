// routers/router.go
package routers

import (
    // "api_rental/controllers"
    "rental_api/controllers"
    beego "github.com/beego/beego/v2/server/web"
)

func init() {
    // Define CRUD routes for properties
    beego.Router("/", &controllers.BookingController{}, "get:Index")
    beego.Router("/v1/property", &controllers.BookingController{}, "post:Create")
    beego.Router("/v1/property/:id", &controllers.BookingController{}, "get:Get;put:Update;delete:Delete")
    beego.Router("/v1/properties", &controllers.BookingController{}, "get:List")

    // Additional routes for specific operations
    beego.Router("/v1/cities/process", &controllers.BookingController{}, "get:ProcessAllCities")
    beego.Router("/v1/property/details/:id", &controllers.BookingController{}, "put:ProcessPropertyDetails")
}



// // @APIVersion 1.0.0
// // @Title beego Test API
// // @Description beego has a very cool tools to autogenerate documents for your API
// // @Contact astaxie@gmail.com
// // @TermsOfServiceUrl http://beego.me/
// // @License Apache 2.0
// // @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
// package routers

// import (
// 	"rental_api/controllers"

// 	beego "github.com/beego/beego/v2/server/web"
// )

// func init() {
// 	ns := beego.NewNamespace("/v1",
// 		beego.NSNamespace("/object",
// 			beego.NSInclude(
// 				&controllers.ObjectController{},
// 			),
// 		),
// 		beego.NSNamespace("/user",
// 			beego.NSInclude(
// 				&controllers.UserController{},
// 			),
// 		),
// 	)
// 	beego.AddNamespace(ns)
// }
