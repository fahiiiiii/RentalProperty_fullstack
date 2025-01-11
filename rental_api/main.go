// main.go
package main

import (
	_ "rental_api/routers"
	"rental_api/models"


	"github.com/beego/beego/v2/server/web"
)

func main() {
	// Initialize the database
	conf.InitDB()

	// Serve static files
	web.SetStaticPath("/static", "static")

	// Start the application
	web.Run()
}




// package main

// import (
// 	_ "rental_api/routers"

// 	beego "github.com/beego/beego/v2/server/web"
// )

// func main() {
// 	if beego.BConfig.RunMode == "dev" {
// 		beego.BConfig.WebConfig.DirectoryIndex = true
// 		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
// 	}
// 	beego.Run()
// }
