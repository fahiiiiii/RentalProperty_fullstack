// Use the PropertiesJSON passed from the controller
const listings = JSON.parse('{{.PropertiesJSON | raw}}');

function createListingCard(listing) {
    return `
        <div class="listing-card">
            <div class="image-gallery">
                <img src="/static/images/${listing.images[0] || 'placeholder.jpg'}" alt="${listing.title}" class="listing-image">
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
                        ${index < listing.amenities.length - 1 ? '<span>•</span>' : ''}
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
    // Implement view listing functionality
    console.log(`Viewing listing ${id}`);
}

// Initial render
document.addEventListener('DOMContentLoaded', () => {
    renderListings();

    // Filter buttons functionality
    document.querySelectorAll('.filter-button').forEach(button => {
        button.addEventListener('click', (e) => {
            document.querySelectorAll('.filter-button').forEach(btn =>
                btn.classList.remove('active'));
            e.target.classList.add('active');

            // Implement filtering logic based on button text
            let filteredListings = [...listings];
            switch (e.target.textContent) {
                case 'Price: Low to High':
                    filteredListings.sort((a, b) => a.price - b.price);
                    break;
                case 'Rating: High to Low':
                    filteredListings.sort((a, b) => b.rating - a.rating);
                    break;
                // Add more filter cases as needed
            }
            renderListings(filteredListings);
        });
    });
});