<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>UAE Vacation Rentals</title>
    <style>
        {{template "static/css/main.css"}}
    </style>
</head>
<body>
    <header class="header">
        <div class="search-container">
            <input type="text" placeholder="Dubai, Dubai, United Arab Emirates" class="search-input" id="location">
            <input type="text" placeholder="Select dates" class="search-input" id="dates">
            <input type="text" placeholder="Guests" class="search-input" id="guests">
            <button class="search-button" onclick="searchListings()">Search</button>
        </div>
    </header>

    <div class="filters-container">
        <button class="filter-button active">All</button>
        <button class="filter-button">Apartments</button>
        <button class="filter-button">Villas</button>
        <button class="filter-button">Hotels</button>
        <button class="filter-button">Price: Low to High</button>
        <button class="filter-button">Rating: High to Low</button>
    </div>

    <h1 class="main-title">United Arab Emirates Vacation Rentals & Rent By Owner Homes</h1>

    <div class="listings-grid" id="listingsContainer">
        {{range .Listings}}
            <div class="listing-card">
                <div class="image-gallery">
                    <img src="/static/images/{{index .Images 0}}" alt="{{.Title}}" class="listing-image">
                    <div class="gallery-nav">
                        {{range $index, $img := .Images}}
                            <div class="gallery-dot {{if eq $index 0}}active{{end}}"></div>
                        {{end}}
                    </div>
                </div>
                <button class="heart-button" onclick="toggleFavorite({{.ID}})">♡</button>
                <div class="listing-details">
                    <div class="price">From ${{.Price}}</div>
                    <div class="rating">
                        <span class="rating-score">{{.Rating}}</span>
                        <span>({{.Reviews}} Reviews)</span>
                    </div>
                    <div class="amenities">
                        {{range $index, $amenity := .Amenities}}
                            <span>{{$amenity}}</span>
                            {{if ne $index (sub (len $.Amenities) 1)}}
                                <span>•</span>
                            {{end}}
                        {{end}}
                    </div>
                    <div class="location">{{.Location}}</div>
                    <button class="view-button" onclick="viewListing({{.ID}})">View Availability</button>
                </div>
            </div>
        {{end}}
    </div>

   
    <script src="/static/js/listings.js"></script>
</body>
</html>