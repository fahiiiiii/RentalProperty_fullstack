package routers

import (
    beego "github.com/beego/beego/v2/server/web"
    "rental_view/controllers"
)

func init() {
    // Property routes
    beego.Router("/", &controllers.MainController{})
    beego.Router("/properties", &controllers.PropertyController{}, "get:List")
    beego.Router("/properties/filter", &controllers.PropertyController{}, "post:Filter")
    beego.Router("/properties/types", &controllers.PropertyController{}, "get:GetPropertyTypes")
}


// package routers

// import (
//     beego "github.com/beego/beego/v2/server/web"
//     "rental_view/controllers"
// )

// func init() {
// 	// Serve static files
//     beego.SetStaticPath("/static", "static")
//     // V1 API namespace
//     ns := beego.NewNamespace("/v1",
//         beego.NSNamespace("/property",
//             // Property list route
//             beego.NSRouter("/list", &controllers.PropertyController{}, "get:List"),
            
//             // Additional property-related routes
//             beego.NSRouter("/filter", &controllers.PropertyController{}, "post:Filter"),
//         ),
//     )
    
//     // Add namespace to beego
//     beego.AddNamespace(ns)

//     // Default route
//     beego.Router("/", &controllers.MainController{})
// }