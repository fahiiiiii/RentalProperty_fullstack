package models

import (
    "database/sql"
    "fmt"
    
    _ "github.com/lib/pq"
    beego "github.com/beego/beego/v2/server/web"
)

// RentalProperty represents a rental property
type RentalProperty struct {
    ID           int      `json:"id"`
    PropertyName string   `json:"title"`
    Type         string   `json:"type"`
    Bedrooms     int      `json:"bedrooms"`
    Bathrooms    int      `json:"bathrooms"`
    Price        float64  `json:"price"`
    Location     string   `json:"location"`
    Amenities    []string `json:"amenities"`
    Images       []string `json:"images"`
    Rating       float64  `json:"rating"`
    Reviews      int      `json:"reviews"`
}

func GetRentalProperties(city string, page, pageSize int) ([]RentalProperty, int, error) {
    // Calculate offset
    offset := (page - 1) * pageSize

    // Count total properties
    countQuery := `
        SELECT COUNT(*) 
        FROM rental_properties 
        WHERE LOWER(location) LIKE LOWER($1)
    `
    var totalCount int
    err := DB.QueryRow(countQuery, "%"+city+"%").Scan(&totalCount)
    if err != nil {
        return nil, 0, fmt.Errorf("error counting properties: %v", err)
    }

    // Query to fetch paginated properties
    query := `
        SELECT 
            id, 
            property_name, 
            type, 
            bedrooms, 
            bathrooms, 
            price, 
            location, 
            amenities, 
            images,
            rating,
            reviews
        FROM 
            rental_properties 
        WHERE 
            LOWER(location) LIKE LOWER($1)
        ORDER BY 
            created_at DESC
        LIMIT $2 OFFSET $3
    `

    // Execute query
    rows, err := DB.Query(query, "%"+city+"%", pageSize, offset)
    if err != nil {
        return nil, 0, fmt.Errorf("error querying properties: %v", err)
    }
    defer rows.Close()

    var properties []RentalProperty

    // Iterate through results
    for rows.Next() {
        var p RentalProperty
        var amenitiesArr, imagesArr []byte

        // Scan row into struct
        err := rows.Scan(
            &p.ID, 
            &p.PropertyName, 
            &p.Type, 
            &p.Bedrooms, 
            &p.Bathrooms, 
            &p.Price, 
            &p.Location, 
            &amenitiesArr, 
            &imagesArr,
            &p.Rating,
            &p.Reviews,
        )
        if err != nil {
            return nil, 0, fmt.Errorf("error scanning property: %v", err)
        }

        // Convert PostgreSQL array strings to Go slices
        p.Amenities = parsePostgresArray(string(amenitiesArr))
        p.Images = parsePostgresArray(string(imagesArr))

        properties = append(properties, p)
    }

    return properties, totalCount, nil
}

func FilterRentalProperties(
    city string, 
    minPrice, maxPrice float64, 
    bedrooms int, 
    amenities []string,
) ([]RentalProperty, error) {
    // Build dynamic query
    query := `
        SELECT 
            id, 
            property_name, 
            type, 
            bedrooms, 
            bathrooms, 
            price, 
            location, 
            amenities, 
            images,
            rating,
            reviews
        FROM 
            rental_properties 
        WHERE 
            LOWER(location) LIKE LOWER($1)
    `

    // Slice to hold query parameters
    params := []interface{}{"%"+city+"%"}
    paramCount := 2

    // Build where conditions dynamically
    if minPrice > 0 {
        query += fmt.Sprintf(" AND price >= $%d", paramCount)
        params = append(params, minPrice)
        paramCount++
    }

    if maxPrice > 0 {
        query += fmt.Sprintf(" AND price <= $%d", paramCount)
        params = append(params, maxPrice)
        paramCount++
    }

    if bedrooms > 0 {
        query += fmt.Sprintf(" AND bedrooms = $%d", paramCount)
        params = append(params, bedrooms)
        paramCount++
    }

    // Amenities filter
    if len(amenities) > 0 {
        query += fmt.Sprintf(" AND (")
        for i, amenity := range amenities {
            if i > 0 {
                query += " OR "
            }
            query += fmt.Sprintf("$%d = ANY(amenities)", paramCount)
            params = append(params, amenity)
            paramCount++
        }
        query += ")"
    }

    // Add ordering
    query += " ORDER BY rating DESC, price ASC"

    // Execute query
    rows, err := DB.Query(query, params...)
    if err != nil {
        return nil, fmt.Errorf("error filtering properties: %v", err)
    }
    defer rows.Close()

    var properties []RentalProperty

    // Similar scanning logic as GetRentalProperties
    for rows.Next() {
        var p RentalProperty
        var amenitiesArr, imagesArr []byte

        err := rows.Scan(
            &p.ID, 
            &p.PropertyName, 
            &p.Type, 
            &p.Bedrooms, 
            &p.Bathrooms, 
            &p.Price, 
            &p.Location, 
            &amenitiesArr, 
            &imagesArr,
            &p.Rating,
            &p.Reviews,
        )
        if err != nil {
            return nil, fmt.Errorf("error scanning property: %v", err)
        }

        p.Amenities = parsePostgresArray(string(amenitiesArr))
        p.Images = parsePostgresArray(string(imagesArr))

        properties = append(properties, p)
    }

    return properties, nil
}

// Utility function to parse PostgreSQL array strings
func parsePostgresArray(arrayStr string) []string {
    if arrayStr == "" || arrayStr == "{}" {
        return []string{}
    }

    // Remove surrounding curly braces
    arrayStr = arrayStr[1 : len(arrayStr)-1]
    
    var result []string
    var current string
    inQuotes := false
    
    for _, char := range arrayStr {
        switch char {
        case '"':
            inQuotes = !inQuotes
        case ',':
            if !inQuotes {
                result = append(result, current)
                current = ""
                continue
            }
            current += string(char)
        default:
            current += string(char)
        }
    }
    
    if current != "" {
        result = append(result, current)
    }
    
    return result
}