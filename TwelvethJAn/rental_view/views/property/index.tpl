<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>UAE Vacation Rentals</title>
    <link rel="stylesheet" href="/static/css/property-listings.css">
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
                    <img src="/static/images/placeholder-image.jpg" alt="{{.PropertyName}}" class="listing-image">
                    <div class="gallery-nav">
                        {{range $index, $image := .Images}}
                            <div class="gallery-dot {{if eq $index 0}}active{{end}}"></div>
                        {{end}}
                    </div>
                </div>
                <button class="heart-button" onclick="toggleFavorite({{.ID}})">♡</button>
                <div class="listing-details">
                    <div class="price">From ${{.Price | formatNumber}}</div>
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

    <script src="/static/js/property_list.js"></script>
    <script>
        // Initialize listings from server-side data
         const listings = {{.PropertiesJSON | safehtml}};
        
        // Override initial render to use server-side data
        document.addEventListener('DOMContentLoaded', () => {
            renderListings(listings);
        });

        function createListingCard(listing) {
            return `
                <div class="listing-card">
                    <div class="image-gallery">
                        <img src="placeholder-image.jpg" alt="${listing.title}" class="listing-image">
                        <div class="gallery-nav">
                            ${listing.images.map((_, i) => `
                                <div class="gallery-dot ${i === 0 ? 'active' : ''}"></div>
                            `).join('')}
                        </div>
                    </div>
                    <button class="heart-button" onclick="toggleFavorite(${listing.id})">♡</button>
                    <div class="listing-details">
                        <div class="price">From $${listing.price.toLocaleString()}</div>
                        <div class="rating">
                            <span class="rating-score">${listing.rating}</span>
                            <span>(${listing.reviews} Reviews)</span>
                        </div>
                        <div class="amenities">
                            ${listing.amenities.map((amenity, index) => `
                                <span>${amenity}</span>
                                ${index !== listing.amenities.length - 1 ? '<span>•</span>' : ''}
                            `).join('')}
                        </div>
                        <div class="location">${listing.location}</div>
                        <button class="view-button" onclick="viewListing(${listing.id})">View Availability</button>
                    </div>
                </div>
            `;
        }

        function renderListings(filteredListings = listings) {
            const container = document.getElementById('listingsContainer');
            container.innerHTML = filteredListings.map(listing => createListingCard(listing)).join('');
        }

        function searchListings() {
            const location = document.getElementById('location').value.toLowerCase();
            const filteredListings = listings.filter(listing =>
                listing.location.toLowerCase().includes(location)
            );
            renderListings(filteredListings);
        }

        function toggleFavorite(id) {
            const button = event.target;
            button.classList.toggle('active');
            button.textContent = button.classList.contains('active') ? '♥' : '♡';
        }

        function viewListing(id) {
            console.log(`Viewing listing ${id}`);
        }

        // Filter buttons functionality
        document.querySelectorAll('.filter-button').forEach(button => {
            button.addEventListener('click', (e) => {
                document.querySelectorAll('.filter-button').forEach(btn =>
                    btn.classList.remove('active'));
                e.target.classList.add('active');

                let filteredListings = [...listings];
                switch (e.target.textContent) {
                    case 'Price: Low to High':
                        filteredListings.sort((a, b) => a.price - b.price);
                        break;
                    case 'Rating: High to Low':
                        filteredListings.sort((a, b) => b.rating - a.rating);
                        break;
                }
                renderListings(filteredListings);
            });
        });
    </script>
</body>
</html>