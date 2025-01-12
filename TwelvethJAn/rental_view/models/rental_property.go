package models

import (
    "database/sql"
    "fmt"
    "strings"
    
    _ "github.com/lib/pq"
    beego "github.com/beego/beego/v2/server/web"
)

// RentalPropertyModel represents a model for rental properties
type RentalPropertyModel struct {
    db *sql.DB
}

// NewRentalPropertyModel creates a new instance of RentalPropertyModel
func NewRentalPropertyModel() *RentalPropertyModel {
    // Ensure DB is initialized
    if DB == nil {
        InitDB()
    }
    
    return &RentalPropertyModel{
        db: DB,
    }
}

// GetProperties fetches properties using the model's database connection
func (m *RentalPropertyModel) GetProperties(city string, page, pageSize int) ([]RentalProperty, int, error) {
    return GetRentalProperties(city, page, pageSize)
}
// GetPropertiesByType fetches properties of a specific type with pagination
func (m *RentalPropertyModel) GetPropertiesByType(city, propertyType string, page, pageSize int) ([]RentalProperty, int, error) {
    // Calculate offset
    offset := (page - 1) * pageSize

    // Count total properties with type filter
    countQuery := `
        SELECT COUNT(*) 
        FROM rental_properties 
        WHERE LOWER(location) LIKE LOWER($1)
        AND LOWER(type) = LOWER($2)
    `
    var totalCount int
    err := m.db.QueryRow(countQuery, "%"+city+"%", propertyType).Scan(&totalCount)
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
            AND LOWER(type) = LOWER($2)
        ORDER BY 
            created_at DESC
        LIMIT $3 OFFSET $4
    `

    // Execute query
    rows, err := m.db.Query(query, "%"+city+"%", propertyType, pageSize, offset)
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

// FilterProperties applies advanced filtering using the model's database connection
func (m *RentalPropertyModel) FilterProperties(
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
            1=1
    `

    // Slice to hold query parameters
    var params []interface{}
    paramCount := 1

    // City filter (if provided)
    if city != "" {
        query += fmt.Sprintf(" AND LOWER(location) LIKE LOWER($%d)", paramCount)
        params = append(params, "%"+city+"%")
        paramCount++
    }

    // Price range filter
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

    // Bedrooms filter
    if bedrooms > 0 {
        query += fmt.Sprintf(" AND bedrooms = $%d", paramCount)
        params = append(params, bedrooms)
        paramCount++
    }

    // Amenities filter
    if len(amenities) > 0 {
        query += " AND ("
        for i, amenity := range amenities {
            if i > 0 {
                query += " AND "
            }
            query += fmt.Sprintf("$%d = ANY(amenities)", paramCount)
            params = append(params, amenity)
            paramCount++
        }
        query += ")"
    }

    // Add ordering and limit
    query += `
        ORDER BY 
            rating DESC, 
            price ASC
        LIMIT 100
    `

    // Execute query
    rows, err := m.db.Query(query, params...)
    if err != nil {
        return nil, fmt.Errorf("error filtering properties: %v", err)
    }
    defer rows.Close()

    var properties []RentalProperty

    // Scan results
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

        // Parse PostgreSQL arrays
        p.Amenities = parsePostgresArray(string(amenitiesArr))
        p.Images = parsePostgresArray(string(imagesArr))

        properties = append(properties, p)
    }

    // Check for any errors encountered during iteration
    if err = rows.Err(); err != nil {
        return nil, fmt.Errorf("error in row iteration: %v", err)
    }

    return properties, nil
}
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

// Database connection variable
var DB *sql.DB

// InitDB initializes the database connection
func InitDB() {
    // Read database configuration from app.conf
    host, _ := beego.AppConfig.String("db.host")
    port, _ := beego.AppConfig.String("db.port")
    user, _ := beego.AppConfig.String("db.user")
    password, _ := beego.AppConfig.String("db.password")
    dbname, _ := beego.AppConfig.String("db.name")

    // Construct connection string
    connStr := fmt.Sprintf(
        "host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
        host, port, user, password, dbname,
    )

    // Open database connection
    var err error
    DB, err = sql.Open("postgres", connStr)
    if err != nil {
        panic(fmt.Sprintf("Error opening database: %v", err))
    }

    // Verify connection
    err = DB.Ping()
    if err != nil {
        panic(fmt.Sprintf("Error connecting to the database: %v", err))
    }
}
// In models/rental_property.go
func (m *RentalPropertyModel) GetUniquePropertyTypes() ([]string, error) {
    query := `
        SELECT DISTINCT type 
        FROM rental_properties 
        WHERE type IS NOT NULL 
        ORDER BY type
    `

    rows, err := m.db.Query(query)
    if err != nil {
        return nil, fmt.Errorf("error fetching property types: %v", err)
    }
    defer rows.Close()

    var types []string
    for rows.Next() {
        var propertyType string
        if err := rows.Scan(&propertyType); err != nil {
            return nil, fmt.Errorf("error scanning property type: %v", err)
        }
        types = append(types, propertyType)
    }

    // Check for any errors encountered during iteration
    if err = rows.Err(); err != nil {
        return nil, fmt.Errorf("error in row iteration: %v", err)
    }

    return types, nil
}
// GetRentalProperties fetches properties with pagination
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

// FilterRentalProperties applies advanced filtering

// FilterRentalProperties applies advanced filtering
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
            1=1
    `

    // Slice to hold query parameters
    var params []interface{}
    paramCount := 1

    // City filter (if provided)
    if city != "" {
        query += fmt.Sprintf(" AND LOWER(location) LIKE LOWER($%d)", paramCount)
        params = append(params, "%"+city+"%")
        paramCount++
    }

    // Price range filter
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

    // Bedrooms filter
    if bedrooms > 0 {
        query += fmt.Sprintf(" AND bedrooms = $%d", paramCount)
        params = append(params, bedrooms)
        paramCount++
    }

    // Amenities filter
    if len(amenities) > 0 {
        query += " AND ("
        for i, amenity := range amenities {
            if i > 0 {
                query += " AND "
            }
            query += fmt.Sprintf("$%d = ANY(amenities)", paramCount)
            params = append(params, amenity)
            paramCount++
        }
        query += ")"
    }

    // Add ordering and limit
    query += `
        ORDER BY 
            rating DESC, 
            price ASC
        LIMIT 100
    `

    // Execute query
    rows, err := DB.Query(query, params...)
    if err != nil {
        return nil, fmt.Errorf("error filtering properties: %v", err)
    }
    defer rows.Close()

    var properties []RentalProperty

    // Scan results
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

        // Parse PostgreSQL arrays
        p.Amenities = parsePostgresArray(string(amenitiesArr))
        p.Images = parsePostgresArray(string(imagesArr))

        properties = append(properties, p)
    }

    // Check for any errors encountered during iteration
    if err = rows.Err(); err != nil {
        return nil, fmt.Errorf("error in row iteration: %v", err)
    }

    return properties, nil
}

// parsePostgresArray parses PostgreSQL array strings
func parsePostgresArray(arrayStr string) []string {
    // Handle empty or null array
    if arrayStr == "" || arrayStr == "{}" {
        return []string{}
    }

    // Remove surrounding curly braces
    arrayStr = strings.Trim(arrayStr, "{}")

    // Split the string, handling quoted elements
    var result []string
    var current strings.Builder
    inQuotes := false
    escaped := false

    for _, char := range arrayStr {
        switch {
        case escaped:
            current.WriteRune(char)
            escaped = false
        case char == '\\':
            escaped = true
        case char == '"':
            inQuotes = !inQuotes
        case char == ',' && !inQuotes:
            // Trim and add current element
            trimmed := strings.Trim(current.String(), " \"")
            if trimmed != "" {
                result = append(result, trimmed)
            }
            current.Reset()
        default:
            current.WriteRune(char)
        }
    }

    // Add last element
    if current.Len() > 0 {
        trimmed := strings.Trim(current.String(), " \"")
        if trimmed != "" {
            result = append(result, trimmed)
        }
    }

    return result
}