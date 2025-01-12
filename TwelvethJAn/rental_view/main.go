package main

import (
    "fmt"
    "html/template"
    "strconv"

    beego "github.com/beego/beego/v2/server/web"
    _ "rental_view/routers"
)

// Custom template functions
func formatNumber(v interface{}) string {
    switch num := v.(type) {
    case float64:
        return strconv.FormatFloat(num, 'f', 0, 64)
    case int:
        return strconv.Itoa(num)
    case int64:
        return strconv.FormatInt(num, 10)
    default:
        return fmt.Sprintf("%v", v)
    }
}

// Custom subtraction function for template
func sub(a, b int) int {
    return a - b
}

// Subtract 1 from length
func minus1(v int) int {
    return v - 1
}

// SafeHTML converts a string to a template.HTML type
func safeHTML(s string) template.HTML {
    return template.HTML(s)
}

func init() {
    // Register custom template functions
    beego.AddFuncMap("formatNumber", formatNumber)
    beego.AddFuncMap("sub", sub)
    beego.AddFuncMap("minus1", minus1)
    beego.AddFuncMap("safehtml", safeHTML)
}


func main() {
   

	beego.BConfig.WebConfig.StaticDir["/static"] = "static"
    beego.Run()
}

// package main

// import (
//     // "fmt"
//     // "log"

//     beego "github.com/beego/beego/v2/server/web"
//     _ "rental_view/routers"
//     "rental_view/models"
// )

// // Custom template function to subtract
// func SubFunc(a, b int) int {
//     return a - b
// }

// func main() {
//     // Initialize database connection
//     models.InitDB()

//     // Add custom template function
//     beego.AddFuncMap("sub", SubFunc)

//     // Run the application
//     beego.Run()
// }