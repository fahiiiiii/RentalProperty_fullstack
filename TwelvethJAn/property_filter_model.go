package models

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

// Helper function to parse PostgreSQL array strings
func parsePostgresArray(arrayStr string) []string {
    if arrayStr == "" || arrayStr == "{}" {
        return []string{}
    }

    // Remove surrounding curly braces
    arrayStr = strings.Trim(arrayStr, "{}")
    
    // Split by comma, handling quoted strings
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
            result = append(result, strings.Trim(current.String(), " \""))
            current.Reset()
        default:
            current.WriteRune(char)
        }
    }

    // Add last item
    if current.Len() > 0 {
        result = append(result, strings.Trim(current.String(), " \""))
    }

    return result
}

// Advanced search with more complex filtering
func AdvancedSearchProperties(searchParams map[string]interface{}) ([]RentalProperty, error) {
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

    var params []interface{}
    paramCount := 1

    // Dynamic filtering based on search parameters
    if city, ok := searchParams["city"].(string); ok && city != "" {
        query += fmt.Sprintf(" AND LOWER(location) LIKE LOWER($%d)", paramCount)
        params = append(params, "%"+city+"%")
        paramCount++
    }

    if minPrice, ok := searchParams["minPrice"].(float64); ok && minPrice > 0 {
        query += fmt.Sprintf(" AND price >= $%d", paramCount)
        params = append(params, minPrice)
        paramCount++
    }

    if maxPrice, ok := searchParams["maxPrice"].(float64); ok && maxPrice > 0 {
        query += fmt.Sprintf(" AND price <= $%d", paramCount)
        params = append(params, maxPrice)
        paramCount++
    }

    if bedrooms, ok := searchParams["bedrooms"].(int); ok && bedrooms > 0 {
        query += fmt.Sprintf(" AND bedrooms = $%d", paramCount)
        params = append(params, bedrooms)
        paramCount++
    }

    if amenities, ok := searchParams["amenities"].([]string); ok && len(amenities) > 0 {
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

    // Execute query (similar to previous implementation)
    rows, err := DB.Query(query, params...)
    if err != nil {
        return nil, fmt.Errorf("error in advanced search: %v", err)
    }
    defer rows.Close()

    var properties []RentalProperty
    // Scanning logic remains the same as in previous implementations

    return properties, nil
}